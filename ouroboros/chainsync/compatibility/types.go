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
	"encoding/json"
	"fmt"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	v5 "github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/v5"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/shared"
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

// Support nextBlock (v6) / RequestNext (v5) universally.
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
		rollForwardBlock := r5.RollForward.Block.GetNonByronBlock()
		if rollForwardBlock == nil {
			return error(fmt.Errorf("rollForwardBlock is nil"))
		}
		block := rollForwardBlock.ConvertToV6()
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

// Frontend for converting v5 JSON responses to v6.
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

	// All we really care about is the result, not the metadata.
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

		rollForwardBlock := r5.Result.RollForward.Block.GetNonByronBlock()
		if rollForwardBlock == nil {
			return error(fmt.Errorf("rollForwardBlock is nil"))
		}
		block := rollForwardBlock.ConvertToV6()

		t := r5.Result.RollForward.Tip.ConvertToV6()

		var nextBlock CompatibleResultNextBlock
		nextBlock.Direction = chainsync.RollForwardString
		nextBlock.Tip = &t
		nextBlock.Block = &block
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

type CompatibleValue shared.Value

func (c *CompatibleValue) UnmarshalJSON(data []byte) error {
	var v shared.Value
	err := json.Unmarshal(data, &v)
	if err == nil {
		*c = CompatibleValue(v)
		return nil
	}

	var r5 v5.ValueV5
	err = json.Unmarshal(data, &r5)
	if err != nil {
		return err
	}

	if r5.Coins.BigInt().BitLen() != 0 {
		shared.Value(*c).AddAsset(shared.Coin{AssetId: shared.AdaAssetID, Amount: r5.Coins})
	}
	for asset, coins := range r5.Assets {
		shared.Value(*c).AddAsset(shared.Coin{AssetId: asset, Amount: coins})
	}

	return nil
}
