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

package ogmios

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"github.com/SundaeSwap-finance/sundae-sync/ouroboros/chainsync"
	"github.com/tj/assert"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestClient_ReadNext(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	var (
		ctx       = context.Background()
		logger, _ = zap.NewDevelopment()
		p         = message.NewPrinter(language.English)
		counter   int64
		read      int64
	)

	client, err := New(ctx, logger, endpoint, 50)
	assert.Nil(t, err)
	defer client.Close()

	for {
		data, err := client.ReadNext(ctx)
		assert.Nil(t, err)

		var response chainsync.Response
		decoder := json.NewDecoder(bytes.NewReader(data)) // use decoder to check for unknown fields
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&response)
		if err != nil {
			fmt.Println(string(data))
		}
		assert.Nil(t, err)

		read += int64(len(data))
		if v := atomic.AddInt64(&counter, 1); v%1e3 == 0 {
			var blockNo uint64
			if response.Result != nil && response.Result.RollForward != nil {
				if ps, ok := response.Result.RollForward.Tip.PointStruct(); ok {
					blockNo = ps.BlockNo
				}
			}
			logger.Info("read",
				zap.Uint64("block", blockNo),
				zap.String("n", p.Sprintf("%d", v)),
				zap.String("read", p.Sprintf("%d", read)),
			)
		}
	}
}
