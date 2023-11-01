package statequery

import (
	"encoding/json"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
	"math/big"
)

type EraStart struct {
	Time  EraSeconds `json:"time,omitempty"`
	Slot  big.Int    `json:"slot,omit"`
	Epoch big.Int    `json:"epoch,omit"`
}

type EraSeconds struct {
	Seconds big.Int `json:"seconds"`
}

type EraMilliseconds struct {
	Milliseconds big.Int `json:"milliseconds"`
}

type Utxo struct {
	Transaction UtxoTxID        `json:"transaction"`
	Index       uint32          `json:"index"`
	Address     string          `json:"address"`
	Value       Value           `json:"value"`
	DatumHash   string          `json:"datumHash,omitempty"`
	Datum       string          `json:"datum,omitempty"`
	Script      json.RawMessage `json:"script,omitempty"`
}

type UtxoTxID struct {
	ID string `json:"id"`
}

type Value map[string]map[string]num.Int

type AssetId struct {
	Policy string
	Token  string
}

const AdaPolicy = "ada"
const AdaName = "lovelace"

func AdaAssetId() AssetId {
	return AssetId{
		Policy: AdaPolicy,
		Token:  AdaName,
	}
}

func (v Value) AdaLovelace() num.Int {
	return v.AssetAmount(AdaAssetId())
}

func (v Value) AssetAmount(asset AssetId) num.Int {
	if nested, ok := v[asset.Policy]; ok {
		return nested[asset.Token]
	}
	return num.Int64(0)
}

func (v Value) Assets() map[string]map[string]num.Int {
	policies := make(map[string]map[string]num.Int, 0)
	for policy, tokenMap := range v {
		if policy == AdaPolicy {
			continue
		}
		policies[policy] = make(map[string]num.Int, 0)
		for token, quantity := range tokenMap {
			policies[policy][token] = quantity
		}
	}
	return policies
}
