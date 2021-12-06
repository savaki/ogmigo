// Copyright 2021 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ogmigo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sort"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/savaki/ogmigo/ouroboros/chainsync"
	"golang.org/x/sync/errgroup"
)

// ChainSyncFunc callback containing json encoded chainsync.Response
type ChainSyncFunc func(ctx context.Context, data []byte) error

// ChainSyncOptions configuration parameters
type ChainSyncOptions struct {
	minSlot uint64           // minSlot to begin invoking ChainSyncFunc; 0 for always invoke func
	points  chainsync.Points // points to attempt initial intersection
	store   Store            // store of points
}

func buildChainSyncOptions(opts ...ChainSyncOption) ChainSyncOptions {
	var options ChainSyncOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.store == nil {
		options.store = nopStore{}
	}
	return options
}

// ChainSyncOption provides functional options for ChainSync
type ChainSyncOption func(opts *ChainSyncOptions)

// WithMinSlot ignores any activity prior to the specified slot
func WithMinSlot(slot uint64) ChainSyncOption {
	return func(opts *ChainSyncOptions) {
		opts.minSlot = slot
	}
}

// WithPoints allows starting from an optional point
func WithPoints(points ...chainsync.Point) ChainSyncOption {
	return func(opts *ChainSyncOptions) {
		opts.points = points
	}
}

// WithStore specifies store to persist points to; defaults to no persistence
func WithStore(store Store) ChainSyncOption {
	return func(opts *ChainSyncOptions) {
		opts.store = store
	}
}

type closeFunc func() error

func (fn closeFunc) Close() error {
	return fn()
}

// ChainSync replays the blockchain by invoking the callback for each block
// By default, ChainSync stores no checkpoints and always restarts from origin.  These can
// be overridden via WithPoints and WithStore
func (c *Client) ChainSync(ctx context.Context, callback ChainSyncFunc, opts ...ChainSyncOption) (io.Closer, error) {
	options := buildChainSyncOptions(opts...)

	conn, _, err := websocket.DefaultDialer.Dial(c.options.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ogmios, %v: %w", c.options.endpoint, err)
	}

	init, err := getInit(ctx, options.store, options.points...)
	if err != nil {
		return nil, fmt.Errorf("failed to create init message: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		c.options.logger.Info("ogmigo chainsync started")
		defer c.options.logger.Info("ogmigo chainsync stopped")
		<-ctx.Done()
		return nil
	})

	var connState int64 // 0 - open, 1 - closing, 2 - closed
	group.Go(func() error {
		<-ctx.Done()
		atomic.AddInt64(&connState, 1)
		if err := conn.Close(); err != nil {
			return err
		}
		atomic.AddInt64(&connState, 1)
		return nil
	})

	// prime the pump
	ch := make(chan struct{}, 64)
	for i := 0; i < c.options.pipeline; i++ {
		select {
		case ch <- struct{}{}:
		default:
		}
	}

	group.Go(func() error {
		if err := conn.WriteMessage(websocket.TextMessage, init); err != nil {
			var oe *net.OpError
			if ok := errors.As(err, &oe); ok {
				if v := atomic.LoadInt64(&connState); v > 0 {
					return nil // connection closed
				}
			}
			return fmt.Errorf("failed to write FindIntersect: %w", err)
		}

		next := []byte(`{"type":"jsonwsp/request","version":"1.0","servicename":"ogmios","methodname":"RequestNext","args":{}}`)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ch:
				if err := conn.WriteMessage(websocket.TextMessage, next); err != nil {
					return fmt.Errorf("failed to write RequestNext: %w", err)
				}
			}
		}
	})

	group.Go(func() error {
		checkSlot := options.minSlot > 0
		last := newCircular(3)
		for n := uint64(1); ; n++ {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				var oe *net.OpError
				if ok := errors.As(err, &oe); ok {
					if v := atomic.LoadInt64(&connState); v > 0 {
						return nil // connection closed
					}
				}
				return fmt.Errorf("failed to read message from ogmios: %w", err)
			}

			select {
			case <-ctx.Done():
				if point, ok := getPoint(last.list()...); ok {
					if err := options.store.Save(context.Background(), point); err != nil {
						return fmt.Errorf("chainsync client failed: %w", err)
					}
				}
				return nil
			case ch <- struct{}{}:
				// request the next message
			default:
				// pump is full
			}

			switch messageType {
			case websocket.BinaryMessage:
				c.options.logger.Info("skipping unexpected binary message")
				continue

			case websocket.CloseMessage:
				if point, ok := getPoint(last.list()...); ok {
					if err := options.store.Save(context.Background(), point); err != nil {
						return fmt.Errorf("chainsync client failed: %w", err)
					}
				}
				return nil

			case websocket.PingMessage:
				if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					return fmt.Errorf("failed to respond with pong to ogmios: %w", err)
				}
				continue

			case websocket.PongMessage:
				continue

			case websocket.TextMessage:
				// ok
			}

			// allow rapid bypassing of earlier slots
			if checkSlot {
				if point, ok := getPoint(data); ok {
					if ps, ok := point.PointStruct(); ok {
						if ps.Slot < options.minSlot {
							continue
						}
						checkSlot = false
					}
				}
			}

			if err := callback(ctx, data); err != nil {
				return fmt.Errorf("chainsync stopped: callback failed: %w", err)
			}

			// periodically save points to the store to allow graceful recovery
			if n%c.options.saveInterval == 0 {
				if point, ok := getPoint(last.prefix(data)...); ok {
					if err := options.store.Save(ctx, point); err != nil {
						return fmt.Errorf("chainsync client failed: %w", err)
					}
				}
			}
			last.add(data)
		}
	})

	shutdown := func() error {
		c.options.logger.Info("ogmigo shutdown requested")
		defer c.options.logger.Info("ogmigo shutdown completed")

		cancel()
		if err := group.Wait(); err != nil {
			return err
		}
		return nil
	}

	return closeFunc(shutdown), nil
}

func getInit(ctx context.Context, store Store, points ...chainsync.Point) (data []byte, err error) {
	if len(points) == 0 {
		points, err = store.Load(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve points from store: %w", err)
		}
		if len(points) == 0 {
			points = append(points, chainsync.Origin)
		}
	}
	sort.Sort(chainsync.Points(points))
	if len(points) > 5 {
		points = points[0:5]
	}

	init := Map{
		"type":        "jsonwsp/request",
		"version":     "1.0",
		"servicename": "ogmios",
		"methodname":  "FindIntersect",
		"args":        Map{"points": points},
		"mirror":      Map{"step": "INIT"},
	}
	return json.Marshal(init)
}

// getPoint returns the first point from the list of json encoded chainsync.Responses provided
// multiple Responses allow for the possibility of a Rollback being included in the set
func getPoint(data ...[]byte) (chainsync.Point, bool) {
	for _, d := range data {
		if len(d) == 0 {
			continue
		}

		var response chainsync.Response
		if err := json.Unmarshal(d, &response); err == nil {
			if response.Result != nil && response.Result.RollForward != nil {
				if point, ok := getPointFromBlock(response.Result.RollForward.Block); ok {
					return point, true
				}
			}
		}
	}
	return chainsync.Point{}, false
}

// getPointFromBlock extracts a point from a block regardless of which era we were in
func getPointFromBlock(block chainsync.RollForwardBlock) (chainsync.Point, bool) {
	if byron := block.Byron; byron != nil {
		return chainsync.PointStruct{
			BlockNo: byron.Header.BlockHeight,
			Hash:    byron.Hash,
			Slot:    byron.Header.Slot,
		}.Point(), true
	}

	var header chainsync.BlockHeader
	switch {
	case block.Allegra != nil:
		header = block.Allegra.Header
	case block.Alonzo != nil:
		header = block.Alonzo.Header
	case block.Byron != nil:
		// already handled above
	case block.Mary != nil:
		header = block.Mary.Header
	case block.Shelley != nil:
		header = block.Shelley.Header
	default:
		return chainsync.Point{}, false
	}

	return chainsync.PointStruct{
		BlockNo: header.BlockHeight,
		Hash:    header.BlockHash,
		Slot:    header.Slot,
	}.Point(), true
}
