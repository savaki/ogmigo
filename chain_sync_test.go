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
	"github.com/savaki/ogmigo/ouroboros/chainsync"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_ChainSync(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var (
		p       = message.NewPrinter(language.English)
		counter int64
		read    int64
	)

	client := New(WithEndpoint(endpoint))
	var callback ChainSyncFunc = func(ctx context.Context, data []byte) error {
		var response chainsync.Response
		decoder := json.NewDecoder(bytes.NewReader(data)) // use decoder to check for unknown fields
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&response)
		if err != nil {
			fmt.Println(string(data))
			t.Fatalf("got %v; want nil", err)
		}

		read += int64(len(data))
		if v := atomic.AddInt64(&counter, 1); v%1e3 == 0 {
			var blockNo uint64
			if response.Result != nil && response.Result.RollForward != nil {
				if ps, ok := response.Result.RollForward.Tip.PointStruct(); ok {
					blockNo = ps.BlockNo
				}
			}
			log.Printf("read: block=%v, n=%v, read=%v", blockNo, p.Sprintf("%d", v), p.Sprintf("%d", read))
		}

		return nil
	}

	closer, err := client.ChainSync(ctx, callback, WithStore(echoStore{}))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	<-time.After(5 * time.Second)
	if err := closer.Close(); err != nil {
		t.Fatalf("got %v; want nil", err)
	}
}

type echoStore struct {
}

func (e echoStore) Save(_ context.Context, p chainsync.Point) error {
	fmt.Print("Save => ")
	return json.NewEncoder(os.Stdout).Encode(p)
}

func (e echoStore) Load(context.Context) (chainsync.Points, error) {
	return nil, nil
}
