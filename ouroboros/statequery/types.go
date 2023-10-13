package statequery

import (
	"encoding/json"
	"fmt"
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

type Value struct {
	Ada    int64
	Assets map[string]map[string]int64
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var m map[string]map[string]int64
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	adaAmt, ok := m["ada"]["lovelace"]
	if !ok {
		return fmt.Errorf("statequery: Value.UnmarshalJSON: key 'ada' missing")
	}
	delete(m, "ada")
	v.Ada = adaAmt
	v.Assets = m
	return nil
}
