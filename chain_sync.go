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
	"os"
	"sort"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thuannguyen2010/ogmigo/ouroboros/chainsync"
	"golang.org/x/sync/errgroup"
)

var (
	delayPeriods = []time.Duration{5, 10, 20, 40, 80, 160, 300}
	periodIndex  = 0
)

// ChainSync provides control over a given ChainSync connection
type ChainSync struct {
	cancel context.CancelFunc
	errs   chan error
	done   chan struct{}
	err    error
	logger Logger
}

// Done indicates the ChainSync has terminated prematurely
func (c *ChainSync) Done() <-chan struct{} {
	return c.done
}

// Close the ChainSync connection
func (c *ChainSync) Close() error {
	c.cancel()
	select {
	case v := <-c.errs:
		c.err = v
	default:
		// err already set
	}
	return c.err
}

// ChainSyncFunc callback containing json encoded chainsync.Response
type ChainSyncFunc func(ctx context.Context, data []byte) error

// ChainSyncOptions configuration parameters
type ChainSyncOptions struct {
	minSlot   uint64           // minSlot to begin invoking ChainSyncFunc; 0 for always invoke func
	points    chainsync.Points // points to attempt initial intersection
	reconnect bool             // reconnect to ogmios if connection drops
	store     Store            // store of points
	useV6     bool             // useV6 decides using v6 or v5 payload
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

// WithUseV6 decides using ogmios v6 or v5 interface
func WithUseV6(useV6 bool) ChainSyncOption {
	return func(opts *ChainSyncOptions) {
		opts.useV6 = useV6
	}
}

// WithReconnect attempt to reconnect to ogmios if connection drops
func WithReconnect(enabled bool) ChainSyncOption {
	return func(opts *ChainSyncOptions) {
		opts.reconnect = enabled
	}
}

// WithStore specifies store to persist points to; defaults to no persistence
func WithStore(store Store) ChainSyncOption {
	return func(opts *ChainSyncOptions) {
		opts.store = store
	}
}

// ChainSync replays the blockchain by invoking the callback for each block
// By default, ChainSync stores no checkpoints and always restarts from origin.  These can
// be overridden via WithPoints and WithStore
func (c *Client) ChainSync(ctx context.Context, callback ChainSyncFunc, opts ...ChainSyncOption) (*ChainSync, error) {
	options := buildChainSyncOptions(opts...)

	done := make(chan struct{})
	errs := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(done)

		var err error

		for {
			err = c.doChainSync(ctx, callback, options)
			if err != nil {
				// always retry
				if options.reconnect {
					timeout := delayPeriods[periodIndex] * time.Second
					if periodIndex < len(delayPeriods)-1 {
						periodIndex += 1
					}
					c.options.logger.Info("websocket connection error: will retry",
						KV("delay", timeout.Round(time.Millisecond).String()),
						KV("err", err.Error()),
					)

					select {
					case <-ctx.Done():
						return
					case <-time.After(timeout):
						continue
					}
				}
			}

			break
		}

		errs <- err
	}()

	return &ChainSync{
		cancel: cancel,
		errs:   errs,
		done:   done,
		logger: c.logger,
	}, nil
}

func (c *Client) doChainSync(ctx context.Context, callback ChainSyncFunc, options ChainSyncOptions) error {
	conn, _, err := websocket.DefaultDialer.Dial(c.options.endpoint, nil)
	if err != nil {
		c.logger.Error(err, "failed to connect to ogmios", KV("endpoint", c.options.endpoint))
		return fmt.Errorf("failed to connect to ogmios, %v: %w", c.options.endpoint, err)
	}
	var init, next []byte
	if options.useV6 {
		next = []byte(`{"jsonrpc":"2.0","method":"nextBlock"}`)
		init, err = getInitV6(ctx, options.store, options.points...)
	} else {
		next = []byte(`{"type":"jsonwsp/request","version":"1.0","servicename":"ogmios","methodname":"RequestNext","args":{}}`)
		init, err = getInit(ctx, options.store, options.points...)
	}
	if err != nil {
		c.logger.Error(err, "failed to create init message")
		return fmt.Errorf("failed to create init message: %w", err)
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		c.options.logger.Info("ogmigo chainsync started")
		periodIndex = 0 // reset delay period
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
			c.logger.Error(err, "failed to write FindIntersect")
			return fmt.Errorf("failed to write FindIntersect: %w", err)
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ch:
				if err := conn.WriteMessage(websocket.TextMessage, next); err != nil {
					c.logger.Error(err, "failed to write RequestNext")
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
				c.logger.Error(err, "failed to read message from ogmios")
				return NewWrappedReadMessageError("failed to read message from ogmios: %w", err)
			}

			select {
			case <-ctx.Done():
				if point, ok := getPoint(last.list()...); ok {
					if err := options.store.Save(context.Background(), point); err != nil {
						c.logger.Error(err, "chainsync client failed")
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
						c.logger.Error(err, "chainsync client failed")
						return fmt.Errorf("chainsync client failed: %w", err)
					}
				}
				return nil

			case websocket.PingMessage:
				if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					c.logger.Error(err, "failed to respond with pong to ogmios")
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
				c.logger.Error(err, "chainsync stopped: callback failed")
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
	return group.Wait()
}

func getInit(ctx context.Context, store Store, pp ...chainsync.Point) (data []byte, err error) {
	points, err := store.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve points from store: %w", err)
	}
	if len(points) == 0 {
		points = append(points, pp...)
	}
	if len(points) == 0 {
		points = append(points, chainsync.Origin)
	}
	sort.Sort(points)
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

func getInitV6(ctx context.Context, store Store, pp ...chainsync.Point) (data []byte, err error) {
	points, err := store.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve points from store: %w", err)
	}
	if len(points) == 0 {
		points = append(points, pp...)
	}
	if len(points) == 0 {
		points = append(points, chainsync.Origin)
	}
	sort.Sort(points)
	if len(points) > 5 {
		points = points[0:5]
	}
	pointsV6 := points.ConvertToV6()
	init := Map{
		"jsonrpc": "2.0",
		"method":  "findIntersection",
		"params": Map{
			"points": pointsV6,
		},
		"id": "findIntersect",
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
				ps := response.Result.RollForward.Block.PointStruct()
				return ps.Point(), true
			}
		}
	}
	return chainsync.Point{}, false
}

// isTemporaryError returns true if the error is recoverable
func isTemporaryError(err error) bool {
	var wce *websocket.CloseError
	if ok := errors.As(err, &wce); ok && wce.Code == websocket.CloseAbnormalClosure {
		return true
	}

	if ok := errors.Is(err, websocket.ErrBadHandshake); ok {
		return true
	}

	var noe *net.OpError
	if ok := errors.As(err, &noe); ok {
		var sce *os.SyscallError
		if ok := errors.As(noe.Err, &sce); ok && sce.Syscall == "connect" {
			return true
		}
		return noe.Temporary()
	}

	// handle the generic temporary error
	var temp interface{ Temporary() bool }
	if ok := errors.As(err, &temp); ok {
		return temp.Temporary()
	}

	return false
}
