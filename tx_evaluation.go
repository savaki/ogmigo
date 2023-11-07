// Copyright 2023 Sundae Labs
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

package ogmigo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
)

type EvaluateResponse struct {
	Type        string
	Version     string
	ServiceName string `json:"servicename"`
	MethodName  string `json:"methodname"`
	Reflection  interface{}
	Result      json.RawMessage
}

type EvaluateTx struct {
	Cbor string `json:"cbor"`
}

// EvaluateTx measures the script execution costs of a transaction.
// TODO: Support additionalUtxoSet
// https://ogmios.dev/mini-protocols/local-tx-submission/
// https://github.com/CardanoSolutions/ogmios/blob/v6.0/docs/content/mini-protocols/local-tx-submission.md
func (c *Client) EvaluateTx(ctx context.Context, data string) (response *EvaluateTxResponse, err error) {
	tx := EvaluateTx{
		Cbor: data,
	}

	var (
		payload = makePayload("evaluateTransaction", Map{"transaction": tx})
		raw     json.RawMessage
	)
	if err := c.query(ctx, payload, &raw); err != nil {
		return nil, fmt.Errorf("failed to evaluate tx: %w", err)
	}

	return readEvaluateTx(raw)
}

type ExUnits struct {
	Validator string        `json:"validator"`
	Budget    ExUnitsBudget `json:"budget"`
}

type ExUnitsBudget struct {
	Memory uint64 `json:"memory"`
	Cpu    uint64 `json:"cpu"`
}

type EvaluateTxError struct {
	Code    int
	Message string
	Data    json.RawMessage
}

type EvaluateTxResponse struct {
	ExUnits []ExUnits
	Error   *EvaluateTxError
}

func readEvaluateTx(data []byte) (r *EvaluateTxResponse, err error) {
	e, err1 := readEvaluateTxError(data)
	u, err2 := readEvaluateTxResult(data)
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("could not parse evaluate tx response; neither error (%w) nor result (%w)", err1, err2)
	}
	if err1 == nil {
		return &EvaluateTxResponse{Error: e}, nil
	}
	if err2 == nil {
		return &EvaluateTxResponse{ExUnits: u}, nil
	}
	return nil, fmt.Errorf("could not parse evaluate tx response: %s", string(data))
}

func readEvaluateTxError(data []byte) (*EvaluateTxError, error) {
	value, _, _, err := jsonparser.Get(data, "error")
	if err != nil {
		return nil, fmt.Errorf("failed to parse EvaluateTx error: %w %s", err, data)
	}
	var e EvaluateTxError
	if err := json.Unmarshal(value, &e); err != nil {
		return nil, fmt.Errorf("failed to parse EvaluateTx error: %w %s", err, data)
	}
	return &e, nil
}

func readEvaluateTxResult(data []byte) ([]ExUnits, error) {
	value, dataType, _, err := jsonparser.Get(data, "result")
	if err != nil {
		return nil, fmt.Errorf("failed to parse EvaluateTx response: %w %s", err, string(data))
	}

	switch dataType {
	case jsonparser.Array:
		var results []ExUnits
		if err := json.Unmarshal(value, &results); err != nil {
			return nil, fmt.Errorf("failed to parse EvaluateTx response: %w", err)
		}
		return results, nil
	default:
		return nil, fmt.Errorf("failed to parser EvaluateTx response: %w", err)
	}
}
