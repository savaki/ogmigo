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

package ogmigo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/thuannguyen2010/ogmigo/ouroboros/chainsync"
)

type EvaluationResponse struct {
	Type        string
	Version     string
	ServiceName string `json:"servicename"`
	MethodName  string `json:"methodname"`
	Reflection  interface{}
	Result      json.RawMessage
}

// EvaluateTx evaluate the execution units of scripts present in a given transaction, without actually submitting the transaction
// https://ogmios.dev/mini-protocols/local-tx-submission/#evaluatetx
func (c *Client) EvaluateTx(ctx context.Context, cborHex string) (redeemer chainsync.Redeemer, err error) {
	var (
		payload = makePayload("EvaluateTx", Map{"evaluate": cborHex})
		raw     json.RawMessage
	)
	if err := c.query(ctx, payload, &raw); err != nil {
		return nil, fmt.Errorf("failed to evaluate tx: %w", err)
	}

	return readEvaluateTx(raw)
}

func readEvaluateTx(data []byte) (chainsync.Redeemer, error) {
	value, dataType, _, err := jsonparser.Get(data, "result", "EvaluationFailure")
	if err != nil {
		if errors.Is(err, jsonparser.KeyPathNotFoundError) {
			redeemerRaw, _, _, err := jsonparser.Get(data, "result", "EvaluationResult")
			if err != nil {
				return nil, err
			}
			var v chainsync.Redeemer
			err = json.Unmarshal(redeemerRaw, &v)
			if err != nil {
				return nil, fmt.Errorf("cannot parse result: %v to redeemer: %w", string(value), err)
			}
			return v, nil
		}
		return nil, fmt.Errorf("failed to parse EvaluateTx response: %w", err)
	}

	switch dataType {
	case jsonparser.Object:
		return nil, EvaluateTxError{message: value}
	default:
		return nil, fmt.Errorf("EvaluateTx failed: %v", string(value))
	}
}

// EvaluateTxV6 evaluate the execution units of scripts present in a given transaction, without actually submitting the transaction
// https://ogmios.dev/mini-protocols/local-tx-submission/#evaluatetx
func (c *Client) EvaluateTxV6(ctx context.Context, cborHex string) (redeemer chainsync.Redeemer, err error) {
	var (
		payload = makePayloadV6("evaluateTransaction", Map{"transaction": Map{"cbor": cborHex}})
		raw     json.RawMessage
	)
	if err := c.query(ctx, payload, &raw); err != nil {
		return nil, fmt.Errorf("failed to evaluate tx: %w", err)
	}

	return readEvaluateTxV6(raw)
}

// EvaluateTxError encapsulates the EvaluateTx errors and allows the results to be parsed
type EvaluateTxError struct {
	message json.RawMessage
}

// Messages returns the raw messages from EvaluateTxError
func (s EvaluateTxError) Messages() json.RawMessage {
	return s.message
}

// Error implements the error interface
func (s EvaluateTxError) Error() string {
	return fmt.Sprintf("EvaluateTx failed: %v", string(s.message))
}

func readEvaluateTxV6(data []byte) (chainsync.Redeemer, error) {
	value, dataType, _, err := jsonparser.Get(data, "error")
	if err != nil {
		if errors.Is(err, jsonparser.KeyPathNotFoundError) {
			redeemerRaw, _, _, err := jsonparser.Get(data, "result")
			if err != nil {
				return nil, err
			}
			var v chainsync.EvaluationResult
			err = json.Unmarshal(redeemerRaw, &v)
			if err != nil {
				return nil, fmt.Errorf("cannot parse result: %v to redeemer: %w", string(value), err)
			}
			result := make(chainsync.Redeemer)
			for _, item := range v {
				var redeemerValue chainsync.RedeemerValue
				redeemerValue.Memory = item.Budget.Memory
				redeemerValue.Steps = item.Budget.Steps
				result[chainsync.RedeemerKey(item.Validator)] = redeemerValue
			}
			return result, nil
		}
		return nil, fmt.Errorf("failed to parse EvaluateTx response: %w", err)
	}

	switch dataType {
	case jsonparser.Object:
		return nil, EvaluateTxError{message: value}
	default:
		return nil, fmt.Errorf("EvaluateTx failed: %v", string(value))
	}
}
