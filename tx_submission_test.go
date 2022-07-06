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
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestClient_SubmitTx(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.SkipNow()
	}

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	err := client.SubmitTx(ctx, []byte("blah"))
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestSubmitTxResult(t *testing.T) {
	err := filepath.Walk("ext/ogmios/server/test/vectors/TxSubmission", testSubmitTxResult(t))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
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

			data, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			err = readSubmitTx(data)
			var ste SubmitTxError
			if ok := errors.As(err, &ste); ok {
				keys, err := ste.ErrorCodes()
				if err != nil {
					t.Fatalf("got %v; want nil", err)
				}
				if len(keys) == 0 {
					t.Fatalf("got 0 keys; want > 0")
				}
				fmt.Println(keys)
				return
			}
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}
		})

		return nil
	}
}
