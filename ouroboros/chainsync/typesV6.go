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

package chainsync

import (
	"encoding/json"
	"strconv"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
)

// TODO - Probably a better way to handle this. For now, just split up and handle
// the differences manually. Also, we need to check the direction. Backward & forward
// are different.
type ResultPraosV6 struct {
	Direction string   `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Block     *BlockV6 `json:"block,omitempty"     dynamodbav:"block,omitempty"`
	Tip       *TipV6   `json:"tip,omitempty"       dynamodbav:"tip,omitempty"`
}

type ResponsePraosV6 struct {
	JsonRpc string          `json:"jsonrpc,omitempty" dynamodbav:"jsonrpc,omitempty"`
	Method  string          `json:"method,omitempty"  dynamodbav:"method,omitempty"`
	Result  *ResultPraosV6  `json:"result,omitempty"  dynamodbav:"result,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"      dynamodbav:"id,omitempty"`
}

type BlockV6 struct {
	Type         string        `json:"type,omitempty"`
	Era          string        `json:"era,omitempty"`
	ID           string        `json:"id,omitempty"`
	Ancestor     string        `json:"ancestor,omitempty"`
	Nonce        NonceV6       `json:"nonce,omitempty"`
	Height       uint64        `json:"height,omitempty"`
	Size         BlockSizeV6   `json:"size,omitempty"`
	Slot         uint64        `json:"slot,omitempty"`
	Transactions []TxV6        `json:"transactions,omitempty"`
	Protocol     ProtocolV6    `json:"protocol,omitempty"`
	Issuer       BlockIssuerV6 `json:"issuer,omitempty"`
}

type BlockIssuerV6 struct {
	VerificationKey        string        `json:"verificationKey,omitempty"`
	VrfVerificationKey     string        `json:"vrfVerificationKey,omitempty"`
	OperationalCertificate OpCertV6      `json:"operationalCertificate,omitempty"`
	LeaderValue            LeaderValueV6 `json:"leaderValue,omitempty"`
}

type OpCertV6 struct {
	Count uint64 `json:"count,omitempty"`
	Kes   KesV6  `json:"kes,omitempty"`
}

type LeaderValueV6 struct {
	Proof  string `json:"proof,omitempty"`
	Output string `json:"output,omitempty"`
}

type KesV6 struct {
	Period          uint64 `json:"period,omitempty"`
	VerificationKey string `json:"verificationKey,omitempty"`
}

type ProtocolV6 struct {
	Version ProtocolVersion `json:"version,omitempty" dynamodbav:"version,omitempty"`
}

type NonceV6 struct {
	Output string `json:"output,omitempty" dynamodbav:"slot,omitempty"`
	Proof  string `json:"proof,omitempty"  dynamodbav:"slot,omitempty"`
}

type BlockSizeV6 struct {
	Bytes int64
}

type TxV6 struct {
	ID                       string                `json:"id,omitempty"                       dynamodbav:"id,omitempty"`
	Spends                   string                `json:"spends,omitempty"                   dynamodbav:"spends,omitempty"`
	Inputs                   []TxInV6              `json:"inputs,omitempty"                   dynamodbav:"inputs,omitempty"`
	References               []TxInV6              `json:"references,omitempty"               dynamodbav:"references,omitempty"`
	Collaterals              []TxInV6              `json:"collaterals,omitempty"              dynamodbav:"collaterals,omitempty"`
	TotalCollateral          *int64                `json:"totalCollateral,omitempty"          dynamodbav:"totalCollateral,omitempty"`
	CollateralReturn         *TxOutV6              `json:"collateralReturn,omitempty"         dynamodbav:"collateralReturn,omitempty"`
	Outputs                  TxOutsV6              `json:"outputs,omitempty"                  dynamodbav:"outputs,omitempty"`
	Certificates             []json.RawMessage     `json:"certificates,omitempty"             dynamodbav:"certificates,omitempty"`
	Withdrawals              map[string]LovelaceV6 `json:"withdrawals,omitempty"              dynamodbav:"withdrawals,omitempty"`
	Fee                      LovelaceV6            `json:"fee,omitempty"                      dynamodbav:"fee,omitempty"`
	ValidityInterval         ValidityInterval      `json:"validityInterval"                   dynamodbav:"validityInterval,omitempty"`
	Mint                     []DoubleNestedInteger `json:"mint,omitempty"                     dynamodbav:"mint,omitempty"`
	Network                  json.RawMessage       `json:"network,omitempty"                  dynamodbav:"network,omitempty"`
	ScriptIntegrityHash      string                `json:"scriptIntegrityHash,omitempty"      dynamodbav:"scriptIntegrityHash,omitempty"`
	RequiredExtraSignatories []string              `json:"requiredExtraSignatories,omitempty" dynamodbav:"requiredExtraSignatories,omitempty"`
	RequiredExtraScripts     []string              `json:"requiredExtraScripts,omitempty"     dynamodbav:"requiredExtraScripts,omitempty"`
	Proposals                json.RawMessage       `json:"proposals,omitempty"                dynamodbav:"proposals,omitempty"`
	Votes                    json.RawMessage       `json:"votes,omitempty"                    dynamodbav:"votes,omitempty"`
	Metadata                 json.RawMessage       `json:"metadata,omitempty"                 dynamodbav:"metadata,omitempty"`
	Signatories              []json.RawMessage     `json:"signatories,omitempty"              dynamodbav:"signatories,omitempty"`
	Scripts                  json.RawMessage       `json:"scripts,omitempty"                  dynamodbav:"scripts,omitempty"`
	Datums                   Datums                `json:"datums,omitempty"                   dynamodbav:"datums,omitempty"`
	Redeemers                json.RawMessage       `json:"redeemers,omitempty"                dynamodbav:"redeemers,omitempty"`
	CBOR                     string                `json:"cbor,omitempty"                     dynamodbav:"cbor,omitempty"`
}

type TxInV6 struct {
	Transaction TxInIdV6 `json:"transaction"    dynamodbav:"transaction"`
	Index       int      `json:"index" dynamodbav:"index"`
}

type TxInIdV6 struct {
	ID string `json:"id" dynamodbav:"id"`
}

func (t TxInV6) String() string {
	return t.Transaction.ID + "#" + strconv.Itoa(t.Index)
}

type TxOutV6 struct {
	Address   string          `json:"address,omitempty"   dynamodbav:"address,omitempty"`
	Value     ValueV6         `json:"value,omitempty"     dynamodbav:"value,omitempty"`
	DatumHash string          `json:"datumHash,omitempty" dynamodbav:"datumHash,omitempty"`
	Datum     string          `json:"datum,omitempty"     dynamodbav:"datum,omitempty"`
	Script    json.RawMessage `json:"script,omitempty"    dynamodbav:"script,omitempty"`
}

type ValueV6 map[string]map[string]uint64

func (v ValueV6) Lovelace() uint64 {
	return v["ada"]["lovelace"]
}

type TxOutsV6 []TxOutV6

type LovelaceV6 struct {
	Lovelace num.Int `json:"lovelace,omitempty"  dynamodbav:"lovelace,omitempty"`
	Extras   []DoubleNestedInteger
}

type DoubleNestedInteger map[string]map[string]num.Int

type TipV6 struct {
	Slot   uint64 `json:"slot,omitempty"   dynamodbav:"slot,omitempty"`
	ID     string `json:"id,omitempty"     dynamodbav:"id,omitempty"`
	Height uint64 `json:"height,omitempty" dynamodbav:"height,omitempty"`
}
