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
	"fmt"
)

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
		payload = makePayload("SubmitTx", Map{"bytes": signedTx})
		got     json.RawMessage
	)
	if err := c.query(ctx, payload, &got); err != nil {
		return fmt.Errorf("failed to submit tx: %w", err)
	}

	fmt.Println(string(got))

	return nil
}
