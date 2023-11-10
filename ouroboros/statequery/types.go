package statequery

import (
	"encoding/json"
	"math/big"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/shared"
)

type EraStart struct {
	Time  EraSeconds `json:"time,omitempty"`
	Slot  big.Int    `json:"slot,omitempty"`
	Epoch big.Int    `json:"epoch,omitempty"`
}

type EraSeconds struct {
	Seconds big.Int `json:"seconds"`
}

type EraMilliseconds struct {
	Milliseconds big.Int `json:"milliseconds"`
}

type UtxoData struct {
	InTransaction shared.UtxoTxID `json:"transaction"`
	InIndex       uint32          `json:"index"`
	OutAddress    string          `json:"address"`
	OutValue      shared.Value    `json:"value"`
	OutDatumHash  string          `json:"datumHash,omitempty"`
	OutDatum      string          `json:"datum,omitempty"`
	OutScript     json.RawMessage `json:"script,omitempty"`
}
