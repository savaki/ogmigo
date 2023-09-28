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

package badgerstore

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/dgraph-io/badger/v3"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
)

type Store struct {
	db      *badger.DB
	counter int64
	prefix  []byte
}

func New(db *badger.DB, prefix string) *Store {
	return &Store{
		db:     db,
		prefix: []byte(strings.TrimRight(prefix, "/") + "/"),
	}
}

// Save the point; save will be called multiple times and should only
// keep track of the most recent points
func (s *Store) Save(_ context.Context, point chainsync.Point) error {
	data, err := json.Marshal(point)
	if err != nil {
		return fmt.Errorf("failed to save point: %w", err)
	}

	v := atomic.AddInt64(&s.counter, 1) % 10
	key := append(s.prefix, []byte(strconv.FormatInt(v, 10))...)

	tx := s.db.NewTransaction(true)
	if err := tx.Set(key, data); err != nil {
		return fmt.Errorf("failed to save point: set failed: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to save point: commit failed: %w", err)
	}

	return s.db.Sync()
}

// Load saved points
func (s *Store) Load(context.Context) (chainsync.Points, error) {
	tx := s.db.NewTransaction(false)
	iter := tx.NewIterator(badger.DefaultIteratorOptions)
	defer iter.Close()

	var pp chainsync.Points
	for iter.Seek(s.prefix); iter.ValidForPrefix(s.prefix); iter.Next() {
		var p chainsync.Point
		unmarshal := func(val []byte) error { return json.Unmarshal(val, &p) }

		if err := iter.Item().Value(unmarshal); err != nil {
			return nil, fmt.Errorf("failed to load points: %w", err)
		}

		pp = append(pp, p)
	}

	sort.Sort(pp)

	return pp, nil
}
