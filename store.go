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
	"strconv"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
)

// Store allows points to be saved and retrieved to allow graceful recovery
// after shutdown
type Store interface {
	// Save the point; save will be called multiple times and should only
	// keep track of the most recent points
	Save(ctx context.Context, point chainsync.Point) error
	// Load saved points
	Load(ctx context.Context) (chainsync.Points, error)
}

type loggingStore struct {
	logger Logger
}

// NewLoggingStore logs Save requests, but does not actually save points
func NewLoggingStore(logger Logger) Store {
	return &loggingStore{
		logger: logger,
	}
}

func (l *loggingStore) Save(_ context.Context, point chainsync.Point) error {
	var kvs []KeyValue
	if ps, ok := point.PointStruct(); ok {
		kvs = append(kvs, KV("slot", strconv.FormatUint(ps.Slot, 10)))
		kvs = append(kvs, KV("block", strconv.FormatUint(ps.BlockNo, 10)))
		kvs = append(kvs, KV("hash", ps.Hash))
	}
	l.logger.Info("save point", kvs...)
	return nil
}

func (l *loggingStore) Load(_ context.Context) (chainsync.Points, error) {
	return nil, nil
}

type nopStore struct {
}

func (n nopStore) Save(context.Context, chainsync.Point) error    { return nil }
func (n nopStore) Load(context.Context) (chainsync.Points, error) { return nil, nil }
