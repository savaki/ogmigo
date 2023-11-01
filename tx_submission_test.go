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
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/tj/assert"
)

func TestClient_SubmitTx(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	err := client.SubmitTx(ctx, "00010203")
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSubmitTxResult(t *testing.T) {
	err := filepath.Walk("ext/ogmios/server/test/vectors/SubmitTransactionResponse", testSubmitTxResult(t))
	assert.Nil(t, err)
}

func testSubmitTxResult(t *testing.T) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		t.Run(filepath.Base(path), func(t *testing.T) {
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}
			if info.IsDir() {
				return
			}

			_, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}
		})

		return nil
	}
}
