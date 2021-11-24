package chainsync

import (
	"encoding/json"
	"github.com/SundaeSwap-finance/sundae-sync/ouroboros/chainsync/num"
)

type ByronBlock struct {
	Body   ByronBody
	Hash   string
	Header ByronHeader
}

type ByronBody struct {
	DlgPayload    []json.RawMessage `json:"dlgPayload,omitempty"`
	TxPayload     []ByronTxPayload  `json:"txPayload,omitempty"`
	UpdatePayload json.RawMessage   `json:"updatePayload,omitempty"`
}

type ByronHeader struct {
	BlockHeight     num.Int `json:"blockHeight,omitempty"`
	GenesisKey      string
	Epoch           uint32
	Proof           json.RawMessage
	PrevHash        string
	ProtocolMagicId uint64
	ProtocolVersion ProtocolVersion
	Signature       json.RawMessage
	Slot            num.Int `json:"slot,omitempty"`
	SoftwareVersion map[string]interface{}
}

type ByronTxBody struct {
	Inputs  []TxIn
	Outputs []TxOut
}

type ByronTxPayload struct {
	ID      string
	Witness []ByronWitness
}

type ByronWitness struct {
	RedeemWitness map[string]string
}
