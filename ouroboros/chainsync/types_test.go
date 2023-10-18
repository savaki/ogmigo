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

package chainsync

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fxamacker/cbor/v2"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	err := filepath.Walk("../../ext/ogmios/server/test/vectors/NextBlockResponse", assertStructMatchesSchema(t))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	decoder := json.NewDecoder(nil)
	decoder.DisallowUnknownFields()
}

func assertStructMatchesSchema(t *testing.T) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		path, _ = filepath.Abs(path)
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		defer f.Close()

		decoder := json.NewDecoder(f)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&ResponsePraos{})
		if err != nil {
			t.Fatalf("got %v; want nil: %v", err, fmt.Sprintf("struct did not match schema for file, %v", path))
		}

		return nil
	}
}

func TestDynamodbSerialize(t *testing.T) {
	err := filepath.Walk("../../ext/ogmios/server/test/vectors/NextBlockResponse", assertDynamoDBSerialize(t))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
}

// TODO - This assumes non-Byron blocks. We're not technically supporting Byron in v6.
// Rework this test to ignore Byron blocks?
func assertDynamoDBSerialize(t *testing.T) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		path, _ = filepath.Abs(path)
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		defer f.Close()

		var want ResponsePraos
		decoder := json.NewDecoder(f)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&want)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		item, err := dynamodbattribute.Marshal(want)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		var got ResponsePraos
		err = dynamodbattribute.Unmarshal(item, &got)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		w, err := json.Marshal(want)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		g, err := json.Marshal(got)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		opts := jsondiff.DefaultConsoleOptions()
		diff, s := jsondiff.Compare(w, g, &opts)
		if diff == jsondiff.FullMatch {
			return nil
		}

		if got, want := diff, jsondiff.FullMatch; !reflect.DeepEqual(got, want) {
			fmt.Println(s)
			t.Fatalf("got %#v; want %#v", got, want)
		}

		return nil
	}
}

func TestPoint_CBOR(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := PointString("origin")
		item, err := cbor.Marshal(want.Point())
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		var point Point
		err = cbor.Unmarshal(item, &point)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := point.PointType(), PointTypeString; got != want {
			t.Fatalf("got %v; want %v", got, want)
		}

		got, ok := point.PointString()
		if !ok {
			t.Fatalf("got false; want true")
		}
		if got != want {
			t.Fatalf("got %v; want %v", got, want)
		}
	})

	t.Run("struct", func(t *testing.T) {
		want := &PointStruct{
			BlockNo: 123,
			ID:      "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
			Slot:    456,
		}
		item, err := cbor.Marshal(want.Point())
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		var point Point
		err = cbor.Unmarshal(item, &point)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := point.PointType(), PointTypeStruct; got != want {
			t.Fatalf("got %v; want %v", got, want)
		}

		got, ok := point.PointStruct()
		if !ok {
			t.Fatalf("got false; want true")
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %#v; want %#v", got, want)
		}
	})
}

func TestPoint_DynamoDB(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := PointString("origin")
		item, err := dynamodbattribute.Marshal(want.Point())
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		var point Point
		err = dynamodbattribute.Unmarshal(item, &point)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := point.PointType(), PointTypeString; got != want {
			t.Fatalf("got %v; want %v", got, want)
		}

		got, ok := point.PointString()
		if !ok {
			t.Fatalf("got false; want true")
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v; want %v", got, want)
		}
	})

	t.Run("struct", func(t *testing.T) {
		want := &PointStruct{
			BlockNo: 123,
			ID:      "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
			Slot:    456,
		}
		item, err := dynamodbattribute.Marshal(want.Point())
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		var point Point
		err = dynamodbattribute.Unmarshal(item, &point)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := point.PointType(), PointTypeStruct; got != want {
			t.Fatalf("got %v; want %v", got, want)
		}

		got, ok := point.PointStruct()
		if !ok {
			t.Fatalf("got false; want true")
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v; want %v", got, want)
		}
	})
}

func TestPoint_JSON(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := PointString("origin")
		data, err := json.Marshal(want.Point())
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}

		var point Point
		err = json.Unmarshal(data, &point)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := point.PointType(), PointTypeString; got != want {
			t.Fatalf("got %v; want %v", got, want)
		}

		got, ok := point.PointString()
		if !ok {
			t.Fatalf("got false; want true")
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v; want %v", got, want)
		}
	})

	t.Run("struct", func(t *testing.T) {
		want := &PointStruct{
			BlockNo: 123,
			ID:      "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
			Slot:    456,
		}
		data, err := json.Marshal(want.Point())
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		var point Point
		err = json.Unmarshal(data, &point)
		if err != nil {
			t.Fatalf("got %v; want nil", err)
		}
		if got, want := point.PointType(), PointTypeStruct; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v; want %v", got, want)
		}

		got, ok := point.PointStruct()
		if !ok {
			t.Fatalf("got false; want true")
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v; want %v", got, want)
		}
	})
}

func TestTxID_Index(t *testing.T) {
	if got, want := TxID("a#3").Index(), 3; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func TestTxID_TxHash(t *testing.T) {
	if got, want := TxID("a#3").TxHash(), "a"; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func TestPoints_Sort(t *testing.T) {
	s1 := PointString("1").Point()
	s2 := PointString("2").Point()
	p1 := PointStruct{Slot: 10}.Point()
	p2 := PointStruct{Slot: 10}.Point()
	tests := map[string]struct {
		Input Points
		Want  Points
	}{
		"string": {
			Input: Points{s1, s2},
			Want:  Points{s2, s1},
		},
		"points": {
			Input: Points{p1, p2},
			Want:  Points{p2, p1},
		},
		"mixed": {
			Input: Points{s1, p1, s2, p2},
			Want:  Points{p2, p1, s2, s1},
		},
	}
	for label, tc := range tests {
		t.Run(label, func(t *testing.T) {
			got := tc.Input
			sort.Sort(got)
			if !reflect.DeepEqual(got, tc.Want) {
				t.Fatalf("got %#v; want %#v", got, tc.Want)
			}
		})
	}
}

func TestPraosResponse(t *testing.T) {
	data := `{
		"jsonrpc": "2.0",
		"method": "nextBlock",
		"result": {
			"direction": "forward",
			"block": {
				"type": "praos",
				"era": "allegra",
				"id": "109dab4b2e94ebb0e1ad5bafa335392ab6d4e9ea01d7878087f825ad4a9cdeb3",
				"ancestor": "82a795bf5dcba27286949627170b40e765e3e1da7d90d20a9bf1075d9517c747",
				"nonce": {
					"output": "02020202000100020001010000010100010201000202010102000101010000000202000200020000000100020201010100020002010201010000010102000100",
					"proof": "98fbb316e8bbd717ca27325d44eaaf90867b8351a7b9839a1b19e8e6b416791363470f0709ba78f6a5c4f50dee2e6295aa04a4fd125143b0dfdbed7422b6efe49ad4b94c34eb7f434c486e1d803e0b0a"
				},
				"height": 0,
				"slot": 2,
				"issuer": {
					"verificationKey": "948b49ced8f0316e7d33a8de788eb3ba3ea88deb2716bf185cae5820591fdabf",
					"vrfVerificationKey": "224655a70fb7d8cb8f900c97f16cfc45645124a1d388b6eaa2c4511ef5da79d8",
					"leaderValue": {
						"output": "00000102020200010101020101010102020202000001000100010001010001020101000202020200000201000001010002020200010202000102010000020000",
						"proof": "a2876927ad069b694d87cf070e49f7786ba1312b25dcc501657b545a1297fe0db1627d96d6a3ca7bbf4e8ede19885e4a164bd8d4e31d747fcd05ac63f17af4a58503a6417ce9554724f30b7dada17f05"
					},
					"operationalCertificate": {
						"count": 0,
						"kes": {
							"period": 2,
							"verificationKey": "b39372534cb0636dfdfcab8e7726aa7bb3160ef9b640a52bb0a043093baa613d"
						}
					}
				},
				"protocol": {
					"version": {
						"major": 4,
						"minor": 1
					}
				},
				"size": {
					"bytes": 1
				},
				"transactions": [
					{
						"id": "2bd3c9bb48d587f9a6bccb8230110ba0c536aa05c5966b1a63052b26775db4f4",
						"spends": "inputs",
						"inputs": [
							{
								"transaction": {
									"id": "7525e434d3b951d3499b87f8cbd02ed0ebb0c0248343af568c29d15e79825bd8"
								},
								"index": 0
							}
						],
						"outputs": [
							{
								"address": "addr_test1ypt3trk84kyekn6zy936rz4aetlndlhkv35ztmfv57te3uzsahhjajdzftrmkmtfnrs5ryefgug64maldke5qjd0yvcs2v68ny",
								"value": {
									"ada": {
										"lovelace": 987617
									}
								}
							},
							{
								"address": "EqGAuA8vHnPEVahYhLV7WYBL4EeaZuTMttq6h5Xc4hHvYpUyyMu1EcdXNKf85Xb4G3hGswH8tqbY6pGhunn1yKrxmfV8aDW4sFfks1ruLM2icn93HRDrwra",
								"value": {
									"ada": {
										"lovelace": 66410
									}
								}
							}
						],
						"withdrawals": {
							"stake_test17re7uuh9zercugh782esdhp4lspnz6k7sspk7254yagzpwq83d45f": {
								"lovelace": 37534
							}
						},
						"fee": {
							"lovelace": 737494
						},
						"validityInterval": {
							"invalidBefore": 1,
							"invalidAfter": 1
						},
						"signatories": [
							{
								"key": "634757c7aecfd4d9b6ba22b68d6410484e6427ac05e2e605e633f19617d1cecb",
								"signature": "7a9aafd169b23153442311f068b3d3ae782d9a7c95f60694453042b4e3725697ca2585270f5e6aa53b914c6f2c79ae363455e326f527bd5fcda627fb3a99dcdc",
								"chainCode": "11de49",
								"addressAttributes": "9068"
							},
							{
								"key": "d3688e9746e6fa68a20fa7a376a947b1adb53b44da72b02eb6d5a00f9bf34fc5",
								"signature": "b16448d177014e269064b8ebc2e6b192ed03b7978fe4278d68c4716a68840664d005aa649e675f0569ef8906df21a5514c00eea3df676bcae4f2d1bad9f4f3b6",
								"chainCode": "07",
								"addressAttributes": "53"
							},
							{
								"key": "7658aec9e30d35256940dab38397becb6712c04f4d79596b0c99dbb6d00892da",
								"signature": "c6215f8eacdbcf09cbb440f5a1bf2553344873b1fd97175c06ccb3920546ee4f40523b58ae2dee86cf526771cc995098b2e015c5af54af36d5c5975170801cbc",
								"chainCode": "a7",
								"addressAttributes": "a993bd"
							}
						],
						"scripts": {
							"5c6dfd90190d6c66547da3debf3d8967340b270feb3f6cf5a5f8cab8": {
								"language": "native",
								"json": {
									"clause": "after",
									"slot": 2
								}
							},
							"9db4dcc531262ba3db77abdb21315f4f68b1d460a74bd87a9ce50ef7": {
								"language": "native",
								"json": {
									"clause": "some",
									"atLeast": 4,
									"from": [
										{
											"clause": "any",
											"from": [
												{
													"clause": "all",
													"from": [
														{
															"clause": "signature",
															"from": "b5ae663aaea8e500157bdf4baafd6f5ba0ce5759f7cd4101fc132f54"
														}
													]
												},
												{
													"clause": "after",
													"slot": 8
												},
												{
													"clause": "before",
													"slot": 16
												}
											]
										},
										{
											"clause": "after",
											"slot": 14
										},
										{
											"clause": "any",
											"from": [
												{
													"clause": "signature",
													"from": "65fc709a5e019b8aba76f6977c1c8770e4b36fa76f434efc588747b7"
												}
											]
										},
										{
											"clause": "all",
											"from": [
												{
													"clause": "signature",
													"from": "0d94e174732ef9aae73f395ab44507bfa983d65023c11a951f0c32e4"
												},
												{
													"clause": "before",
													"slot": 6
												},
												{
													"clause": "some",
													"atLeast": 2,
													"from": [
														{
															"clause": "signature",
															"from": "3542acb3a64d80c29302260d62c3b87a742ad14abf855ebc6733081e"
														},
														{
															"clause": "signature",
															"from": "a646474b8f5431261506b6c273d307c7569a4eb6c96b42dd4a29520a"
														},
														{
															"clause": "signature",
															"from": "4acf2773917c7b547c576a7ff110d2ba5733c1f1ca9cdc659aea3a56"
														}
													]
												},
												{
													"clause": "before",
													"slot": 0
												},
												{
													"clause": "any",
													"from": [
														{
															"clause": "signature",
															"from": "a646474b8f5431261506b6c273d307c7569a4eb6c96b42dd4a29520a"
														}
													]
												}
											]
										},
										{
											"clause": "signature",
											"from": "e0a714319812c3f773ba04ec5d6b3ffcd5aad85006805b047b082541"
										}
									]
								}
							},
							"acfae5054c2216d7b6e1985637d35fbe4f0520bc84f0e8c60b9a9789": {
								"language": "native",
								"json": {
									"clause": "some",
									"atLeast": 1,
									"from": [
										{
											"clause": "before",
											"slot": 16
										}
									]
								}
							}
						}
					}
				]
			},
			"tip": {
				"slot": 47137,
				"id": "12a3b5451db4b82932f1e4045e1be8a829bf8d30f61010cfac42781f29a47564",
				"height": 8953265
			}
		},
		"id": "H07GyFhAxTb4"
	}
`
	var response ResponsePraos
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
}

func TestVasil_DatumParsing_Base64(t *testing.T) {
	data := `{"datums": {"a": "2HmfWBzIboNaGwk6qBYQ/Tk19GPOUpkpze2Ldfe1HOZEQpwK/w=="}}`
	var response Witness
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	datumHex := response.Datums["a"]
	_, err = hex.DecodeString(datumHex)
	if err != nil {
		t.Fatalf("error decoding hex string: %v", err)
	}
}

func TestVasil_DatumParsing_Hex(t *testing.T) {
	data := `{"datums": {"a": "d8799f581cc86e835a1b093aa81610fd3935f463ce529929cded8b75f7b51ce644429c0aff"}}`
	var response Witness
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	datumHex := response.Datums["a"]
	_, err = hex.DecodeString(datumHex)
	if err != nil {
		t.Fatalf("error decoding hex string: %v", err)
	}
}

func TestVasil_BackwardsCompatibleWithExistingDynamoDB(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/scoop.json")
	assert.Nil(t, err)

	var item map[string]*dynamodb.AttributeValue
	err = json.Unmarshal(data, &item)
	assert.NoError(t, err)

	var response Tx
	err = dynamodbattribute.Unmarshal(item["tx"], &response)
	assert.NoError(t, err)
	fmt.Println(response.Datums)
}
