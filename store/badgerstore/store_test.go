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
	"reflect"
	"testing"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/dgraph-io/badger/v3"
)

func TestStore_Load(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("test-db"))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	defer db.Close()

	var (
		ctx   = context.Background()
		a     = chainsync.PointStruct{Slot: 10}
		b     = chainsync.PointStruct{Slot: 20}
		c     = chainsync.PointStruct{Slot: 30}
		store = New(db, "points")
	)

	err = store.Save(ctx, a.Point())
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	err = store.Save(ctx, b.Point())
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	err = store.Save(ctx, c.Point())
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	points, err := store.Load(ctx)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if got, want := len(points), 3; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	want := chainsync.Points{c.Point(), b.Point(), a.Point()}
	if got := points; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}
