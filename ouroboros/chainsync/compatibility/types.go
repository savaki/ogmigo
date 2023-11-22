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

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Support findIntersect (v6) / FindIntersection (v5) universally.
type CompatibleResultFindIntersection chainsync.ResultFindIntersectionPraos

func (c *CompatibleResultFindIntersection) UnmarshalJSON(data []byte) error {
	// Assume v6 responses first, then fall back to manual v5 processing.
	var r chainsync.ResultFindIntersectionPraos
	err := json.Unmarshal(data, &r)
	// We check intersection here, as that key is distinct from the other result types
	if err == nil && r.Intersection != nil {
		*c = CompatibleResultFindIntersection(r)
		return nil
	}

	var r5 v5.ResultFindIntersectionV5
	err = json.Unmarshal(data, &r5)
	if err == nil && (r5.IntersectionFound != nil || r5.IntersectionNotFound != nil) {
		*c = CompatibleResultFindIntersection(r5.ConvertToV6())
		return nil
	} else {
		return fmt.Errorf("unable to parse as either v5 or v6 FindIntersection: %w", err)
	}
}

func (c *CompatibleResultFindIntersection) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	var s chainsync.ResultFindIntersectionPraos
	err := dynamodbattribute.Unmarshal(item, &s)
	if err == nil && s.Intersection != nil {
		*c = CompatibleResultFindIntersection(s)
		return nil
	}

	var v v5.ResultFindIntersectionV5
	err = dynamodbattribute.Unmarshal(item, &v)
	if err == nil && s.Intersection != nil {
		*c = CompatibleResultFindIntersection(v.ConvertToV6())
		return nil
	} else {
		return fmt.Errorf("unable to parse as either v5 or v6 FindIntersection: %w", err)
	}
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
	if err == nil && r.Direction != "" {
		*c = CompatibleResultNextBlock(r)
		return nil
	}

	var v v5.ResultNextBlockV5
	err = json.Unmarshal(data, &v)
	if err == nil && (v.RollBackward != nil || v.RollForward != nil) {
		*c = CompatibleResultNextBlock(v.ConvertToV6())
		return nil
	} else {
		return fmt.Errorf("unable to parse as either v5 of v6 NextBlock: %w", err)
	}
}

func (c *CompatibleResultNextBlock) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	var s chainsync.ResultNextBlockPraos
	err := dynamodbattribute.Unmarshal(item, &s)
	if err == nil && s.Direction != "" {
		*c = CompatibleResultNextBlock(s)
		return nil
	}

	var v v5.ResultNextBlockV5
	err = dynamodbattribute.Unmarshal(item, &v)
	if err == nil && (v.RollBackward != nil || v.RollForward != nil) {
		*c = CompatibleResultNextBlock(v.ConvertToV6())
		return nil
	} else {
		return fmt.Errorf("unable to parse as either v5 or v6 NextBlock: %w", err)
	}
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
	if err != nil {
		// Just skip all the data processing, as it's useless.
		return err
	}

	*c = CompatibleResponsePraos(r5.ConvertToV6())
	return nil
}

func (c *CompatibleResponsePraos) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	var s chainsync.ResponsePraos
	if err := dynamodbattribute.Unmarshal(item, &s); err != nil {
		var v v5.ResponseV5
		if err := dynamodbattribute.Unmarshal(item, &v); err != nil {
			return err
		}
		*c = CompatibleResponsePraos(v.ConvertToV6())
		return nil
	}
	*c = CompatibleResponsePraos(s)
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

func (c *CompatibleValue) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	var s shared.Value
	if err := dynamodbattribute.Unmarshal(item, &s); err != nil {
		var v v5.ValueV5
		if err := dynamodbattribute.Unmarshal(item, &v); err != nil {
			return err
		}
		*c = CompatibleValue(v.ConvertToV6())
		return nil
	}
	*c = CompatibleValue(s)
	return nil
}

type CompatibleResult struct {
	NextBlock        *CompatibleResultNextBlock
	FindIntersection *CompatibleResultFindIntersection
}

func (c *CompatibleResult) UnmarshalJSON(data []byte) error {
	var rfi CompatibleResultFindIntersection
	err := json.Unmarshal(data, &rfi)
	if err == nil {
		*c.FindIntersection = rfi
		return nil
	}

	var rnb CompatibleResultNextBlock
	err = json.Unmarshal(data, &rnb)
	if err == nil {
		*c.NextBlock = rnb
		return nil
	}

	return fmt.Errorf("unable to find an appropriate result")
}

func (c *CompatibleResult) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	var rfi CompatibleResultFindIntersection
	if err := dynamodbattribute.Unmarshal(item, &rfi); err != nil {
		var rnb CompatibleResultNextBlock
		if err := dynamodbattribute.Unmarshal(item, &rnb); err != nil {
			return err
		}
		*c.NextBlock = rnb
		return nil
	}
	*c.FindIntersection = rfi
	return nil
}
