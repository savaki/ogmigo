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
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fxamacker/cbor/v2"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
)

var (
	encOptions = cbor.CoreDetEncOptions()
	bNil       = []byte("nil")
)

type AssetID string

func (a AssetID) HasPolicyID(s string) bool {
	return len(s) == 56 && strings.HasPrefix(string(a), s)
}

func (a AssetID) HasAssetID(re *regexp.Regexp) bool {
	return re.MatchString(string(a))
}

func (a AssetID) IsZero() bool {
	return a == ""
}

func (a AssetID) MatchAssetName(re *regexp.Regexp) ([]string, bool) {
	if assetName := a.AssetName(); len(assetName) > 0 {
		ss := re.FindStringSubmatch(assetName)
		return ss, len(ss) > 0
	}
	return nil, false
}

func (a AssetID) String() string {
	return string(a)
}

func (a AssetID) AssetName() string {
	s := string(a)
	if index := strings.Index(s, "."); index > 0 {
		return s[index+1:]
	}
	return ""
}

func (a AssetID) AssetNameUTF8() (string, bool) {
	if data, err := hex.DecodeString(a.AssetName()); err == nil {
		if utf8.Valid(data) {
			return string(data), true
		}
	}
	return "", false
}

func (a AssetID) PolicyID() string {
	s := string(a)
	if index := strings.Index(s, "."); index > 0 {
		return s[:index]
	}
	return s // Assets with empty-string name come back as just the policy ID
}

type IntersectionFound struct {
	Point PointV5
	Tip   TipV5
}

type IntersectionNotFound struct {
	Tip TipV5
}

// Use V5 materials only for JSON backwards compatibility.
type TxV5 struct {
	ID          string          `json:"id,omitempty"       dynamodbav:"id,omitempty"`
	InputSource string          `json:"inputSource,omitempty"  dynamodbav:"inputSource,omitempty"`
	Body        TxBodyV5        `json:"body,omitempty"     dynamodbav:"body,omitempty"`
	Witness     Witness         `json:"witness,omitempty"  dynamodbav:"witness,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty" dynamodbav:"metadata,omitempty"`
	// Raw serialized transaction, base64.
	Raw string `json:"raw,omitempty" dynamodbav:"raw,omitempty"`
}

type TxBodyV5 struct {
	Certificates            []json.RawMessage `json:"certificates,omitempty"            dynamodbav:"certificates,omitempty"`
	Collaterals             []TxIn            `json:"collaterals,omitempty"             dynamodbav:"collaterals,omitempty"`
	Fee                     num.Int           `json:"fee,omitempty"                     dynamodbav:"fee,omitempty"`
	Inputs                  []TxIn            `json:"inputs,omitempty"                  dynamodbav:"inputs,omitempty"`
	Mint                    *Value            `json:"mint,omitempty"                    dynamodbav:"mint,omitempty"`
	Network                 json.RawMessage   `json:"network,omitempty"                 dynamodbav:"network,omitempty"`
	Outputs                 TxOuts            `json:"outputs,omitempty"                 dynamodbav:"outputs,omitempty"`
	RequiredExtraSignatures []string          `json:"requiredExtraSignatures,omitempty" dynamodbav:"requiredExtraSignatures,omitempty"`
	ScriptIntegrityHash     string            `json:"scriptIntegrityHash,omitempty"     dynamodbav:"scriptIntegrityHash,omitempty"`
	TimeToLive              int64             `json:"timeToLive,omitempty"              dynamodbav:"timeToLive,omitempty"`
	Update                  json.RawMessage   `json:"update,omitempty"                  dynamodbav:"update,omitempty"`
	ValidityInterval        ValidityInterval  `json:"validityInterval"                  dynamodbav:"validityInterval,omitempty"`
	Withdrawals             map[string]int64  `json:"withdrawals,omitempty"             dynamodbav:"withdrawals,omitempty"`
	CollateralReturn        *TxOut            `json:"collateralReturn,omitempty"        dynamodbav:"collateralReturn,omitempty"`
	TotalCollateral         *int64            `json:"totalCollateral,omitempty"         dynamodbav:"totalCollateral,omitempty"`
	References              []TxIn            `json:"references,omitempty"              dynamodbav:"references,omitempty"`
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

// All blocks except Byron-era blocks.
type Block struct {
	Type         string      `json:"type,omitempty"`
	Era          string      `json:"era,omitempty"`
	ID           string      `json:"id,omitempty"`
	Ancestor     string      `json:"ancestor,omitempty"`
	Nonce        Nonce       `json:"nonce,omitempty"`
	Height       uint64      `json:"height,omitempty"`
	Size         BlockSize   `json:"size,omitempty"`
	Slot         uint64      `json:"slot,omitempty"`
	Transactions []Tx        `json:"transactions,omitempty"`
	Protocol     Protocol    `json:"protocol,omitempty"`
	Issuer       BlockIssuer `json:"issuer,omitempty"`
}

type Nonce struct {
	Output string `json:"output,omitempty" dynamodbav:"slot,omitempty"`
	Proof  string `json:"proof,omitempty"  dynamodbav:"slot,omitempty"`
}

type BlockSize struct {
	Bytes int64
}

type Protocol struct {
	Version ProtocolVersion `json:"version,omitempty" dynamodbav:"version,omitempty"`
}

type BlockIssuer struct {
	VerificationKey        string      `json:"verificationKey,omitempty"`
	VrfVerificationKey     string      `json:"vrfVerificationKey,omitempty"`
	OperationalCertificate OpCert      `json:"operationalCertificate,omitempty"`
	LeaderValue            LeaderValue `json:"leaderValue,omitempty"`
}

type OpCert struct {
	Count uint64 `json:"count,omitempty"`
	Kes   Kes    `json:"kes,omitempty"`
}

type Kes struct {
	Period          uint64 `json:"period,omitempty"`
	VerificationKey string `json:"verificationKey,omitempty"`
}

type LeaderValue struct {
	Proof  string `json:"proof,omitempty"`
	Output string `json:"output,omitempty"`
}

type PointType int

const (
	PointTypeString PointType = 1
	PointTypeStruct PointType = 2
)

var Origin = PointString("origin").Point()

type PointString string

func (p PointString) Point() Point {
	return Point{
		pointType:   PointTypeString,
		pointString: p,
	}
}

type PointStruct struct {
	BlockNo uint64 `json:"blockNo,omitempty" dynamodbav:"blockNo,omitempty"` // Not part of RollBackward.
	ID      string `json:"id,omitempty"      dynamodbav:"id,omitempty"`      // BLAKE2b_256 hash
	Slot    uint64 `json:"slot,omitempty"    dynamodbav:"slot,omitempty"`
}

type PointStructV5 struct {
	Hash string `json:"hash,omitempty"    dynamodbav:"hash,omitempty"` // BLAKE2b_256 hash
	Slot uint64 `json:"slot,omitempty"    dynamodbav:"slot,omitempty"`
}

func (p PointStruct) Point() Point {
	return Point{
		pointType:   PointTypeStruct,
		pointStruct: &p,
	}
}

type Point struct {
	pointType   PointType
	pointString PointString
	pointStruct *PointStruct
}

func (p Point) String() string {
	switch p.pointType {
	case PointTypeString:
		return string(p.pointString)
	case PointTypeStruct:
		if p.pointStruct.BlockNo == 0 {
			return fmt.Sprintf("slot=%v id=%v", p.pointStruct.Slot, p.pointStruct.ID)
		}
		return fmt.Sprintf("slot=%v id=%v block=%v", p.pointStruct.Slot, p.pointStruct.ID, p.pointStruct.BlockNo)
	default:
		return "invalid point"
	}
}

type Points []Point

type PointV5 struct {
	pointType   PointType
	pointString PointString
	pointStruct *PointStructV5
}

func (p PointV5) String() string {
	switch p.pointType {
	case PointTypeString:
		return string(p.pointString)
	case PointTypeStruct:
		return fmt.Sprintf("slot=%v hash=%v block=%v", p.pointStruct.Slot, p.pointStruct.Hash)
	default:
		return "invalid point"
	}
}

type PointsV5 []Point

func (pp PointsV5) String() string {
	var ss []string
	for _, p := range pp {
		ss = append(ss, p.String())
	}
	return strings.Join(ss, ", ")
}

// pointCBOR provide simplified internal wrapper
type pointCBORV5 struct {
	String PointString    `cbor:"1,keyasint,omitempty"`
	Struct *PointStructV5 `cbor:"2,keyasint,omitempty"`
}

func (p PointV5) PointType() PointType             { return p.pointType }
func (p PointV5) PointString() (PointString, bool) { return p.pointString, p.pointString != "" }

func (p PointV5) PointStruct() (*PointStructV5, bool) { return p.pointStruct, p.pointStruct != nil }

func (p PointV5) MarshalCBOR() ([]byte, error) {
	switch p.pointType {
	case PointTypeString, PointTypeStruct:
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
	case PointTypeString:
		return json.Marshal(p.pointString)
	case PointTypeStruct:
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
		pointType:   PointTypeString,
		pointString: v.String,
		pointStruct: v.Struct,
	}
	if point.pointStruct != nil {
		point.pointType = PointTypeStruct
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
			pointType:   PointTypeString,
			pointString: PointString(s),
		}

	default:
		var ps PointStructV5
		if err := json.Unmarshal(data, &ps); err != nil {
			return fmt.Errorf("failed to unmarshal Point, %v: %w", string(data), err)
		}

		*p = PointV5{
			pointType:   PointTypeStruct,
			pointStruct: &ps,
		}
	}

	return nil
}

func (pp Points) String() string {
	var ss []string
	for _, p := range pp {
		ss = append(ss, p.String())
	}
	return strings.Join(ss, ", ")
}

func (pp Points) Len() int      { return len(pp) }
func (pp Points) Swap(i, j int) { pp[i], pp[j] = pp[j], pp[i] }
func (pp Points) Less(i, j int) bool {
	pi, pj := pp[i], pp[j]
	switch {
	case pi.pointType == PointTypeStruct && pj.pointType == PointTypeStruct:
		return pi.pointStruct.Slot > pj.pointStruct.Slot
	case pi.pointType == PointTypeStruct:
		return true
	case pj.pointType == PointTypeStruct:
		return false
	default:
		return pi.pointString > pj.pointString
	}
}

// pointCBOR provide simplified internal wrapper
type pointCBOR struct {
	String PointString  `cbor:"1,keyasint,omitempty"`
	Struct *PointStruct `cbor:"2,keyasint,omitempty"`
}

func (p Point) PointType() PointType             { return p.pointType }
func (p Point) PointString() (PointString, bool) { return p.pointString, p.pointString != "" }

func (p Point) PointStruct() (*PointStruct, bool) { return p.pointStruct, p.pointStruct != nil }

func (p Point) MarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	switch p.pointType {
	case PointTypeString:
		item.S = aws.String(string(p.pointString))
	case PointTypeStruct:
		m, err := dynamodbattribute.MarshalMap(p.pointStruct)
		if err != nil {
			return fmt.Errorf("failed to marshal point struct: %w", err)
		}
		item.M = m
	default:
		return fmt.Errorf("unable to unmarshal Point: unknown type")
	}
	return nil
}

func (p Point) MarshalCBOR() ([]byte, error) {
	switch p.pointType {
	case PointTypeString, PointTypeStruct:
		v := pointCBOR{
			String: p.pointString,
			Struct: p.pointStruct,
		}
		return cbor.Marshal(v)
	default:
		return nil, fmt.Errorf("unable to unmarshal Point: unknown type")
	}
}

func (p Point) MarshalJSON() ([]byte, error) {
	switch p.pointType {
	case PointTypeString:
		return json.Marshal(p.pointString)
	case PointTypeStruct:
		return json.Marshal(p.pointStruct)
	default:
		return nil, fmt.Errorf("unable to unmarshal Point: unknown type")
	}
}

func (p *Point) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, bNil) {
		return nil
	}

	var v pointCBOR
	if err := cbor.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("failed to unmarshal Point: %w", err)
	}

	point := Point{
		pointType:   PointTypeString,
		pointString: v.String,
		pointStruct: v.Struct,
	}
	if point.pointStruct != nil {
		point.pointType = PointTypeStruct
	}

	*p = point

	return nil
}

func (p *Point) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	switch {
	case item == nil:
		return nil
	case item.S != nil:
		*p = Point{
			pointType:   PointTypeString,
			pointString: PointString(aws.StringValue(item.S)),
		}
	case len(item.M) > 0:
		var point PointStruct
		if err := dynamodbattribute.UnmarshalMap(item.M, &point); err != nil {
			return fmt.Errorf("failed to unmarshal point struct: %w", err)
		}
		*p = Point{
			pointType:   PointTypeStruct,
			pointStruct: &point,
		}
	}
	return nil
}

func (p *Point) UnmarshalJSON(data []byte) error {
	switch {
	case data[0] == '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("failed to unmarshal Point, %v: %w", string(data), err)
		}

		*p = Point{
			pointType:   PointTypeString,
			pointString: PointString(s),
		}

	default:
		var ps PointStruct
		if err := json.Unmarshal(data, &ps); err != nil {
			return fmt.Errorf("failed to unmarshal Point, %v: %w", string(data), err)
		}

		*p = Point{
			pointType:   PointTypeStruct,
			pointStruct: &ps,
		}
	}

	return nil
}

type ProtocolVersion struct {
	Major uint32
	Minor uint32
	Patch uint32 `json:"patch,omitempty"`
}

type RollBackward struct {
	Direction string            `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Tip       Tip               `json:"tip,omitempty"   dynamodbav:"tip,omitempty"`
	Point     RollBackwardPoint `json:"point,omitempty" dynamodbav:"point,omitempty"`
}

type RollBackwardPoint struct {
	Slot uint64 `json:"slot,omitempty"    dynamodbav:"slot,omitempty"`
	ID   string `json:"id,omitempty"      dynamodbav:"id,omitempty"` // BLAKE2b_256 hash
}

// Assume non-Byron blocks.
type RollForward struct {
	Direction string `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Tip       Tip    `json:"tip,omitempty"   dynamodbav:"tip,omitempty"`
	Block     Block  `json:"block,omitempty" dynamodbav:"block,omitempty"`
}

func (r Block) PointStruct() PointStruct {
	return PointStruct{
		BlockNo: r.Height,
		ID:      r.ID,
		Slot:    r.Slot,
	}
}

type ResultV5 struct {
	IntersectionFound    *IntersectionFound    `json:",omitempty" dynamodbav:",omitempty"`
	IntersectionNotFound *IntersectionNotFound `json:",omitempty" dynamodbav:",omitempty"`
	RollForward          *RollForwardV5        `json:",omitempty" dynamodbav:",omitempty"`
	RollBackward         *RollBackwardV5       `json:",omitempty" dynamodbav:",omitempty"`
}

// Covers everything except Byron-era blocks.
type ResultFindIntersectionPraos struct {
	Intersection *Point          `json:"intersection,omitempty" dynamodbav:"intersection,omitempty"`
	Tip          *Tip            `json:"tip,omitempty"          dynamodbav:"tip,omitempty"`
	Error        *ResultError    `json:"error,omitempty"        dynamodbav:"error,omitempty"`
	ID           json.RawMessage `json:"id,omitempty"           dynamodbav:"id,omitempty"`
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
	Point *Point
	Tip   *TipV5
}

type IntersectionNotFoundV5 struct {
	Tip *TipV5
}

type ResultError struct {
	Code    uint32          `json:"code,omitempty"    dynamodbav:"code,omitempty"`
	Message string          `json:"message,omitempty" dynamodbav:"message,omitempty"`
	Data    *Tip            `json:"data,omitempty"    dynamodbav:"data,omitempty"` // Forward
	ID      json.RawMessage `json:"id,omitempty"      dynamodbav:"id,omitempty"`
}

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResultFindIntersection ResultFindIntersectionPraos

func (c *CompatibleResultFindIntersection) UnmarshalJSON(data []byte) error {
	// Assume v6 responses first, then fall back to manual v5 processing.
	var r ResultFindIntersectionPraos
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
		tip := Tip{
			Height: r5.IntersectionFound.Tip.BlockNo,
			ID:     r5.IntersectionFound.Tip.Hash,
			Slot:   r5.IntersectionFound.Tip.Slot,
		}
		c.Intersection = r5.IntersectionFound.Point
		c.Tip = &tip
		c.Error = nil
		c.ID = nil
		return nil
	} else if r5.IntersectionNotFound != nil {
		// Emulate the v6 IntersectionNotFound error as best as possible.
		tip := Tip{
			Height: r5.IntersectionFound.Tip.BlockNo,
			ID:     r5.IntersectionFound.Tip.Hash,
			Slot:   r5.IntersectionFound.Tip.Slot,
		}
		err := ResultError{Code: 1000, Message: "Intersection not found", Data: &tip}
		c.Error = &err
		return nil
	}

	// TODO: Further error handling here.
	return nil
}

func (c CompatibleResultFindIntersection) String() string {
	return fmt.Sprintf("intersection=[%v] tip=[%v] error=[%v] id=[%v]", c.Intersection, c.Tip, c.Error, c.ID)
}

// Covers all blocks except Byron-era blocks.
type ResultNextBlockPraos struct {
	Direction string `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Tip       *Tip   `json:"tip,omitempty"       dynamodbav:"tip,omitempty"`
	Block     *Block `json:"block,omitempty"     dynamodbav:"block,omitempty"` // Forward
	Point     *Point `json:"point,omitempty"     dynamodbav:"point,omitempty"` // Backward
}

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResultNextBlock ResultNextBlockPraos

func (c *CompatibleResultNextBlock) UnmarshalJSON(data []byte) error {
	// Assume v6 responses first, then fall back to manual v5 processing.
	var r ResultNextBlockPraos
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
		tip := Tip{Height: r5.RollForward.Tip.BlockNo, ID: r5.RollForward.Tip.Hash, Slot: r5.RollForward.Tip.Slot}

		txArray := make([]Tx, len(r5.RollForward.Block.Body))
		for _, t := range r5.RollForward.Block.Body {
			withdrawals := make(map[string]Lovelace)
			for txid, amt := range t.Body.Withdrawals {
				withdrawals[txid] = Lovelace{Lovelace: num.Int64(amt), Extras: nil}
			}

			tx := Tx{
				ID:                       t.ID,
				Spends:                   t.InputSource,
				Inputs:                   t.Body.Inputs,
				References:               t.Body.References,
				Collaterals:              t.Body.Collaterals,
				TotalCollateral:          t.Body.TotalCollateral,
				CollateralReturn:         (*TxOut)(t.Body.CollateralReturn),
				Outputs:                  t.Body.Outputs,
				Certificates:             t.Body.Certificates,
				Withdrawals:              withdrawals,
				Fee:                      Lovelace{Lovelace: t.Body.Fee, Extras: nil},
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

		fakeNonce := Nonce{Output: "fake", Proof: "fake"}
		protocolVersion := ProtocolVersion{
			Major: uint32(r5.RollForward.Block.Header.ProtocolVersion["Major"]),
			Minor: uint32(r5.RollForward.Block.Header.ProtocolVersion["Minor"]),
			Patch: uint32(r5.RollForward.Block.Header.ProtocolVersion["Patch"]),
		}
		protocol := Protocol{Version: protocolVersion}
		leaderValue := LeaderValue{Proof: "fake", Output: "fake"}
		vk, _ := base64.StdEncoding.DecodeString(r5.RollForward.Block.Header.OpCert["hotVk"].(string))
		opCert := OpCert{
			Count: r5.RollForward.Block.Header.OpCert["count"].(uint64),
			Kes:   Kes{Period: r5.RollForward.Block.Header.OpCert["kesPeriod"].(uint64), VerificationKey: string(vk)},
		}
		issuer := BlockIssuer{VerificationKey: r5.RollForward.Block.Header.IssuerVK, VrfVerificationKey: r5.RollForward.Block.Header.IssuerVrf, OperationalCertificate: opCert, LeaderValue: leaderValue}
		block := Block{
			Type:         "praos",
			Era:          "babbage", // TODO - Get from V5 entry - Not trivial as designed
			ID:           r5.RollForward.Block.HeaderHash,
			Ancestor:     r5.RollForward.Block.Header.PrevHash,
			Nonce:        fakeNonce,
			Height:       r5.RollForward.Block.Header.BlockHeight,
			Size:         BlockSize{Bytes: int64(r5.RollForward.Block.Header.BlockSize)},
			Slot:         r5.RollForward.Block.Header.Slot,
			Transactions: txArray,
			Protocol:     protocol,
			Issuer:       issuer,
		}
		c.Direction = "forward"
		c.Tip = &tip
		c.Block = &block
		c.Point = nil

		return nil
	} else if r5.RollBackward != nil {
		tip := Tip{Height: r5.RollForward.Tip.BlockNo, ID: r5.RollForward.Tip.Hash, Slot: r5.RollForward.Tip.Slot}
		c.Direction = "backward"
		c.Tip = &tip
		c.Block = nil
		c.Point = nil // TODO

		return nil
	}

	// TODO: Further error handling here.
	return nil
}

func (c CompatibleResultNextBlock) String() string {
	return fmt.Sprintf("direction=[%v] tip=[%v] block=[%v] point=[%v]", c.Direction, c.Tip, c.Block, c.Point)
}

type Tip struct {
	Slot   uint64 `json:"slot,omitempty"   dynamodbav:"slot,omitempty"`
	ID     string `json:"id,omitempty"     dynamodbav:"id,omitempty"`
	Height uint64 `json:"height,omitempty" dynamodbav:"height,omitempty"`
}

func (t Tip) String() string {
	return fmt.Sprintf("slot=%v id=%v height=%v", t.Slot, t.ID, t.Height)
}

type TipV5 struct {
	Slot    uint64 `json:"slot,omitempty"    dynamodbav:"slot,omitempty"`
	Hash    string `json:"hash,omitempty"    dynamodbav:"hash,omitempty"`
	BlockNo uint64 `json:"blockNo,omitempty" dynamodbav:"blockNo,omitempty"`
}

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResponsePraos ResponsePraos

func (c *CompatibleResponsePraos) UnmarshalJSON(data []byte) error {
	var r ResponsePraos
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
		return nil
	} else {
		// All we really care about is the result.
		if r5.Result.IntersectionFound != nil {
			c.Method = FindIntersectionMethod

			var p Point
			p.pointType = r5.Result.IntersectionFound.Point.pointType
			if p.pointType == PointTypeString {
				p.pointString = r5.Result.IntersectionFound.Point.pointString
			} else if p.pointType == PointTypeStruct {
				p.pointStruct.Slot = r5.Result.IntersectionFound.Point.pointStruct.Slot
				p.pointStruct.ID = r5.Result.IntersectionFound.Point.pointStruct.Hash
			} else {
				panic("Invalid point type")
			}

			var t Tip
			t.Slot = r5.Result.IntersectionFound.Tip.Slot
			t.ID = r5.Result.IntersectionFound.Tip.Hash
			t.Height = r5.Result.IntersectionFound.Tip.BlockNo

			var findIntersection CompatibleResultFindIntersection
			findIntersection.Intersection = &p
			findIntersection.Tip = &t
			c.Result = &findIntersection
		} else if r5.Result.IntersectionNotFound != nil {
			c.Method = FindIntersectionMethod

			var t Tip
			t.Slot = r5.Result.IntersectionNotFound.Tip.Slot
			t.ID = r5.Result.IntersectionNotFound.Tip.Hash
			t.Height = r5.Result.IntersectionFound.Tip.BlockNo

			var e ResultError
			e.Data = &t
			e.Code = 1000
			e.Message = "Intersection not found - Conversion from a v5 Ogmigo call"
			c.Error = &e
		} else if r5.Result.RollForward != nil {
			c.Method = NextBlockMethod
			var t Tip
			t.Slot = r5.Result.RollForward.Tip.Slot
			t.ID = r5.Result.RollForward.Tip.Hash
			t.Height = r5.Result.RollForward.Tip.BlockNo

			txArray := make([]Tx, len(r5.Result.RollForward.Block.Body))
			for _, t := range r5.Result.RollForward.Block.Body {
				withdrawals := make(map[string]Lovelace)
				for txid, amt := range t.Body.Withdrawals {
					withdrawals[txid] = Lovelace{Lovelace: num.Int64(amt), Extras: nil}
				}

				tx := Tx{
					ID:                       t.ID,
					Spends:                   t.InputSource,
					Inputs:                   t.Body.Inputs,
					References:               t.Body.References,
					Collaterals:              t.Body.Collaterals,
					TotalCollateral:          t.Body.TotalCollateral,
					CollateralReturn:         (*TxOut)(t.Body.CollateralReturn),
					Outputs:                  t.Body.Outputs,
					Certificates:             t.Body.Certificates,
					Withdrawals:              withdrawals,
					Fee:                      Lovelace{Lovelace: t.Body.Fee, Extras: nil},
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

			fakeNonce := Nonce{Output: "fake", Proof: "fake"}
			protocolVersion := ProtocolVersion{
				Major: uint32(r5.Result.RollForward.Block.Header.ProtocolVersion["Major"]),
				Minor: uint32(r5.Result.RollForward.Block.Header.ProtocolVersion["Minor"]),
				Patch: uint32(r5.Result.RollForward.Block.Header.ProtocolVersion["Patch"]),
			}
			protocol := Protocol{Version: protocolVersion}
			leaderValue := LeaderValue{Proof: "fake", Output: "fake"}
			var opCert OpCert
			if r5.Result.RollForward.Block.Header.OpCert != nil {
				var vk []byte
				if r5.Result.RollForward.Block.Header.OpCert["hotVk"].([]byte) != nil {
					vk, _ = base64.StdEncoding.DecodeString(r5.Result.RollForward.Block.Header.OpCert["hotVk"].(string))
				}
				opCert = OpCert{
					Count: r5.Result.RollForward.Block.Header.OpCert["count"].(uint64),
					Kes:   Kes{Period: r5.Result.RollForward.Block.Header.OpCert["kesPeriod"].(uint64), VerificationKey: string(vk)},
				}
			}
			issuer := BlockIssuer{VerificationKey: r5.Result.RollForward.Block.Header.IssuerVK, VrfVerificationKey: r5.Result.RollForward.Block.Header.IssuerVrf, OperationalCertificate: opCert, LeaderValue: leaderValue}
			b := Block{
				Type:         "praos",
				Era:          "babbage", // TODO - Get from V5 entry - Not trivial as designed
				ID:           r5.Result.RollForward.Block.HeaderHash,
				Ancestor:     r5.Result.RollForward.Block.Header.PrevHash,
				Nonce:        fakeNonce,
				Height:       r5.Result.RollForward.Block.Header.BlockHeight,
				Size:         BlockSize{Bytes: int64(r5.Result.RollForward.Block.Header.BlockSize)},
				Slot:         r5.Result.RollForward.Block.Header.Slot,
				Transactions: txArray,
				Protocol:     protocol,
				Issuer:       issuer,
			}

			var nextBlock CompatibleResultNextBlock
			nextBlock.Direction = "forward"
			nextBlock.Tip = &t
			nextBlock.Block = &b
			c.Result = &nextBlock
		} else if r5.Result.RollBackward != nil {
			c.Method = NextBlockMethod
			var t Tip
			t.Slot = r5.Result.RollBackward.Tip.Slot
			t.ID = r5.Result.RollBackward.Tip.Hash
			t.Height = r5.Result.RollBackward.Tip.BlockNo

			var p Point
			p.pointType = r5.Result.RollBackward.Point.pointType
			if p.pointType == PointTypeString {
				p.pointString = r5.Result.RollBackward.Point.pointString
			} else {
				p.pointStruct.Slot = r5.Result.RollBackward.Point.pointStruct.Slot
				p.pointStruct.ID = r5.Result.RollBackward.Point.pointStruct.Hash
			}

			var nextBlock CompatibleResultNextBlock
			nextBlock.Direction = "backward"
			nextBlock.Tip = &t
			nextBlock.Point = &p
			c.Result = &nextBlock
		}
		c.ID = r5.Reflection
		return nil
	}

	// TODO: Further error handling here.
	return nil
}

type ResponseV5 struct {
	Type        string          `json:"type,omitempty"        dynamodbav:"type,omitempty"`
	Version     string          `json:"version,omitempty"     dynamodbav:"version,omitempty"`
	ServiceName string          `json:"servicename,omitempty" dynamodbav:"servicename,omitempty"`
	MethodName  string          `json:"methodname,omitempty"  dynamodbav:"methodname,omitempty"`
	Result      *ResultV5       `json:"result,omitempty"      dynamodbav:"result,omitempty"`
	Reflection  json.RawMessage `json:"reflection,omitempty"  dynamodbav:"reflection,omitempty"`
}

type ResponsePraos struct {
	JsonRpc string          `json:"jsonrpc,omitempty" dynamodbav:"jsonrpc,omitempty"`
	Method  string          `json:"method,omitempty"  dynamodbav:"method,omitempty"`
	Result  interface{}     `json:"result,omitempty"  dynamodbav:"result,omitempty"`
	Error   *ResultError    `json:"error,omitempty"   dynamodbav:"error,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"      dynamodbav:"id,omitempty"`
}

const FindIntersectionMethod = "findIntersection"
const NextBlockMethod = "nextBlock"
const FindIntersectMethod = "FindIntersect"
const RequestNextMethod = "RequestNext"

func (r *ResponsePraos) UnmarshalJSON(b []byte) error {
	var m struct {
		JsonRpc string          `json:"jsonrpc"`
		Method  string          `json:"method" json:"methodname"`
		ID      json.RawMessage `json:"ID"`
		Result  json.RawMessage `json:"result"`
		Error   json.RawMessage `json:"error"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	r.JsonRpc = m.JsonRpc
	r.ID = m.ID

	if m.Error != nil {
		var resultError ResultError
		if err := json.Unmarshal(m.Error, &resultError); err != nil {
			return err
		}
		r.Error = &resultError
	} else {
		switch m.Method {
		case FindIntersectionMethod, FindIntersectMethod:
			r.Method = FindIntersectionMethod
			var findIntersection CompatibleResultFindIntersection
			if err := json.Unmarshal(m.Result, &findIntersection); err != nil {
				return err
			}
			r.Result = findIntersection

		case NextBlockMethod, RequestNextMethod:
			r.Method = NextBlockMethod
			var nextBlock CompatibleResultNextBlock
			if err := json.Unmarshal(m.Result, &nextBlock); err != nil {
				return err
			}
			r.Result = nextBlock

		default:
			return fmt.Errorf("unknown method: '%v'", r.Method)
		}
	}

	return nil
}

func (r ResponsePraos) MustFindIntersectResult() CompatibleResultFindIntersection {
	if r.Method != FindIntersectionMethod {
		panic(fmt.Errorf("must only use *Must* methods after switching on the findIntersection method; called on %v", r.Method))
	}
	return r.Result.(CompatibleResultFindIntersection)
}

func (r ResponsePraos) MustNextBlockResult() CompatibleResultNextBlock {
	if r.Method != NextBlockMethod {
		panic(fmt.Errorf("must only use *Must* methods after switching on the nextBlock method; called on %v", r.Method))
	}
	fmt.Printf("type of r.Result is %T\n", r.Result)
	return r.Result.(CompatibleResultNextBlock)
}

type Tx struct {
	ID                       string                `json:"id,omitempty"                       dynamodbav:"id,omitempty"`
	Spends                   string                `json:"spends,omitempty"                   dynamodbav:"spends,omitempty"`
	Inputs                   []TxIn                `json:"inputs,omitempty"                   dynamodbav:"inputs,omitempty"`
	References               []TxIn                `json:"references,omitempty"               dynamodbav:"references,omitempty"`
	Collaterals              []TxIn                `json:"collaterals,omitempty"              dynamodbav:"collaterals,omitempty"`
	TotalCollateral          *int64                `json:"totalCollateral,omitempty"          dynamodbav:"totalCollateral,omitempty"`
	CollateralReturn         *TxOut                `json:"collateralReturn,omitempty"         dynamodbav:"collateralReturn,omitempty"`
	Outputs                  TxOuts                `json:"outputs,omitempty"                  dynamodbav:"outputs,omitempty"`
	Certificates             []json.RawMessage     `json:"certificates,omitempty"             dynamodbav:"certificates,omitempty"`
	Withdrawals              map[string]Lovelace   `json:"withdrawals,omitempty"              dynamodbav:"withdrawals,omitempty"`
	Fee                      Lovelace              `json:"fee,omitempty"                      dynamodbav:"fee,omitempty"`
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

type DoubleNestedInteger map[string]map[string]num.Int

type Lovelace struct {
	Lovelace num.Int `json:"lovelace,omitempty"  dynamodbav:"lovelace,omitempty"`
	Extras   []DoubleNestedInteger
}

type TxID string

func NewTxID(txHash string, index int) TxID {
	return TxID(txHash + "#" + strconv.Itoa(index))
}

func (t TxID) String() string {
	return string(t)
}

func (t TxID) Index() int {
	if index := strings.Index(string(t), "#"); index > 0 {
		if v, err := strconv.Atoi(string(t[index+1:])); err == nil {
			return v
		}
	}
	return -1
}

func (t TxID) TxHash() string {
	if index := strings.Index(string(t), "#"); index > 0 {
		return string(t[0:index])
	}
	return ""
}

type TxIn struct {
	TxHash string `json:"txId"  dynamodbav:"txId"`
	Index  int    `json:"index" dynamodbav:"index"`
}

func (t TxIn) String() string {
	return t.TxHash + "#" + strconv.Itoa(t.Index)
}

func (t TxIn) TxID() TxID {
	return NewTxID(t.TxHash, t.Index)
}

type TxOut struct {
	Address   string          `json:"address,omitempty"   dynamodbav:"address,omitempty"`
	Datum     string          `json:"datum,omitempty"     dynamodbav:"datum,omitempty"`
	DatumHash string          `json:"datumHash,omitempty" dynamodbav:"datumHash,omitempty"`
	Value     Value           `json:"value,omitempty"     dynamodbav:"value,omitempty"`
	Script    json.RawMessage `json:"script,omitempty"    dynamodbav:"script,omitempty"`
}

type TxOuts []TxOut

func (tt TxOuts) FindByAssetID(assetID AssetID) (TxOut, bool) {
	for _, t := range tt {
		for gotAssetID := range t.Value.Assets {
			if gotAssetID == assetID {
				return t, true
			}
		}
	}
	return TxOut{}, false
}

type Datums map[string]string

type TxInQuery struct {
	Transaction UtxoTxID `json:"transaction"  dynamodbav:"transaction"`
	Index       uint32   `json:"index" dynamodbav:"index"`
}

type UtxoTxID struct {
	ID string `json:"id"`
}

func (d *Datums) UnmarshalJSON(i []byte) error {
	if i == nil {
		return nil
	}

	var raw map[string]interface{}
	err := json.Unmarshal(i, &raw)
	if err != nil {
		return fmt.Errorf("unable to unmarshal as raw map: %w", err)
	}

	results := make(Datums, len(raw))
	// for backwards compatibility, since ogmios switched Datum values from []byte to hex string
	// this should be safe to remove after we upgrade all ogmios nodes to >= 5.5.0
	for k, v := range raw {
		s, isString := v.(string)
		if !isString {
			return fmt.Errorf("expecting string, got %v", v)
		}
		asHex := s
		// if it's base64 encoded, convert it to a hex string.
		if _, err := hex.DecodeString(s); err != nil {
			rawDatum, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return fmt.Errorf("unable to decode string %v: %w", s, err)
			}
			asHex = hex.EncodeToString(rawDatum)
		}
		results[k] = asHex
	}

	*d = results
	return nil
}

func (d *Datums) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if item == nil {
		return nil
	}

	var raw map[string]interface{}
	if err := dynamodbattribute.UnmarshalMap(item.M, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal map: %w", err)
	}

	results := make(Datums, len(raw))
	// for backwards compatibility, since ogmios switched Datum values from []byte to hex string
	for k, v := range raw {
		if hexString, ok := v.(string); ok {
			results[k] = hexString
		} else {
			results[k] = hex.EncodeToString(v.([]byte))
		}
	}

	*d = results
	return nil
}

type Witness struct {
	Bootstrap  []json.RawMessage `json:"bootstrap,omitempty"  dynamodbav:"bootstrap,omitempty"`
	Datums     Datums            `json:"datums,omitempty"     dynamodbav:"datums,omitempty"`
	Redeemers  json.RawMessage   `json:"redeemers,omitempty"  dynamodbav:"redeemers,omitempty"`
	Scripts    json.RawMessage   `json:"scripts,omitempty"    dynamodbav:"scripts,omitempty"`
	Signatures map[string]string `json:"signatures,omitempty" dynamodbav:"signatures,omitempty"`
}

type ValidityInterval struct {
	InvalidBefore    uint64 `json:"invalidBefore,omitempty"    dynamodbav:"invalidBefore,omitempty"`
	InvalidHereafter uint64 `json:"invalidHereafter,omitempty" dynamodbav:"invalidHereafter,omitempty"`
}

type Value struct {
	Coins  num.Int             `json:"coins,omitempty"  dynamodbav:"coins,omitempty"`
	Assets map[AssetID]num.Int `json:"assets,omitempty" dynamodbav:"assets,omitempty"`
}

func Add(a Value, b Value) Value {
	var result Value
	result.Coins = a.Coins.Add(b.Coins)
	result.Assets = map[AssetID]num.Int{}
	for assetId, amt := range a.Assets {
		result.Assets[assetId] = amt
	}
	for assetId, amt := range b.Assets {
		result.Assets[assetId] = result.Assets[assetId].Add(amt)
	}
	return result
}
func Subtract(a Value, b Value) Value {
	var result Value
	result.Coins = a.Coins.Sub(b.Coins)
	result.Assets = map[AssetID]num.Int{}
	for assetId, amt := range a.Assets {
		result.Assets[assetId] = amt
	}
	for assetId, amt := range b.Assets {
		result.Assets[assetId] = result.Assets[assetId].Sub(amt)
	}
	return result
}
func Enough(have Value, want Value) (bool, error) {
	if have.Coins.Int64() < want.Coins.Int64() {
		return false, fmt.Errorf("not enough ADA to meet demand")
	}
	for asset, amt := range want.Assets {
		if have.Assets[asset].Int64() < amt.Int64() {
			return false, fmt.Errorf("not enough %v to meet demand", asset)
		}
	}
	return true, nil
}
