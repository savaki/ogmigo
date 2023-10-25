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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
	"github.com/fxamacker/cbor/v2"
)

var (
	bNil = []byte("nil")
)

type IntersectionFound struct {
	Point PointV5
	Tip   TipV5
}

type IntersectionNotFound struct {
	Tip TipV5
}

// Use V5 materials only for JSON backwards compatibility.
type TxV5 struct {
	ID          string            `json:"id,omitempty"       dynamodbav:"id,omitempty"`
	InputSource string            `json:"inputSource,omitempty"  dynamodbav:"inputSource,omitempty"`
	Body        TxBodyV5          `json:"body,omitempty"     dynamodbav:"body,omitempty"`
	Witness     chainsync.Witness `json:"witness,omitempty"  dynamodbav:"witness,omitempty"`
	Metadata    json.RawMessage   `json:"metadata,omitempty" dynamodbav:"metadata,omitempty"`
	// Raw serialized transaction, base64.
	Raw string `json:"raw,omitempty" dynamodbav:"raw,omitempty"`
}

func (t TxV5) ConvertToV6() chainsync.Tx {
	withdrawals := make(map[string]chainsync.Lovelace)
	for txid, amt := range t.Body.Withdrawals {
		withdrawals[txid] = chainsync.Lovelace{Lovelace: num.Int64(amt), Extras: nil}
	}

	tx := chainsync.Tx{
		ID:                       t.ID,
		Spends:                   t.InputSource,
		Inputs:                   t.Body.Inputs,
		References:               t.Body.References,
		Collaterals:              t.Body.Collaterals,
		TotalCollateral:          t.Body.TotalCollateral,
		CollateralReturn:         (*chainsync.TxOut)(t.Body.CollateralReturn),
		Outputs:                  t.Body.Outputs,
		Certificates:             t.Body.Certificates,
		Withdrawals:              withdrawals,
		Fee:                      chainsync.Lovelace{Lovelace: t.Body.Fee, Extras: nil},
		ValidityInterval:         t.Body.ValidityInterval,
		Mint:                     nil, // TODO - Differences appear to be too much to handle.
		Network:                  t.Body.Network,
		ScriptIntegrityHash:      t.Body.ScriptIntegrityHash,
		RequiredExtraSignatories: t.Body.RequiredExtraSignatures,
		RequiredExtraScripts:     nil,
		Proposals:                t.Body.Update,
		Votes:                    nil,
		Metadata:                 t.Metadata,
		Signatories:              t.Witness.Bootstrap,
		Scripts:                  t.Witness.Scripts,
		Datums:                   t.Witness.Datums,
		Redeemers:                t.Witness.Redeemers,
		CBOR:                     t.Raw,
	}

	return tx
}

type TxBodyV5 struct {
	Certificates            []json.RawMessage          `json:"certificates,omitempty"            dynamodbav:"certificates,omitempty"`
	Collaterals             []chainsync.TxIn           `json:"collaterals,omitempty"             dynamodbav:"collaterals,omitempty"`
	Fee                     num.Int                    `json:"fee,omitempty"                     dynamodbav:"fee,omitempty"`
	Inputs                  []chainsync.TxIn           `json:"inputs,omitempty"                  dynamodbav:"inputs,omitempty"`
	Mint                    *chainsync.Value           `json:"mint,omitempty"                    dynamodbav:"mint,omitempty"`
	Network                 json.RawMessage            `json:"network,omitempty"                 dynamodbav:"network,omitempty"`
	Outputs                 chainsync.TxOuts           `json:"outputs,omitempty"                 dynamodbav:"outputs,omitempty"`
	RequiredExtraSignatures []string                   `json:"requiredExtraSignatures,omitempty" dynamodbav:"requiredExtraSignatures,omitempty"`
	ScriptIntegrityHash     string                     `json:"scriptIntegrityHash,omitempty"     dynamodbav:"scriptIntegrityHash,omitempty"`
	TimeToLive              int64                      `json:"timeToLive,omitempty"              dynamodbav:"timeToLive,omitempty"`
	Update                  json.RawMessage            `json:"update,omitempty"                  dynamodbav:"update,omitempty"`
	ValidityInterval        chainsync.ValidityInterval `json:"validityInterval"                  dynamodbav:"validityInterval,omitempty"`
	Withdrawals             map[string]int64           `json:"withdrawals,omitempty"             dynamodbav:"withdrawals,omitempty"`
	CollateralReturn        *chainsync.TxOut           `json:"collateralReturn,omitempty"        dynamodbav:"collateralReturn,omitempty"`
	TotalCollateral         *int64                     `json:"totalCollateral,omitempty"         dynamodbav:"totalCollateral,omitempty"`
	References              []chainsync.TxIn           `json:"references,omitempty"              dynamodbav:"references,omitempty"`
}

type BlockV5 struct {
	Body       []TxV5        `json:"body,omitempty"       dynamodbav:"body,omitempty"`
	Header     BlockHeaderV5 `json:"header,omitempty"     dynamodbav:"header,omitempty"`
	HeaderHash string        `json:"headerHash,omitempty" dynamodbav:"headerHash,omitempty"`
}

type BlockHeaderV5 struct {
	BlockHash       string                 `json:"blockHash,omitempty"       dynamodbav:"blockHash,omitempty"`
	BlockHeight     uint64                 `json:"blockHeight,omitempty"     dynamodbav:"blockHeight,omitempty"`
	BlockSize       uint64                 `json:"blockSize,omitempty"       dynamodbav:"blockSize,omitempty"`
	IssuerVK        string                 `json:"issuerVK,omitempty"        dynamodbav:"issuerVK,omitempty"`
	IssuerVrf       string                 `json:"issuerVrf,omitempty"       dynamodbav:"issuerVrf,omitempty"`
	LeaderValue     map[string][]byte      `json:"leaderValue,omitempty"     dynamodbav:"leaderValue,omitempty"`
	Nonce           map[string]string      `json:"nonce,omitempty"           dynamodbav:"nonce,omitempty"`
	OpCert          map[string]interface{} `json:"opCert,omitempty"          dynamodbav:"opCert,omitempty"`
	PrevHash        string                 `json:"prevHash,omitempty"        dynamodbav:"prevHash,omitempty"`
	ProtocolVersion map[string]int         `json:"protocolVersion,omitempty" dynamodbav:"protocolVersion,omitempty"`
	Signature       string                 `json:"signature,omitempty"       dynamodbav:"signature,omitempty"`
	Slot            uint64                 `json:"slot,omitempty"            dynamodbav:"slot,omitempty"`
}

// Assume no Byron support.
func (r BlockV5) PointStruct() PointStructV5 {
	return PointStructV5{
		BlockNo: r.Header.BlockHeight,
		Hash:    r.HeaderHash,
		Slot:    r.Header.Slot,
	}
}

func (b BlockV5) ConvertToV6() chainsync.Block {
	txArray := make([]chainsync.Tx, len(b.Body))
	for _, t := range b.Body {
		txArray = append(txArray, t.ConvertToV6())
	}

	fakeNonce := chainsync.Nonce{Output: "fake", Proof: "fake"}
	protocolVersion := chainsync.ProtocolVersion{
		Major: uint32(b.Header.ProtocolVersion["Major"]),
		Minor: uint32(b.Header.ProtocolVersion["Minor"]),
		Patch: uint32(b.Header.ProtocolVersion["Patch"]),
	}
	protocol := chainsync.Protocol{Version: protocolVersion}
	leaderValue := chainsync.LeaderValue{Proof: "fake", Output: "fake"}
	var opCert chainsync.OpCert
	if b.Header.OpCert != nil {
		var vk []byte
		if b.Header.OpCert["hotVk"].([]byte) != nil {
			vk, _ = base64.StdEncoding.DecodeString(b.Header.OpCert["hotVk"].(string))
		}
		opCert = chainsync.OpCert{
			Count: b.Header.OpCert["count"].(uint64),
			Kes:   chainsync.Kes{Period: b.Header.OpCert["kesPeriod"].(uint64), VerificationKey: string(vk)},
		}
	}
	issuer := chainsync.BlockIssuer{
		VerificationKey:        b.Header.IssuerVK,
		VrfVerificationKey:     b.Header.IssuerVrf,
		OperationalCertificate: opCert,
		LeaderValue:            leaderValue}
	b6 := chainsync.Block{
		Type:         "praos",
		Era:          "babbage", // TODO - Get from V5 entry - Not trivial as designed
		ID:           b.HeaderHash,
		Ancestor:     b.Header.PrevHash,
		Nonce:        fakeNonce,
		Height:       b.Header.BlockHeight,
		Size:         chainsync.BlockSize{Bytes: int64(b.Header.BlockSize)},
		Slot:         b.Header.Slot,
		Transactions: txArray,
		Protocol:     protocol,
		Issuer:       issuer,
	}

	return b6
}

type PointStructV5 struct {
	BlockNo uint64 `json:"blockNo,omitempty" dynamodbav:"blockNo,omitempty"`
	Hash    string `json:"hash,omitempty"    dynamodbav:"hash,omitempty"` // BLAKE2b_256 hash
	Slot    uint64 `json:"slot,omitempty"    dynamodbav:"slot,omitempty"`
}

func (p PointStructV5) Point() PointV5 {
	return PointV5{
		pointType:   chainsync.PointTypeStruct,
		pointStruct: &p,
	}
}

type PointV5 struct {
	pointType   chainsync.PointType
	pointString chainsync.PointString
	pointStruct *PointStructV5
}

func (p PointV5) String() string {
	switch p.pointType {
	case chainsync.PointTypeString:
		return string(p.pointString)
	case chainsync.PointTypeStruct:
		return fmt.Sprintf("slot=%v hash=%v", p.pointStruct.Slot, p.pointStruct.Hash)
	default:
		return "invalid point"
	}
}

func (p PointV5) ConvertToV6() chainsync.Point {
	var p6 chainsync.Point
	if p.pointType == chainsync.PointTypeString {
		p6 = p.pointString.Point()
	} else {
		ps := chainsync.PointStruct{Slot: p.pointStruct.Slot, ID: p.pointStruct.Hash}
		p6 = ps.Point()
	}

	return p6
}

type PointsV5 []PointV5

func (pp PointsV5) String() string {
	var ss []string
	for _, p := range pp {
		ss = append(ss, p.String())
	}
	return strings.Join(ss, ", ")
}

func (pp PointsV5) ConvertToV6() chainsync.Points {
	var points chainsync.Points
	for _, p := range pp {
		points = append(points, p.ConvertToV6())
	}
	return points
}

// pointCBOR provide simplified internal wrapper
type pointCBORV5 struct {
	String chainsync.PointString `cbor:"1,keyasint,omitempty"`
	Struct *PointStructV5        `cbor:"2,keyasint,omitempty"`
}

func (p PointV5) PointType() chainsync.PointType { return p.pointType }
func (p PointV5) PointString() (chainsync.PointString, bool) {
	return p.pointString, p.pointString != ""
}

func (p PointV5) PointStruct() (*PointStructV5, bool) { return p.pointStruct, p.pointStruct != nil }

func (p PointV5) MarshalCBOR() ([]byte, error) {
	switch p.pointType {
	case chainsync.PointTypeString, chainsync.PointTypeStruct:
		v := pointCBORV5{
			String: p.pointString,
			Struct: p.pointStruct,
		}
		return cbor.Marshal(v)
	default:
		return nil, fmt.Errorf("unable to unmarshal Point: unknown type")
	}
}

func (p PointV5) MarshalJSON() ([]byte, error) {
	switch p.pointType {
	case chainsync.PointTypeString:
		return json.Marshal(p.pointString)
	case chainsync.PointTypeStruct:
		return json.Marshal(p.pointStruct)
	default:
		return nil, fmt.Errorf("unable to unmarshal Point: unknown type")
	}
}

func (p *PointV5) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, bNil) {
		return nil
	}

	var v pointCBORV5
	if err := cbor.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("failed to unmarshal Point: %w", err)
	}

	point := PointV5{
		pointType:   chainsync.PointTypeString,
		pointString: v.String,
		pointStruct: v.Struct,
	}
	if point.pointStruct != nil {
		point.pointType = chainsync.PointTypeStruct
	}

	*p = point

	return nil
}

func (p *PointV5) UnmarshalJSON(data []byte) error {
	switch {
	case data[0] == '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("failed to unmarshal Point, %v: %w", string(data), err)
		}

		*p = PointV5{
			pointType:   chainsync.PointTypeString,
			pointString: chainsync.PointString(s),
		}

	default:
		var ps PointStructV5
		if err := json.Unmarshal(data, &ps); err != nil {
			return fmt.Errorf("failed to unmarshal Point, %v: %w", string(data), err)
		}

		*p = PointV5{
			pointType:   chainsync.PointTypeStruct,
			pointStruct: &ps,
		}
	}

	return nil
}

type ResultV5 struct {
	IntersectionFound    *IntersectionFound    `json:",omitempty" dynamodbav:",omitempty"`
	IntersectionNotFound *IntersectionNotFound `json:",omitempty" dynamodbav:",omitempty"`
	RollForward          *RollForwardV5        `json:",omitempty" dynamodbav:",omitempty"`
	RollBackward         *RollBackwardV5       `json:",omitempty" dynamodbav:",omitempty"`
}

type ResultFindIntersectionV5 struct {
	IntersectionFound    *IntersectionFoundV5    `json:",omitempty" dynamodbav:",omitempty"`
	IntersectionNotFound *IntersectionNotFoundV5 `json:",omitempty" dynamodbav:",omitempty"`
}

type RollBackwardV5 struct {
	Point PointV5 `json:"point,omitempty" dynamodbav:"point,omitempty"`
	Tip   TipV5   `json:"tip,omitempty"   dynamodbav:"tip,omitempty"`
}

type RollForwardV5 struct {
	Block BlockV5 `json:"block,omitempty" dynamodbav:"block,omitempty"`
	Tip   TipV5   `json:"tip,omitempty"   dynamodbav:"tip,omitempty"`
}

type ResultNextBlockV5 struct {
	RollForward  *RollForwardV5  `json:",omitempty" dynamodbav:",omitempty"`
	RollBackward *RollBackwardV5 `json:",omitempty" dynamodbav:",omitempty"`
}

type IntersectionFoundV5 struct {
	Point *chainsync.Point
	Tip   *TipV5
}

type IntersectionNotFoundV5 struct {
	Tip *TipV5
}

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResultFindIntersection chainsync.ResultFindIntersectionPraos

func (c *CompatibleResultFindIntersection) UnmarshalJSON(data []byte) error {
	// Assume v6 responses first, then fall back to manual v5 processing.
	var r chainsync.ResultFindIntersectionPraos
	err := json.Unmarshal(data, &r)
	if err == nil && r.Tip != nil {
		*c = CompatibleResultFindIntersection(r)
		return nil
	}

	var r5 ResultFindIntersectionV5
	err = json.Unmarshal(data, &r5)
	if err != nil {
		return err
	} else if r5.IntersectionFound != nil {
		tip := r5.IntersectionFound.Tip.ConvertToV6()
		c.Intersection = r5.IntersectionFound.Point
		c.Tip = &tip
		c.Error = nil
		c.ID = nil

		return nil
	} else if r5.IntersectionNotFound != nil {
		// Emulate the v6 IntersectionNotFound error as best as possible.
		tip := r5.IntersectionNotFound.Tip.ConvertToV6()
		err := chainsync.ResultError{Code: 1000, Message: "Intersection not found", Data: &tip}
		c.Error = &err

		return nil
	}

	// TODO: Further error handling here.
	return nil
}

func (c CompatibleResultFindIntersection) String() string {
	return fmt.Sprintf("intersection=[%v] tip=[%v] error=[%v] id=[%v]", c.Intersection, c.Tip, c.Error, c.ID)
}

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResultNextBlock chainsync.ResultNextBlockPraos

func (c *CompatibleResultNextBlock) UnmarshalJSON(data []byte) error {
	// Assume v6 responses first, then fall back to manual v5 processing.
	var r chainsync.ResultNextBlockPraos
	err := json.Unmarshal(data, &r)
	if err == nil && r.Tip != nil {
		*c = CompatibleResultNextBlock(r)
		return nil
	}

	var r5 ResultNextBlockV5
	err = json.Unmarshal(data, &r5)
	if err != nil {
		return err
	} else if r5.RollForward != nil {
		tip := r5.RollForward.Tip.ConvertToV6()
		block := r5.RollForward.Block.ConvertToV6()
		c.Direction = chainsync.RollForwardString
		c.Tip = &tip
		c.Block = &block

		return nil
	} else if r5.RollBackward != nil {
		tip := r5.RollBackward.Tip.ConvertToV6()
		point := r5.RollBackward.Point.ConvertToV6()
		c.Direction = chainsync.RollBackwardString
		c.Tip = &tip
		c.Point = &point

		return nil
	}

	// TODO: Further error handling here.
	return nil
}

func (c CompatibleResultNextBlock) String() string {
	return fmt.Sprintf("direction=[%v] tip=[%v] block=[%v] point=[%v]", c.Direction, c.Tip, c.Block, c.Point)
}

type TipV5 struct {
	Slot    uint64 `json:"slot,omitempty"    dynamodbav:"slot,omitempty"`
	Hash    string `json:"hash,omitempty"    dynamodbav:"hash,omitempty"`
	BlockNo uint64 `json:"blockNo,omitempty" dynamodbav:"blockNo,omitempty"`
}

func (t TipV5) ConvertToV6() chainsync.Tip {
	return chainsync.Tip{
		Slot:   t.Slot,
		ID:     t.Hash,
		Height: t.BlockNo,
	}
}

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResponsePraos chainsync.ResponsePraos

func (c *CompatibleResponsePraos) UnmarshalJSON(data []byte) error {
	var r chainsync.ResponsePraos
	err := json.Unmarshal(data, &r)
	if err == nil && r.Result != nil {
		*c = CompatibleResponsePraos(r)
		return nil
	}

	var r5 ResponseV5
	err = json.Unmarshal(data, &r5)
	c.JsonRpc = "2.0"
	if err != nil {
		// Just skip all the data processing, as it's useless.
		return err
	}

	var p chainsync.Point
	var t chainsync.Tip
	var e chainsync.ResultError

	// All we really care about is the result.
	if r5.Result.IntersectionFound != nil {
		c.Method = chainsync.FindIntersectionMethod

		pt := r5.Result.IntersectionFound.Point.pointType
		if pt == chainsync.PointTypeString {
			p = r5.Result.IntersectionFound.Point.pointString.Point()
		} else if pt == chainsync.PointTypeStruct {
			ps := chainsync.PointStruct{Slot: r5.Result.IntersectionFound.Point.pointStruct.Slot, ID: r5.Result.IntersectionFound.Point.pointStruct.Hash}
			p = ps.Point()
		} else {
			panic("Invalid point type")
		}

		t.Slot = r5.Result.IntersectionFound.Tip.Slot
		t.ID = r5.Result.IntersectionFound.Tip.Hash
		t.Height = r5.Result.IntersectionFound.Tip.BlockNo

		var findIntersection CompatibleResultFindIntersection
		findIntersection.Intersection = &p
		findIntersection.Tip = &t
		c.Result = &findIntersection
	} else if r5.Result.IntersectionNotFound != nil {
		c.Method = chainsync.FindIntersectionMethod

		t.Slot = r5.Result.IntersectionNotFound.Tip.Slot
		t.ID = r5.Result.IntersectionNotFound.Tip.Hash
		t.Height = r5.Result.IntersectionFound.Tip.BlockNo

		e.Data = &t
		e.Code = 1000
		e.Message = "Intersection not found - Conversion from a v5 Ogmigo call"
		c.Error = &e
	} else if r5.Result.RollForward != nil {
		c.Method = chainsync.NextBlockMethod

		t.Slot = r5.Result.RollForward.Tip.Slot
		t.ID = r5.Result.RollForward.Tip.Hash
		t.Height = r5.Result.RollForward.Tip.BlockNo

		txArray := make([]chainsync.Tx, len(r5.Result.RollForward.Block.Body))
		for _, t := range r5.Result.RollForward.Block.Body {
			withdrawals := make(map[string]chainsync.Lovelace)
			for txid, amt := range t.Body.Withdrawals {
				withdrawals[txid] = chainsync.Lovelace{Lovelace: num.Int64(amt), Extras: nil}
			}

			tx := chainsync.Tx{
				ID:                       t.ID,
				Spends:                   t.InputSource,
				Inputs:                   t.Body.Inputs,
				References:               t.Body.References,
				Collaterals:              t.Body.Collaterals,
				TotalCollateral:          t.Body.TotalCollateral,
				CollateralReturn:         (*chainsync.TxOut)(t.Body.CollateralReturn),
				Outputs:                  t.Body.Outputs,
				Certificates:             t.Body.Certificates,
				Withdrawals:              withdrawals,
				Fee:                      chainsync.Lovelace{Lovelace: t.Body.Fee, Extras: nil},
				ValidityInterval:         t.Body.ValidityInterval,
				Mint:                     nil, // TODO - Differences appear to be too much to handle.
				Network:                  t.Body.Network,
				ScriptIntegrityHash:      t.Body.ScriptIntegrityHash,
				RequiredExtraSignatories: t.Body.RequiredExtraSignatures,
				RequiredExtraScripts:     nil,
				Proposals:                t.Body.Update,
				Votes:                    nil,
				Metadata:                 t.Metadata,
				Signatories:              t.Witness.Bootstrap,
				Scripts:                  t.Witness.Scripts,
				Datums:                   t.Witness.Datums,
				Redeemers:                t.Witness.Redeemers,
				CBOR:                     t.Raw,
			}
			txArray = append(txArray, tx)
		}

		fakeNonce := chainsync.Nonce{Output: "fake", Proof: "fake"}
		protocolVersion := chainsync.ProtocolVersion{
			Major: uint32(r5.Result.RollForward.Block.Header.ProtocolVersion["Major"]),
			Minor: uint32(r5.Result.RollForward.Block.Header.ProtocolVersion["Minor"]),
			Patch: uint32(r5.Result.RollForward.Block.Header.ProtocolVersion["Patch"]),
		}
		protocol := chainsync.Protocol{Version: protocolVersion}
		leaderValue := chainsync.LeaderValue{Proof: "fake", Output: "fake"}
		var opCert chainsync.OpCert
		if r5.Result.RollForward.Block.Header.OpCert != nil {
			var vk []byte
			if r5.Result.RollForward.Block.Header.OpCert["hotVk"].([]byte) != nil {
				vk, _ = base64.StdEncoding.DecodeString(r5.Result.RollForward.Block.Header.OpCert["hotVk"].(string))
			}
			opCert = chainsync.OpCert{
				Count: r5.Result.RollForward.Block.Header.OpCert["count"].(uint64),
				Kes:   chainsync.Kes{Period: r5.Result.RollForward.Block.Header.OpCert["kesPeriod"].(uint64), VerificationKey: string(vk)},
			}
		}
		issuer := chainsync.BlockIssuer{VerificationKey: r5.Result.RollForward.Block.Header.IssuerVK, VrfVerificationKey: r5.Result.RollForward.Block.Header.IssuerVrf, OperationalCertificate: opCert, LeaderValue: leaderValue}
		b := chainsync.Block{
			Type:         "praos",
			Era:          "babbage", // TODO - Get from V5 entry - Not trivial as designed
			ID:           r5.Result.RollForward.Block.HeaderHash,
			Ancestor:     r5.Result.RollForward.Block.Header.PrevHash,
			Nonce:        fakeNonce,
			Height:       r5.Result.RollForward.Block.Header.BlockHeight,
			Size:         chainsync.BlockSize{Bytes: int64(r5.Result.RollForward.Block.Header.BlockSize)},
			Slot:         r5.Result.RollForward.Block.Header.Slot,
			Transactions: txArray,
			Protocol:     protocol,
			Issuer:       issuer,
		}

		var nextBlock CompatibleResultNextBlock
		nextBlock.Direction = chainsync.RollForwardString
		nextBlock.Tip = &t
		nextBlock.Block = &b
		c.Result = &nextBlock
	} else if r5.Result.RollBackward != nil {
		c.Method = chainsync.NextBlockMethod
		var t chainsync.Tip
		t.Slot = r5.Result.RollBackward.Tip.Slot
		t.ID = r5.Result.RollBackward.Tip.Hash
		t.Height = r5.Result.RollBackward.Tip.BlockNo

		var p chainsync.Point
		pt := r5.Result.RollBackward.Point.pointType
		if pt == chainsync.PointTypeString {
			p = r5.Result.RollBackward.Point.pointString.Point()
		} else {
			ps := chainsync.PointStruct{Slot: r5.Result.RollBackward.Point.pointStruct.Slot, ID: r5.Result.RollBackward.Point.pointStruct.Hash}
			p = ps.Point()
		}

		var nextBlock CompatibleResultNextBlock
		nextBlock.Direction = chainsync.RollBackwardString
		nextBlock.Tip = &t
		nextBlock.Point = &p
		c.Result = &nextBlock
	}

	c.ID = r5.Reflection
	return nil
}

func (r CompatibleResponsePraos) MustFindIntersectResult() CompatibleResultFindIntersection {
	if r.Method != chainsync.FindIntersectionMethod {
		panic(fmt.Errorf("must only use *Must* methods after switching on the findIntersection method; called on %v", r.Method))
	}
	return r.Result.(CompatibleResultFindIntersection)
}

func (r CompatibleResponsePraos) MustNextBlockResult() CompatibleResultNextBlock {
	if r.Method != chainsync.NextBlockMethod {
		panic(fmt.Errorf("must only use *Must* methods after switching on the nextBlock method; called on %v", r.Method))
	}
	fmt.Printf("type of r.Result is %T\n", r.Result)
	return r.Result.(CompatibleResultNextBlock)
}

type ResponseV5 struct {
	Type        string          `json:"type,omitempty"        dynamodbav:"type,omitempty"`
	Version     string          `json:"version,omitempty"     dynamodbav:"version,omitempty"`
	ServiceName string          `json:"servicename,omitempty" dynamodbav:"servicename,omitempty"`
	MethodName  string          `json:"methodname,omitempty"  dynamodbav:"methodname,omitempty"`
	Result      *ResultV5       `json:"result,omitempty"      dynamodbav:"result,omitempty"`
	Reflection  json.RawMessage `json:"reflection,omitempty"  dynamodbav:"reflection,omitempty"`
}
