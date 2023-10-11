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
	Direction string `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Tip       Tip    `json:"tip,omitempty"   dynamodbav:"tip,omitempty"`
	Point     Point  `json:"point,omitempty" dynamodbav:"point,omitempty"`
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

// Covers all blocks except Byron-era blocks.
type ResultPraos struct {
	Direction string `json:"direction,omitempty" dynamodbav:"direction,omitempty"`
	Tip       *Tip   `json:"tip,omitempty"       dynamodbav:"tip,omitempty"`
	Block     *Block `json:"block,omitempty"     dynamodbav:"block,omitempty"` // Forward
	Point     Point  `json:"point,omitempty"     dynamodbav:"point,omitempty"` // Backward
}

type Tip struct {
	Slot   uint64 `json:"slot,omitempty"   dynamodbav:"slot,omitempty"`
	ID     string `json:"id,omitempty"     dynamodbav:"id,omitempty"`
	Height uint64 `json:"height,omitempty" dynamodbav:"height,omitempty"`
}

type ResponsePraos struct {
	JsonRpc string          `json:"jsonrpc,omitempty" dynamodbav:"jsonrpc,omitempty"`
	Method  string          `json:"method,omitempty"  dynamodbav:"method,omitempty"`
	Result  *ResultPraos    `json:"result,omitempty"  dynamodbav:"result,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"      dynamodbav:"id,omitempty"`
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
