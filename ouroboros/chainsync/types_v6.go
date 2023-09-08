package chainsync

import (
	"encoding/json"
	"fmt"
	"github.com/thuannguyen2010/ogmigo/ouroboros/chainsync/num"
)

type ResponseV6 struct {
	Jsonrpc string    `json:"jsonrpc,omitempty"`
	Method  string    `json:"method,omitempty"`
	Result  *ResultV6 `json:"result,omitempty"`
}

type ResultV6 struct {
	Direction    string   `json:"direction,omitempty"`
	BlockV6      *BlockV6 `json:"block,omitempty"`
	Point        *PointV6 `json:"point,omitempty"`        // roll backward (nextBlock)
	Intersection *PointV6 `json:"intersection,omitempty"` // intersection found (findIntersection)
	Error        *Error   `json:"error,omitempty"`        // intersection not found (findIntersection)
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type BlockV6 struct {
	Type         string        `json:"type,omitempty"`
	Era          string        `json:"era,omitempty"`
	ID           string        `json:"id,omitempty"`
	Height       uint64        `json:"height,omitempty"`
	Slot         uint64        `json:"slot,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

type Transaction struct {
	ID        string             `json:"id,omitempty"`
	Spends    string             `json:"spends,omitempty"`
	TxInputs  []TxInV6           `json:"inputs,omitempty"`
	TxOutputs []TxOutV6          `json:"outputs,omitempty"`
	Datums    map[string]HexData `json:"datums,omitempty"`
	Redeemers json.RawMessage    `json:"redeemers,omitempty"`
	Fee       Fee                `json:"fee"`
}

type Fee struct {
	Lovelace num.Int `json:"lovelace,omitempty"`
}

type TxInV6 struct {
	Transaction struct {
		ID string `json:"id,omitempty"`
	} `json:"transaction,omitempty"`
	Index int `json:"index,omitempty"`
}

type TxOutV6 struct {
	Address   string `json:"address,omitempty"`
	Value     map[string]map[string]num.Int
	Datum     string `json:"datum"`
	DatumHash string `json:"datumHash"`
}

type PointV6 struct {
	pointType   PointType
	pointString PointString
	pointStruct *PointStructV6
}

type PointsV6 []PointV6

type PointStructV6 struct {
	ID   string `json:"id,omitempty"`
	Slot uint64 `json:"slot,omitempty"`
}

func (p *PointV6) UnmarshalJSON(data []byte) error {
	switch {
	case data[0] == '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("failed to unmarshal Point, %v: %w", string(data), err)
		}

		*p = PointV6{
			pointType:   PointTypeString,
			pointString: PointString(s),
		}

	default:
		var ps PointStructV6
		if err := json.Unmarshal(data, &ps); err != nil {
			return fmt.Errorf("failed to unmarshal Point, %v: %w", string(data), err)
		}

		*p = PointV6{
			pointType:   PointTypeStruct,
			pointStruct: &ps,
		}
	}

	return nil
}

func (p *PointV6) convertToV5() Point {
	var pointStruct *PointStruct
	if p.pointStruct != nil {
		pointStruct = &PointStruct{
			Hash: p.pointStruct.ID,
			Slot: p.pointStruct.Slot,
		}
	}
	return Point{
		pointType:   p.pointType,
		pointString: p.pointString,
		pointStruct: pointStruct,
	}
}

func (p PointV6) MarshalJSON() ([]byte, error) {
	switch p.pointType {
	case PointTypeString:
		return json.Marshal(p.pointString)
	case PointTypeStruct:
		return json.Marshal(p.pointStruct)
	default:
		return nil, fmt.Errorf("unable to unmarshal Point: unknown type")
	}
}

func (p *Point) ConvertToV6() PointV6 {
	var pointStruct *PointStructV6
	if p.pointStruct != nil {
		pointStruct = &PointStructV6{
			ID:   p.pointStruct.Hash,
			Slot: p.pointStruct.Slot,
		}
	}
	return PointV6{
		pointType:   p.pointType,
		pointString: p.pointString,
		pointStruct: pointStruct,
	}
}

func (pp Points) ConvertToV6() PointsV6 {
	var result PointsV6
	for _, p := range pp {
		result = append(result, p.ConvertToV6())
	}
	return result
}

func (responseV6 ResponseV6) ConvertToV5() (response Response) {
	response = Response{
		Type:        responseV6.Jsonrpc,
		Version:     responseV6.Jsonrpc,
		ServiceName: "ogmios",
		MethodName:  responseV6.Method,
		Result:      nil,
	}
	// intersection not found
	if responseV6.Result.Error != nil {
		response.Result = &Result{
			IntersectionNotFound: &IntersectionNotFound{},
		}
		return
	}
	// defend result is nil
	if responseV6.Result == nil {
		return
	}
	// intersection found
	if responseV6.Result.Intersection != nil {
		response.Result = &Result{
			IntersectionFound: &IntersectionFound{
				Point: responseV6.Result.Intersection.convertToV5(),
			},
		}
		return
	}
	// roll backward
	if responseV6.Result.Direction == "backward" {
		response.Result = &Result{
			RollBackward: &RollBackward{},
		}
	}
	// roll forward block
	if responseV6.Result.BlockV6 == nil {
		return
	}
	var txs []Tx
	for _, txV6 := range responseV6.Result.BlockV6.Transactions {
		var txIns []TxIn
		for _, txInV6 := range txV6.TxInputs {
			txIns = append(txIns, TxIn{
				TxHash: txInV6.Transaction.ID,
				Index:  txInV6.Index,
			})
		}

		var txOuts TxOuts
		for _, txOutV6 := range txV6.TxOutputs {
			var coins num.Int
			assets := make(map[AssetID]num.Int)
			for policyID, assetNameMap := range txOutV6.Value {
				if policyID == "ada" {
					for assetName, val := range assetNameMap {
						if assetName == "lovelace" {
							coins = val
						}
					}
				} else {
					for assetName, val := range assetNameMap {
						assetID := AssetID(fmt.Sprintf("%s.%s", policyID, assetName))
						assets[assetID] = val
					}
				}
			}
			txOuts = append(txOuts, TxOut{
				Address:   txOutV6.Address,
				Datum:     txOutV6.Datum,
				DatumHash: txOutV6.DatumHash,
				Value: Value{
					Coins:  coins,
					Assets: assets,
				},
			})
		}
		txs = append(txs, Tx{
			ID: txV6.ID,
			Body: TxBody{
				Fee:     txV6.Fee.Lovelace,
				Inputs:  txIns,
				Outputs: txOuts,
			},
			Witness: Witness{
				Datums:    txV6.Datums,
				Redeemers: txV6.Redeemers,
			},
		})
	}
	block := Block{
		Body: txs,
		Header: BlockHeader{
			BlockHeight: responseV6.Result.BlockV6.Height,
			Slot:        responseV6.Result.BlockV6.Slot,
		},
		HeaderHash: responseV6.Result.BlockV6.ID,
	}
	var rollForwardBlock RollForwardBlock
	switch responseV6.Result.BlockV6.Era {
	case "shelley", "allegra", "mary", "alonzo", "babbage", "conway":
		rollForwardBlock = RollForwardBlock{
			Babbage: &block,
		}
	default:
		rollForwardBlock = RollForwardBlock{}
	}

	result := Result{
		IntersectionFound:    nil,
		IntersectionNotFound: nil,
		RollForward: &RollForward{
			Block: rollForwardBlock,
		},
		RollBackward: nil,
	}

	response.Result = &result
	return response

}
