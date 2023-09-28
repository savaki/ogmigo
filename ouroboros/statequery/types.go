package statequery

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
)

type EraStart struct {
	Time  time.Duration `json:"time,omitempty"`
	Slot  uint64        `json:"slot,omit"`
	Epoch uint64        `json:"epoch,omit"`
}

type Utxo struct {
	TxIn  chainsync.TxIn
	TxOut chainsync.TxOut
}

func (u Utxo) MarshalJSON() ([]byte, error) {
	txIn, err := json.Marshal(u.TxIn)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Utxo: %w", err)
	}

	txOut, err := json.Marshal(u.TxOut)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Utxo: %w", err)
	}

	utxo := []json.RawMessage{txIn, txOut}
	return json.Marshal(utxo)
}

func (u *Utxo) UnmarshalJSON(data []byte) (err error) {
	var items [2]json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return fmt.Errorf("failed to unmarshal Utxo: %w", err)
	}

	var txIn chainsync.TxIn
	if err := json.Unmarshal(items[0], &txIn); err != nil {
		return fmt.Errorf("failed to unmarshal Utxo: failed to unmarshal TxIn: %w", err)
	}
	var txOut chainsync.TxOut
	if err := json.Unmarshal(items[1], &txOut); err != nil {
		return fmt.Errorf("failed to unmarshal Utxo: failed to unmarshal TxOut: %w", err)
	}

	*u = Utxo{
		TxIn:  txIn,
		TxOut: txOut,
	}

	return nil
}
