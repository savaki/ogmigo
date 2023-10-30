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

package compatibility

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
	v5 "github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/v5"
)

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

	var r5 v5.ResultFindIntersectionV5
	err = json.Unmarshal(data, &r5)
	if err != nil {
		return err
	} else if r5.IntersectionFound != nil {
		p := r5.IntersectionFound.Point.ConvertToV6()
		tip := r5.IntersectionFound.Tip.ConvertToV6()
		c.Intersection = &p
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

	var r5 v5.ResultNextBlockV5
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

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResponsePraos chainsync.ResponsePraos

func (c *CompatibleResponsePraos) UnmarshalJSON(data []byte) error {
	var r chainsync.ResponsePraos
	err := json.Unmarshal(data, &r)
	if err == nil && r.Result != nil {
		*c = CompatibleResponsePraos(r)
		return nil
	}

	var r5 v5.ResponseV5
	err = json.Unmarshal(data, &r5)
	c.JsonRpc = "2.0"
	if err != nil {
		// Just skip all the data processing, as it's useless.
		return err
	}

	// var p Point
	// var t Tip
	// var e ResultError

	// All we really care about is the result.
	if r5.Result.IntersectionFound != nil {
		c.Method = chainsync.FindIntersectionMethod

		p := r5.Result.IntersectionFound.Point.ConvertToV6()
		t := r5.Result.IntersectionFound.Tip.ConvertToV6()

		var findIntersection CompatibleResultFindIntersection
		findIntersection.Intersection = &p
		findIntersection.Tip = &t
		c.Result = &findIntersection
	} else if r5.Result.IntersectionNotFound != nil {
		c.Method = chainsync.FindIntersectionMethod
		t := r5.Result.IntersectionNotFound.Tip.ConvertToV6()
		var e chainsync.ResultError
		e.Data = &t
		e.Code = 1000
		e.Message = "Intersection not found - Conversion from a v5 Ogmigo call"
		c.Error = &e
	} else if r5.Result.RollForward != nil {
		c.Method = chainsync.NextBlockMethod

		t := r5.Result.RollForward.Tip.ConvertToV6()
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
				TotalCollateral:          &chainsync.Lovelace{Lovelace: num.Int64(*t.Body.TotalCollateral)},
				CollateralReturn:         (*chainsync.TxOut)(t.Body.CollateralReturn),
				Outputs:                  t.Body.Outputs,
				Certificates:             t.Body.Certificates,
				Withdrawals:              withdrawals,
				Fee:                      chainsync.Lovelace{Lovelace: t.Body.Fee},
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

		p := r5.Result.RollBackward.Point.ConvertToV6()
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
	return CompatibleResultFindIntersection(r.Result.(chainsync.ResultFindIntersectionPraos))
}

func (r CompatibleResponsePraos) MustNextBlockResult() CompatibleResultNextBlock {
	if r.Method != chainsync.NextBlockMethod {
		panic(fmt.Errorf("must only use *Must* methods after switching on the nextBlock method; called on %v", r.Method))
	}
	return CompatibleResultNextBlock(r.Result.(chainsync.ResultNextBlockPraos))
}

func ConvertPointStructV5ToV6(p v5.PointStructV5) chainsync.PointStruct {
	return chainsync.PointStruct{
		BlockNo: p.BlockNo,
		ID:      p.Hash,
		Slot:    p.Slot,
	}
}
