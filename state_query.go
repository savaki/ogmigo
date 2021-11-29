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

	"github.com/savaki/ogmigo/ouroboros/chainsync"
	"github.com/savaki/ogmigo/ouroboros/statequery"
)

func (c *Client) ChainTip(ctx context.Context) (chainsync.Point, error) {
	var (
		payload = makePayload("Query", Map{"query": "ledgerTip"})
		content struct{ Result chainsync.Point }
	)

	if err := c.query(ctx, payload, &content); err != nil {
		return chainsync.Point{}, err
	}

	return content.Result, nil
}

func (c *Client) CurrentEpoch(ctx context.Context) (uint64, error) {
	var (
		payload = makePayload("Query", Map{"query": "currentEpoch"})
		content struct{ Result uint64 }
	)

	if err := c.query(ctx, payload, &content); err != nil {
		return 0, err
	}

	return content.Result, nil
}

func (c *Client) CurrentProtocolParameters(ctx context.Context) (json.RawMessage, error) {
	var (
		payload = makePayload("Query", Map{"query": "currentProtocolParameters"})
		content struct{ Result json.RawMessage }
	)

	if err := c.query(ctx, payload, &content); err != nil {
		return nil, err
	}

	return content.Result, nil
}

func (c *Client) EraStart(ctx context.Context) (statequery.EraStart, error) {
	var (
		payload = makePayload("Query", Map{"query": "eraStart"})
		content struct{ Result statequery.EraStart }
	)

	if err := c.query(ctx, payload, &content); err != nil {
		return statequery.EraStart{}, err
	}

	return content.Result, nil
}

func (c *Client) UtxosByAddress(ctx context.Context, addresses ...string) ([]statequery.Utxo, error) {
	var (
		payload = makePayload("Query", Map{"query": Map{"utxo": addresses}})
		content struct{ Result []statequery.Utxo }
	)

	if err := c.query(ctx, payload, &content); err != nil {
		return nil, fmt.Errorf("failed to query utxos by address: %w", err)
	}

	return content.Result, nil
}

func (c *Client) UtxosByTxIn(ctx context.Context, txIns ...chainsync.TxIn) ([]statequery.Utxo, error) {
	var (
		payload = makePayload("Query", Map{"query": Map{"utxo": txIns}})
		content struct{ Result []statequery.Utxo }
	)

	if err := c.query(ctx, payload, &content); err != nil {
		return nil, fmt.Errorf("failed to query utxos by address: %w", err)
	}

	return content.Result, nil
}
