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

package ogmios

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	"github.com/savaki/ogmigo/ouroboros/chainsync"
	"golang.org/x/sync/errgroup"
)

// ChainSyncClient provides a client for the chain sync protocol only
type ChainSyncClient struct {
	blocks chan json.RawMessage
	ch     chan struct{}
	conn   *websocket.Conn
	tip    chan struct{} // tip will be published to whenever
	group  *errgroup.Group
}

// NewChainSyncClient returns a new ChainSyncClient
func NewChainSyncClient(ctx context.Context, opts ...Option) (*ChainSyncClient, error) {
	options := buildOptions(opts...)
	logger := options.logger.With(KV("service", "ogmios"))

	conn, _, err := websocket.DefaultDialer.Dial(options.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ogmios, %v: %w", options.endpoint, err)
	}

	logger.Info(ctx, "ogmigo client started")
	defer logger.Info(ctx, "ogmigo client stopped")

	group, ctx := errgroup.WithContext(ctx)
	client := &ChainSyncClient{
		blocks: make(chan json.RawMessage, 8),
		ch:     make(chan struct{}, 64),
		tip:    make(chan struct{}, 1),
		conn:   conn,
		group:  group,
	}

	group.Go(func() error {
		init := []byte(`{"type":"jsonwsp/request","version":"1.0","servicename":"ogmios","methodname":"FindIntersect","args":{"points":["origin"]},"mirror":{"step":"INIT"}}`)
		if err := conn.WriteMessage(websocket.TextMessage, init); err != nil {
			return fmt.Errorf("failed to write FindIntersect: %w", err)
		}

		next := []byte(`{"type":"jsonwsp/request","version":"1.0","servicename":"ogmios","methodname":"RequestNext","args":{}}`)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-client.ch:
				if err := conn.WriteMessage(websocket.TextMessage, next); err != nil {
					return fmt.Errorf("failed to write RequestNext: %w", err)
				}
			}
		}
	})

	group.Go(func() error {
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return fmt.Errorf("failed to read message from ogmios: %w", err)
			}

			select {
			case client.ch <- struct{}{}:
			default:
			}

			if messageType == websocket.PingMessage {
				if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
					return fmt.Errorf("failed to respond with pong to ogmios: %w", err)
				}
			}

			select {
			case <-ctx.Done():
				return nil
			case client.blocks <- data:
				// ok
			}
		}
	})

	for i := 0; i < options.pipeline; i++ {
		select {
		case client.ch <- struct{}{}:
		default:
		}
	}

	return client, nil
}

// Blocks retrieves the next block encoded as JSON
func (c *ChainSyncClient) Blocks() <-chan json.RawMessage {
	return c.blocks
}

// Close the client
func (c *ChainSyncClient) Close() error {
	defer c.group.Wait()
	return c.conn.Close()
}

// ReadNext returns the next chainsync.Response
func (c *ChainSyncClient) ReadNext(ctx context.Context) (chainsync.Response, error) {
	message, err := c.ReadNextMessage(ctx)
	if err != nil {
		return chainsync.Response{}, nil
	}

	var response chainsync.Response
	if err := json.Unmarshal(message, &response); err != nil {
		return chainsync.Response{}, fmt.Errorf("failed to read next: %w", err)
	}

	return response, nil
}

// ReadNextMessage returns the next message for those cases where
// unmarshalling the content isn't required
func (c *ChainSyncClient) ReadNextMessage(ctx context.Context) (json.RawMessage, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case block := <-c.blocks:
		return block, nil
	}
}

func (c *ChainSyncClient) Tip() <-chan struct{} {
	return c.tip
}
