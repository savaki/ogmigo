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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/shared"
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
	withdrawals := map[string]chainsync.Lovelace{}
	for txid, amt := range t.Body.Withdrawals {
		withdrawals[txid] = chainsync.Lovelace{Lovelace: num.Int64(amt)}
	}

	var tc *chainsync.Lovelace
	if t.Body.TotalCollateral != nil {
		tc = &chainsync.Lovelace{Lovelace: num.Int64(*t.Body.TotalCollateral)}
	}
	var cr *chainsync.TxOut
	if t.Body.CollateralReturn != nil {
		*cr = t.Body.CollateralReturn.ConvertToV6()
	}

	cbor, _ := base64.StdEncoding.DecodeString(t.Raw)
	cborHex := hex.EncodeToString(cbor)
	n, _ := json.Marshal(t.Body.Network)
	tx := chainsync.Tx{
		ID:                       t.ID,
		Spends:                   t.InputSource,
		Inputs:                   t.Body.Inputs.ConvertToV6(),
		References:               t.Body.References.ConvertToV6(),
		Collaterals:              t.Body.Collaterals.ConvertToV6(),
		TotalCollateral:          tc,
		CollateralReturn:         cr,
		Outputs:                  t.Body.Outputs.ConvertToV6(),
		Certificates:             t.Body.Certificates,
		Withdrawals:              withdrawals,
		Fee:                      chainsync.Lovelace{Lovelace: t.Body.Fee},
		ValidityInterval:         t.Body.ValidityInterval.ConvertToV6(),
		Mint:                     t.Body.Mint.ConvertToV6(),
		Network:                  string(n),
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
		CBOR:                     cborHex,
	}

	return tx
}

type TxBodyV5 struct {
	Certificates            []json.RawMessage  `json:"certificates,omitempty"            dynamodbav:"certificates,omitempty"`
	Collaterals             TxInsV5            `json:"collaterals,omitempty"             dynamodbav:"collaterals,omitempty"`
	Fee                     num.Int            `json:"fee,omitempty"                     dynamodbav:"fee,omitempty"`
	Inputs                  TxInsV5            `json:"inputs,omitempty"                  dynamodbav:"inputs,omitempty"`
	Mint                    *ValueV5           `json:"mint,omitempty"                    dynamodbav:"mint,omitempty"`
	Network                 json.RawMessage    `json:"network,omitempty"                 dynamodbav:"network,omitempty"`
	Outputs                 TxOutsV5           `json:"outputs,omitempty"                 dynamodbav:"outputs,omitempty"`
	RequiredExtraSignatures []string           `json:"requiredExtraSignatures,omitempty" dynamodbav:"requiredExtraSignatures,omitempty"`
	ScriptIntegrityHash     string             `json:"scriptIntegrityHash,omitempty"     dynamodbav:"scriptIntegrityHash,omitempty"`
	TimeToLive              int64              `json:"timeToLive,omitempty"              dynamodbav:"timeToLive,omitempty"`
	Update                  json.RawMessage    `json:"update,omitempty"                  dynamodbav:"update,omitempty"`
	ValidityInterval        ValidityIntervalV5 `json:"validityInterval"                  dynamodbav:"validityInterval,omitempty"`
	Withdrawals             map[string]int64   `json:"withdrawals,omitempty"             dynamodbav:"withdrawals,omitempty"`
	CollateralReturn        *TxOutV5           `json:"collateralReturn,omitempty"        dynamodbav:"collateralReturn,omitempty"`
	TotalCollateral         *int64             `json:"totalCollateral,omitempty"         dynamodbav:"totalCollateral,omitempty"`
	References              TxInsV5            `json:"references,omitempty"              dynamodbav:"references,omitempty"`
}

type TxInsV5 []TxInV5

func (t TxInsV5) ConvertToV6() chainsync.TxIns {
	var txIns chainsync.TxIns
	for _, txIn := range t {
		txIns = append(txIns, txIn.ConvertToV6())
	}
	return txIns
}

type TxInV5 struct {
	TxHash string `json:"txId"  dynamodbav:"txId"`
	Index  int    `json:"index" dynamodbav:"index"`
}

func (t TxInV5) String() string {
	return t.TxHash + "#" + strconv.Itoa(t.Index)
}

func (t TxInV5) ConvertToV6() chainsync.TxIn {
	id := chainsync.TxInID{ID: t.TxHash}
	return chainsync.TxIn{Transaction: id, Index: t.Index}
}

type TxOutV5 struct {
	Address   string          `json:"address,omitempty"   dynamodbav:"address,omitempty"`
	Datum     string          `json:"datum,omitempty"     dynamodbav:"datum,omitempty"`
	DatumHash string          `json:"datumHash,omitempty" dynamodbav:"datumHash,omitempty"`
	Value     ValueV5         `json:"value,omitempty"     dynamodbav:"value,omitempty"`
	Script    json.RawMessage `json:"script,omitempty"    dynamodbav:"script,omitempty"`
}

func (t TxOutV5) ConvertToV6() chainsync.TxOut {
	return chainsync.TxOut{
		Address:   t.Address,
		Datum:     t.Datum,
		DatumHash: t.DatumHash,
		Value:     t.Value.ConvertToV6(),
		Script:    t.Script,
	}
}

type TxOutsV5 []TxOutV5

func (t TxOutsV5) ConvertToV6() chainsync.TxOuts {
	txOuts := make(chainsync.TxOuts, len(t))
	for _, txOut := range t {
		txOuts = append(txOuts, txOut.ConvertToV6())
	}
	return txOuts
}

func (tt TxOutsV5) FindByAssetID(assetID shared.AssetID) (TxOutV5, bool) {
	for _, t := range tt {
		for gotAssetID := range t.Value.Assets {
			if gotAssetID == assetID {
				return t, true
			}
		}
	}
	return TxOutV5{}, false
}

type ValidityIntervalV5 struct {
	InvalidBefore    uint64 `json:"invalidBefore,omitempty"     dynamodbav:"invalidBefore,omitempty"`
	InvalidHereafter uint64 `json:"invalidHereafter,omitempty"  dynamodbav:"invalidHereafter,omitempty"`
}

func (v ValidityIntervalV5) ConvertToV6() chainsync.ValidityInterval {
	return chainsync.ValidityInterval{
		InvalidBefore: v.InvalidBefore,
		InvalidAfter:  v.InvalidHereafter,
	}
}

type ValueV5 struct {
	Coins  num.Int                    `json:"coins,omitempty"  dynamodbav:"coins,omitempty"`
	Assets map[shared.AssetID]num.Int `json:"assets,omitempty" dynamodbav:"assets,omitempty"`
}

func (v ValueV5) ConvertToV6() shared.Value {
	assets := shared.Value{}
	if v.Coins.Uint64() != 0 {
		assets[shared.AdaPolicy] = map[string]num.Int{
			shared.AdaAsset: v.Coins,
		}
	}
	for assetId, assetNum := range v.Assets {
		policySplit := strings.Split(string(assetId), ".")
		var (
			policyId  string
			assetName string
		)
		if len(policySplit) == 2 {
			policyId = policySplit[0]
			assetName = policySplit[1]
		} else {
			policyId = policySplit[0]
		}
		if assets[policyId] == nil {
			assets[policyId] = map[string]num.Int{}
		}
		assets[policyId][assetName] = assetNum
	}

	return assets
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
func (b BlockV5) PointStruct() PointStructV5 {
	return PointStructV5{
		BlockNo: b.Header.BlockHeight,
		Hash:    b.HeaderHash,
		Slot:    b.Header.Slot,
	}
}

func (b BlockV5) ConvertToV6() chainsync.Block {
	txArray := make([]chainsync.Tx, len(b.Body))
	for _, t := range b.Body {
		txArray = append(txArray, t.ConvertToV6())
	}

	// The v5 spec indicates that both nonce entries are optional. We'll create a v6
	// entry (which also indicates both are optional) if either is present.
	nonceOutput := b.Header.Nonce["output"]
	nonceProof := b.Header.Nonce["output"]
	var nonce *chainsync.Nonce
	if nonceOutput != "" || nonceProof != "" {
		*nonce = chainsync.Nonce{Output: nonceOutput, Proof: nonceProof}
	}
	majorVer := uint32(b.Header.ProtocolVersion["major"])
	protocolVersion := chainsync.ProtocolVersion{
		Major: majorVer,
		Minor: uint32(b.Header.ProtocolVersion["minor"]),
		Patch: uint32(b.Header.ProtocolVersion["patch"]),
	}
	protocol := chainsync.Protocol{Version: protocolVersion}

	var opCert chainsync.OpCert
	if b.Header.OpCert != nil {
		var vk []byte
		if b.Header.OpCert["hotVk"] != nil {
			vk, _ = base64.StdEncoding.DecodeString(b.Header.OpCert["hotVk"].(string))
		}
		count := b.Header.OpCert["count"]
		kesPd := b.Header.OpCert["kesPeriod"]

		// Yes, the uint64 casts are ugly. JSON covers floats but not ints. Unmarshalling
		// into interface{} creates float64. If we treat interface{} as uint64, the code
		// compiles but crashes at runtime. So, we cast float64 to uint64.
		opCert = chainsync.OpCert{
			Count: uint64(count.(float64)),
			Kes:   chainsync.Kes{Period: uint64(kesPd.(float64)), VerificationKey: string(vk)},
		}
	}

	var leaderValue *chainsync.LeaderValue
	if b.Header.LeaderValue["output"] != nil && b.Header.LeaderValue["proof"] != nil {
		decodedOutput, _ := base64.StdEncoding.DecodeString(string(b.Header.LeaderValue["output"]))
		decodedProof, _ := base64.StdEncoding.DecodeString(string(b.Header.LeaderValue["proof"]))
		leaderValue = &chainsync.LeaderValue{
			Output: string(decodedOutput),
			Proof:  string(decodedProof),
		}
	}

	issuerVrf, _ := base64.StdEncoding.DecodeString(b.Header.IssuerVrf)
	issuer := chainsync.BlockIssuer{
		VerificationKey:        b.Header.IssuerVK,
		VrfVerificationKey:     string(issuerVrf),
		OperationalCertificate: opCert,
		LeaderValue:            leaderValue,
	}

	// Unfortunately, due to how v5 data is formatted on the wire and stored by Ogmigo,
	// we have to use other data to determine the era. Because the protocol version can
	// be determined before a fork, unlike the slot number, we use the protocol version.
	// This will make it easier to add support for future forks before the fork occurs,
	// while also supporting testnet and other networks that use the same versioning.
	b6 := chainsync.Block{
		Type:         "praos",
		Era:          shared.MajorEras[majorVer],
		ID:           b.HeaderHash,
		Ancestor:     b.Header.PrevHash,
		Nonce:        nonce,
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

func (p PointStructV5) ConvertToV6() chainsync.PointStruct {
	return chainsync.PointStruct{
		BlockNo: p.BlockNo,
		ID:      p.Hash,
		Slot:    p.Slot,
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
	IntersectionFound    *IntersectionFoundV5    `json:",omitempty" dynamodbav:",omitempty"`
	IntersectionNotFound *IntersectionNotFoundV5 `json:",omitempty" dynamodbav:",omitempty"`
	RollForward          *RollForwardV5          `json:",omitempty" dynamodbav:",omitempty"`
	RollBackward         *RollBackwardV5         `json:",omitempty" dynamodbav:",omitempty"`
}

type ResultFindIntersectionV5 struct {
	IntersectionFound    *IntersectionFoundV5    `json:",omitempty" dynamodbav:",omitempty"`
	IntersectionNotFound *IntersectionNotFoundV5 `json:",omitempty" dynamodbav:",omitempty"`
}

type RollBackwardV5 struct {
	Point PointV5 `json:"point,omitempty" dynamodbav:"point,omitempty"`
	Tip   TipV5   `json:"tip,omitempty"   dynamodbav:"tip,omitempty"`
}

type RollForwardBlockV5 struct {
	Allegra *BlockV5    `json:"allegra,omitempty" dynamodbav:"allegra,omitempty"`
	Alonzo  *BlockV5    `json:"alonzo,omitempty"  dynamodbav:"alonzo,omitempty"`
	Babbage *BlockV5    `json:"babbage,omitempty" dynamodbav:"babbage,omitempty"`
	Byron   *ByronBlock `json:"byron,omitempty"   dynamodbav:"byron,omitempty"`
	Mary    *BlockV5    `json:"mary,omitempty"    dynamodbav:"mary,omitempty"`
	Shelley *BlockV5    `json:"shelley,omitempty" dynamodbav:"shelley,omitempty"`
}

func (b RollForwardBlockV5) GetNonByronBlock() *BlockV5 {
	if b.Shelley != nil {
		return b.Shelley
	} else if b.Allegra != nil {
		return b.Allegra
	} else if b.Mary != nil {
		return b.Mary
	} else if b.Alonzo != nil {
		return b.Alonzo
	} else if b.Babbage != nil {
		return b.Babbage
	} else {
		return nil
	}
}

type RollForwardV5 struct {
	Block RollForwardBlockV5 `json:"block,omitempty" dynamodbav:"block,omitempty"`
	Tip   TipV5              `json:"tip,omitempty"   dynamodbav:"tip,omitempty"`
}

type ResultNextBlockV5 struct {
	RollForward  *RollForwardV5  `json:",omitempty" dynamodbav:",omitempty"`
	RollBackward *RollBackwardV5 `json:",omitempty" dynamodbav:",omitempty"`
}

type IntersectionFoundV5 struct {
	Point *PointV5
	Tip   *TipV5
}

type IntersectionNotFoundV5 struct {
	Tip *TipV5
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

type ResponseV5 struct {
	Type        string          `json:"type,omitempty"        dynamodbav:"type,omitempty"`
	Version     string          `json:"version,omitempty"     dynamodbav:"version,omitempty"`
	ServiceName string          `json:"servicename,omitempty" dynamodbav:"servicename,omitempty"`
	MethodName  string          `json:"methodname,omitempty"  dynamodbav:"methodname,omitempty"`
	Result      *ResultV5       `json:"result,omitempty"      dynamodbav:"result,omitempty"`
	Reflection  json.RawMessage `json:"reflection,omitempty"  dynamodbav:"reflection,omitempty"`
}