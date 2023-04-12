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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var fault = []byte(`jsonwsp/fault`)

func (c *Client) query(ctx context.Context, payload interface{}, v interface{}) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		ch     = make(chan error, 1)
		conn   *websocket.Conn
		closed int64 // ensures close is only called once
	)
	go func() {
		<-ctx.Done()
		if conn != nil {
			ch <- ctx.Err()
			if v := atomic.AddInt64(&closed, 1); v == 1 {
				conn.Close()
			}
		}
	}()

	conn, _, err = c.dialer.DialContext(ctx, c.options.endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to ogmios, %v: %w", c.options.endpoint, err)
	}
	defer func() {
		if v := atomic.AddInt64(&closed, 1); v == 1 {
			conn.Close()
		} else {
			err = <-ch
		}
	}()

	if err := conn.WriteJSON(payload); err != nil {
		return fmt.Errorf("failed to submit request: %w", err)
	}

	var raw json.RawMessage
	if err := conn.ReadJSON(&raw); err != nil {
		return fmt.Errorf("failed to read json response: %w", err)
	}

	if bytes.Contains(raw, fault) {
		var e Error
		if err := json.Unmarshal(raw, &e); err != nil {
			return fmt.Errorf("failed to decode error: %w", err)
		}
		return e
	}

	if v != nil {
		if err := json.Unmarshal(raw, v); err != nil {
			return fmt.Errorf("failed to unmarshal contents: %w", err)
		}
	}

	return nil
}
