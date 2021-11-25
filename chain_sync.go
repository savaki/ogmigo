package ogmigo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/gorilla/websocket"
	"github.com/savaki/ogmigo/ouroboros/chainsync"
	"golang.org/x/sync/errgroup"
)

// ChainSyncFunc callback containing json encoded chainsync.Response
type ChainSyncFunc func(ctx context.Context, data []byte) error

func (c *Client) ChainSync(ctx context.Context, store Store, callback ChainSyncFunc, points ...chainsync.Point) (func() error, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.options.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ogmios, %v: %w", c.options.endpoint, err)
	}

	init, err := getInit(ctx, store, points...)
	if err != nil {
		return nil, fmt.Errorf("failed to create init message: %w", err)
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		c.options.logger.Info(ctx, "ogmigo chainsync started")
		defer c.options.logger.Info(ctx, "ogmigo chainsync stopped")
		<-ctx.Done()
		return nil
	})
	group.Go(func() error {
		<-ctx.Done()
		return conn.Close()
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
			fmt.Printf("WriteMessage: %T\n", err)
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
		last := newCircular(3)
		for n := uint64(1); ; n++ {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("ReadMessage: 	%T\n", err)
				if errors.Is(err, io.EOF) {
					return nil
				}
				return fmt.Errorf("failed to read message from ogmios: %w", err)
			}

			select {
			case <-ctx.Done():
				if point, ok := getPoint(last.list()...); ok {
					if err := store.Save(context.Background(), point); err != nil {
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
				c.options.logger.Info(ctx, "skipping unexpected binary message")
				continue

			case websocket.CloseMessage:
				if point, ok := getPoint(last.list()...); ok {
					if err := store.Save(context.Background(), point); err != nil {
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

			if err := callback(ctx, data); err != nil {
				return fmt.Errorf("chainsync stopped: callback failed: %w", err)
			}

			// periodically save points to the store to allow graceful recovery
			if n%c.options.saveInterval == 0 {
				if point, ok := getPoint(last.prefix(data)...); ok {
					if err := store.Save(ctx, point); err != nil {
						return fmt.Errorf("chainsync client failed: %w", err)
					}
				}
			}
			last.add(data)
		}
	})

	return group.Wait, nil
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

	init := map[string]interface{}{
		"type":        "jsonwsp/request",
		"version":     "1.0",
		"servicename": "ogmios",
		"methodname":  "FindIntersect",
		"args":        map[string]interface{}{"points": points},
		"mirror":      map[string]interface{}{"step": "INIT"},
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
