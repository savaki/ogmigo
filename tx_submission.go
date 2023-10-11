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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/buger/jsonparser"
)

type Response struct {
	Type        string
	Version     string
	ServiceName string `json:"servicename"`
	MethodName  string `json:"methodname"`
	Reflection  interface{}
	Result      json.RawMessage
}

// SubmitTx submits the transaction via ogmios
// https://ogmios.dev/mini-protocols/local-tx-submission/
func (c *Client) SubmitTx(ctx context.Context, data []byte) (err error) {
	var content struct{ CborHex string }
	if err := json.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("failed to decode signed tx: %w", err)
	}

	signedTx := content.CborHex
	if signedTx == "" {
		signedTx = string(data)
	}

	var (
		payload = makePayload("submitTransaction", Map{"transaction": signedTx})
		raw     json.RawMessage
	)
	if err := c.query(ctx, payload, &raw); err != nil {
		return fmt.Errorf("failed to submit tx: %w", err)
	}

	return readSubmitTx(raw)
}

// SubmitTxError encapsulates the SubmitTx errors and allows the results to be parsed
type SubmitTxError struct {
	messages []json.RawMessage
}

// HasErrorCode returns true if the error contains the provided code
func (s SubmitTxError) HasErrorCode(errorCode string) bool {
	errorCodes, _ := s.ErrorCodes()
	for _, ec := range errorCodes {
		if ec == errorCode {
			return true
		}
	}
	return false
}

// ErrorCodes the list of errors codes
func (s SubmitTxError) ErrorCodes() (keys []string, err error) {
	for _, data := range s.messages {
		if bytes.HasPrefix(data, []byte(`"`)) {
			var key string
			if err := json.Unmarshal(data, &key); err != nil {
				return nil, fmt.Errorf("failed to decode string, %v", string(data))
			}
			keys = append(keys, key)
			continue
		}

		var messages map[string]json.RawMessage
		if err := json.Unmarshal(data, &messages); err != nil {
			return nil, fmt.Errorf("failed to decode object, %v", string(data))
		}

		for key := range messages {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys, nil
}

// Messages returns the raw messages from SubmitTxError
func (s SubmitTxError) Messages() []json.RawMessage {
	return s.messages
}

// Error implements the error interface
func (s SubmitTxError) Error() string {
	keys, _ := s.ErrorCodes()
	return fmt.Sprintf("SubmitTx failed: %v", strings.Join(keys, ", "))
}

func readSubmitTx(data []byte) error {
	value, dataType, _, err := jsonparser.Get(data, "error")
	if err != nil {
		if errors.Is(err, jsonparser.KeyPathNotFoundError) {
			return nil
		}
		return fmt.Errorf("failed to parse SubmitTx response: %w", err)
	}

	switch dataType {
	case jsonparser.Array:
		var messages []json.RawMessage
		if err := json.Unmarshal(value, &messages); err != nil {
			return fmt.Errorf("failed to parse SubmitTx response: array: %w", err)
		}
		if len(messages) == 0 {
			return nil
		}
		return SubmitTxError{messages: messages}

	case jsonparser.Object:
		return SubmitTxError{messages: []json.RawMessage{value}}

	default:
		return fmt.Errorf("SubmitTx failed: %v", string(value))
	}
}
