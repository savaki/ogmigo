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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/shared"
)

func TestClient_ChainTip(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	point, err := client.ChainTip(ctx)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	ps, ok := point.PointStruct()
	if !ok {
		t.Fatalf("got false; want true")
	}
	if ps.ID == "" {
		t.Fatalf("got blank; want not blank")
	}
	if ps.Slot == 0 {
		t.Fatalf("got zero; want not zero")
	}
}

func TestClient_EraSummaries(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	history, err := client.EraSummaries(ctx)
	if err != nil {
		t.Fatalf("got %#v; want nil", err)
	}
	if len(history.Summaries) == 0 {
		t.Fatalf("got empty era history; want nonempty")
	}
}

func TestClient_CurrentEpoch(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	params, err := client.CurrentEpoch(ctx)
	if err != nil {
		t.Fatalf("got %#v; want nil", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(params)
}

func TestClient_CurrentProtocolParameters(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	params, err := client.CurrentProtocolParameters(ctx)
	if err != nil {
		t.Fatalf("got %#v; want nil", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(params)
}

func TestClient_GenesisConfig(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	params, err := client.GenesisConfig(ctx, shared.ShelleyEra)
	if err != nil {
		t.Fatalf("got %#v; want nil", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(params)
}

func TestClient_EraStart(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	eraStart, err := client.EraStart(ctx)
	if err != nil {
		t.Fatalf("got %#v; want nil", err)
	}

	start := time.Now().Add(-time.Duration(eraStart.Time.Seconds.Uint64()))
	fmt.Println(start)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(eraStart)
}

func TestClient_UtxosByAddress(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	utxos, err := client.UtxosByAddress(ctx, "addr_test1qz6m03tdfm5raxr00fsw7p8v79ptfveaptar9a56zqz09kqkazwhq98h9v8gnk3wm5uvevzvd642zm7778afv0evwqgqfuy84f")
	if err != nil {
		t.Fatalf("got %#v; want nil", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(utxos)
}

func TestClient_UtxosByTxIn(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	utxos, err := client.UtxosByTxIn(ctx, chainsync.TxInQuery{
		Transaction: chainsync.UtxoTxID{
			ID: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		Index: 0,
	})

	if err != nil {
		t.Fatalf("got %#v; want nil", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(utxos)
}
