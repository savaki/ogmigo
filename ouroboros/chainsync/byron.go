// Copyright 2023 SundaeSwap Labs
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

// Byron block code for Ogmios v6.
// https://github.com/CardanoSolutions/ogmios/blob/d2b1d70ab5e676b5d053817d57ea8220f2b07317/docs/static/api/specification.yaml

package chainsync

import (
	"encoding/json"
)

// BFT Block Root
type ByronBlockBFT struct {
	Type                    string             `json:"type,omitempty"`
	Era                     string             `json:"era,omitempty"`
	ID                      string             `json:"id,omitempty"`
	Ancestor                string             `json:"ancestor,omitempty"`
	Height                  uint64             `json:"height,omitempty"`
	Slot                    uint64             `json:"slot,omitempty"`
	Size                    BlockSize          `json:"size,omitempty"`
	Transactions            []Tx               `json:"transactions,omitempty"`
	OperationalCertificates []json.RawMessage  `json:"operationalCertificates,omitempty"`
	Protocol                ByronProtocol      `json:"protocol,omitempty"`
	Issuer                  ByronBlockIssuer   `json:"issuer,omitempty"`
	Delegate                ByronBlockDelegate `json:"delegate,omitempty"`
}

// EBB Block Type
type ByronBlockEBB struct {
	Type     string `json:"type,omitempty"`
	Era      string `json:"era,omitempty"`
	ID       string `json:"id,omitempty"`
	Ancestor string `json:"ancestor,omitempty"`
	Height   uint64 `json:"height,omitempty"`
}

type ByronBlockDelegate struct {
	VerificationKey string
}

type ByronBlockIssuer struct {
	VerificationKey string
}

type ByronProtocol struct {
	Version  ProtocolVersion
	Id       uint64 // aka magic
	Software map[string]interface{}
	Update   json.RawMessage
}

type ResultByronEBB struct {
	Direction string         `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Block     *ByronBlockEBB `json:"block,omitempty"     dynamodbav:"block,omitempty"`
	Tip       *Tip           `json:"tip,omitempty"       dynamodbav:"tip,omitempty"`
}

type ResultByronBFT struct {
	Direction string         `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Block     *ByronBlockBFT `json:"block,omitempty"     dynamodbav:"block,omitempty"`
	Tip       *Tip           `json:"tip,omitempty"       dynamodbav:"tip,omitempty"`
}

type ResponseByronEBB struct {
	JsonRpc string          `json:"jsonrpc,omitempty" dynamodbav:"jsonrpc,omitempty"`
	Method  string          `json:"method,omitempty"  dynamodbav:"method,omitempty"`
	Result  *ResultByronEBB `json:"result,omitempty"  dynamodbav:"result,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"      dynamodbav:"id,omitempty"`
}

type ResponseByronBFT struct {
	JsonRpc string          `json:"jsonrpc,omitempty" dynamodbav:"jsonrpc,omitempty"`
	Method  string          `json:"method,omitempty"  dynamodbav:"method,omitempty"`
	Result  *ResultByronBFT `json:"result,omitempty"  dynamodbav:"result,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"      dynamodbav:"id,omitempty"`
}
