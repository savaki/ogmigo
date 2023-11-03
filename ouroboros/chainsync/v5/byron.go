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

package v5

import (
	"encoding/json"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
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
	BlockHeight     uint64 `json:"blockHeight,omitempty"`
	GenesisKey      string
	Epoch           uint32
	Proof           json.RawMessage
	PrevHash        string
	ProtocolMagicId uint64
	ProtocolVersion chainsync.ProtocolVersion
	Signature       json.RawMessage
	Slot            uint64 `json:"slot,omitempty"`
	SoftwareVersion map[string]interface{}
}

type ByronTxBody struct {
	Inputs  []chainsync.TxIn
	Outputs []chainsync.TxOut
}

type ByronTxPayload struct {
	ID      string
	Witness []ByronWitness
}

type ByronWitness struct {
	RedeemWitness map[string]string
}
