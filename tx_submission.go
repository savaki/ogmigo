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
	"fmt"
)

type Response struct {
	Transaction ResponseTx `json:"transaction,omitempty"  dynamodbav:"transaction,omitempty"`
}

type ResponseTx struct {
	ID string `json:"id,omitempty" dynamodbav:"id,omitempty"`
}

// type Response struct {
// 	Type        string
// 	Version     string
// 	ServiceName string `json:"servicename"`
// 	MethodName  string `json:"methodname"`
// 	Reflection  interface{}
// 	Result      json.RawMessage
// }

// SubmitTx submits the transaction via ogmios
// https://ogmios.dev/mini-protocols/local-tx-submission/
func (c *Client) SubmitTx(ctx context.Context, data string) (err error) {
	var (
		payload = makeSubmitTxPayload(data, Map{})
		content struct{ Result Response }
	)
	if err := c.query(ctx, payload, &content); err != nil {
		return fmt.Errorf("failed to submit TX: %w", err)
	}

	return nil
}
