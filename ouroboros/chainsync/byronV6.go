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
type ByronBlockBFTV6 struct {
	Type                    string               `json:"type,omitempty"`
	Era                     string               `json:"era,omitempty"`
	ID                      string               `json:"id,omitempty"`
	Ancestor                string               `json:"ancestor,omitempty"`
	Height                  uint64               `json:"height,omitempty"`
	Slot                    uint64               `json:"slot,omitempty"`
	Size                    BlockSizeV6          `json:"size,omitempty"`
	Transactions            []TxV6               `json:"transactions,omitempty"`
	OperationalCertificates []json.RawMessage    `json:"operationalCertificates,omitempty"`
	Protocol                ByronProtocolV6      `json:"protocol,omitempty"`
	Issuer                  ByronBlockIssuerV6   `json:"issuer,omitempty"`
	Delegate                ByronBlockDelegateV6 `json:"delegate,omitempty"`
}

// EBB Block Type
type ByronBlockEBBV6 struct {
	Type     string `json:"type,omitempty"`
	Era      string `json:"era,omitempty"`
	ID       string `json:"id,omitempty"`
	Ancestor string `json:"ancestor,omitempty"`
	Height   uint64 `json:"height,omitempty"`
}

type ByronBlockDelegateV6 struct {
	VerificationKey string
}

type ByronBlockIssuerV6 struct {
	VerificationKey string
}

type ByronProtocolV6 struct {
	Version  ProtocolVersion
	Id       uint64 // aka magic
	Software map[string]interface{}
	Update   json.RawMessage
}

// TODO - Probably a better way to handle this. For now, just split up and handle
// the differences manually. Also, we need to check the direction. Backward & forward
// are different.
type ResultByronEBBV6 struct {
	Direction string           `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Block     *ByronBlockEBBV6 `json:"block,omitempty"     dynamodbav:"block,omitempty"`
	Tip       *TipV6           `json:"tip,omitempty"       dynamodbav:"tip,omitempty"`
}

type ResultByronBFTV6 struct {
	Direction string           `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Block     *ByronBlockBFTV6 `json:"block,omitempty"     dynamodbav:"block,omitempty"`
	Tip       *TipV6           `json:"tip,omitempty"       dynamodbav:"tip,omitempty"`
}

type ResponseByronEBBV6 struct {
	JsonRpc string            `json:"jsonrpc,omitempty" dynamodbav:"jsonrpc,omitempty"`
	Method  string            `json:"method,omitempty"  dynamodbav:"method,omitempty"`
	Result  *ResultByronEBBV6 `json:"result,omitempty"  dynamodbav:"result,omitempty"`
	ID      json.RawMessage   `json:"id,omitempty"      dynamodbav:"id,omitempty"`
}

type ResponseByronBFTV6 struct {
	JsonRpc string            `json:"jsonrpc,omitempty" dynamodbav:"jsonrpc,omitempty"`
	Method  string            `json:"method,omitempty"  dynamodbav:"method,omitempty"`
	Result  *ResultByronBFTV6 `json:"result,omitempty"  dynamodbav:"result,omitempty"`
	ID      json.RawMessage   `json:"id,omitempty"      dynamodbav:"id,omitempty"`
}
