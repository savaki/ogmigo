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
	"os"
	"testing"
)

func TestClient_EvaluateTx(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	if endpoint == "" {
		t.Fatalf("ogmios not configured")
	}
	testTx := "84a30081825820aafc6c20740e42caed5f97e527b1aa1a4fb5fc13a16b68f141b88440d403d512010182a300581d707fa2a9a246c648573168390652b61abeae2dc761a66e363e37b2b179011a0053ec60028201d81858b8d8799fd8799f581fa9077b2871760b856883e1f7b668e0590b30943420a9ab6608e521fd8e2425ffd8799f581c035dee66d57cc271697711d63c8c35ffa0b6c4468a6a98024feac73bff1a002625a0d8799fd8799fd8799f581c035dee66d57cc271697711d63c8c35ffa0b6c4468a6a98024feac73bffd87a80ffd87980ffd8799f9f40401a000f4240ff9f581cd441227553a0f1a965fee7d60a0f724b368dd1bddbc208730fccebcf4652424552525900ffffd87980ff82581d60035dee66d57cc271697711d63c8c35ffa0b6c4468a6a98024feac73b1b000000024e76c5c9021a0002a8b1a100818258206f2b757b39b783977e0306bd2751fa3db19f2f8d52d478f44a4a03efb0fa2b295840727e21f70ab833001d5fd0dd496714ebc6c40e5a8f7eedd2d0a3a3c017fde0f7f0cebf77f8a88d77fb06259c1cec617b4d35f473784c06ddc185249683a49c02f5f6"

	ctx := context.Background()
	client := New(WithEndpoint(endpoint), WithLogger(DefaultLogger))
	units, err := client.EvaluateTx(ctx, testTx)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("%v\n", units)
}
