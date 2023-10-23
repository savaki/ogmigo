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

func TestCompatibleResult(t *testing.T) {
	dataRequestNext := `{
		"RollForward": {
			"block": {
				"babbage": {
					"body": [
						{
							"witness": {
								"signatures": {
									"06f398b7fdd5a32d3de808fb1d3b32be177f2d2878e6b19845aeaabf8b94661c": "huZkCWpiS7YvXEmnfr5QNrvD9Lfe5hXOO3uFdXb+GtzO4DsnD9P28hvsowrgKBs5g6C0M5hHfNWTzUfcUm25/g==",
									"b420fafa7256edf97efa07c31f28fd072b528b81040cc8585b4a261337b2cde0": "HMGGGy/U8v+LQy+yunqz0BtZv2CQoHWfCi2KeZCXCDZ2nKtt1DTJEQz9Hdb0affxpMIcV4560Uw9I67ZRMVR/w=="
								},
								"scripts": {
									"6988f97425839f8caa8757e054a28667567b59c19dd38c3f280f25e9": {
										"native": {
											"startsAt": 46090
										}
									},
									"a37f946a84f2ec6d21929be44eb0bca522c560ef17094ecccfe2f6d1": {
										"native": {
											"startsAt": 79106
										}
									}
								},
								"datums": {},
								"redeemers": {
									"certificate:0": {
										"redeemer": "d8668219032d80",
										"executionUnits": {
											"memory": 2009715596563484410,
											"steps": 6764152604174301894
										}
									}
								},
								"bootstrap": []
							},
							"raw": "hKsAgYJYILswpCweYvCv2l8KTopWL3oTokzqAO6BkXuGuJ6AExSqAg2Bglgg7hVazpxAKSB0y2r/jJzN0nPIFkj/EUnvNrzqbruKPiUCAYAQowBYOSANlOF0cy75quc/OVq0RQe/qYPWUCPBGpUfDDLkDZThdHMu+arnPzlatEUHv6mD1lAjwRqVHwwy5AEAAoIAWCCuhdJFo9AL/eAfWfPE/gtL+uHLN+nPkZKerc6kmFcR3gIZAr0DGdR/BIGKA1gc4KcUMZgSw/dzugTsXWs//NWq2FAGgFsEewglQVggAmi+nb0ERuqiF+HeyPOZJJMF5VHX/BQ33YRSH3SqYhwZA+AY3NgeggECWB3h4KcUMZgSw/dzugTsXWs//NWq2FAGgFsEewglQYCA9gaCoBlFbggZ1/cOgVgcDZThdHMu+arnPzlatEUHv6mD1lAjwRqVHwwy5AdYIAMXCi51l7e349hMBTkdE5pisVfnh4bYwILync9MERMUowCCglggBvOYt/3Voy096Aj7HTsyvhd/LSh45rGYRa6qv4uUZhxYQIbmZAlqYku2L1xJp36+UDa7w/S33uYVzjt7hXV2/hrczuA7Jw/T9vIb7KMK4CgbOYOgtDOYR3zVk81H3FJtuf6CWCC0IPr6clbt+X76B8MfKP0HK1KLgQQMyFhbSiYTN7LN4FhAHMGGGy/U8v+LQy+yunqz0BtZv2CQoHWfCi2KeZCXCDZ2nKtt1DTJEQz9Hdb0affxpMIcV4560Uw9I67ZRMVR/wGCggQZtAqCBBoAATUCBYGEAgDYZoIZAy2Aghsb4/Gv6LqO+htd3xj0mle6xvXZAQOhAKIAYlA3AgA=",
							"id": "8bc0405aa28ef6c43883e7e96600bbffc1d20e932607c7b1be514296dd025aec",
							"body": {
								"inputs": [
									{
										"txId": "bb30a42c1e62f0afda5f0a4e8a562f7a13a24cea00ee81917b86b89e801314aa",
										"index": 2
									}
								],
								"collaterals": [
									{
										"txId": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25",
										"index": 2
									}
								],
								"references": [],
								"collateralReturn": {
									"address": "addr_test1yqxefct5wvh0n2h88uu44dz9q7l6nq7k2q3uzx54ruxr9eqdjnshguewlx4ww0eet26y2pal4xpav5prcydf28cvxtjqzdy9kg",
									"value": {
										"coins": 0,
										"assets": {}
									},
									"datumHash": "ae85d245a3d00bfde01f59f3c4fe0b4bfae1cb37e9cf91929eadcea4985711de",
									"datum": null,
									"script": null
								},
								"totalCollateral": null,
								"outputs": [],
								"certificates": [
									{
										"poolRegistration": {
											"id": "pool1uzn3gvvcztplwua6qnk966elln264kzsq6q9kprmpqj5zytzn03",
											"vrf": "0268be9dbd0446eaa217e1dec8f399249305e551d7fc1437dd84521f74aa621c",
											"pledge": 992,
											"cost": 220,
											"margin": "1/2",
											"rewardAccount": "stake1u8s2w9p3nqfv8amnhgzwchtt8l7dt2kc2qrgqkcy0vyz2sgltd8xl",
											"owners": [],
											"relays": [],
											"metadata": null
										}
									}
								],
								"withdrawals": {},
								"fee": 701,
								"validityInterval": {
									"invalidBefore": 55287,
									"invalidHereafter": 54399
								},
								"update": {
									"proposal": {},
									"epoch": 17774
								},
								"mint": {
									"coins": 0,
									"assets": {}
								},
								"network": null,
								"scriptIntegrityHash": null,
								"requiredExtraSignatures": [
									"0d94e174732ef9aae73f395ab44507bfa983d65023c11a951f0c32e4"
								]
							},
							"metadata": {
								"hash": "03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314",
								"body": {
									"blob": {
										"0": {
											"string": "P7"
										},
										"2": {
											"int": 0
										}
									},
									"scripts": []
								}
							},
							"inputSource": "inputs"
						}
					],
					"header": {
						"protocolVersion": {
							"major": 345,
							"minor": 373
						},
						"opCert": {
							"hotVk": "75eWZ631+ApZTU7CILYZgCQKkhchDH7PaMdyqcrYA8E=",
							"count": 0,
							"kesPeriod": 2,
							"sigma": "dTXoCl4HOLeieRvHy24lCNYIUVW135+iokEMbfMJu6zmdbxW9Swg0LUFD2yy86qu1p0BSjJWNvckLuI8XOGwnw=="
						},
						"signature": "yK1nKo9qm3oLLVJMNrSl/gtN1qmk/7hKVj+/b3w+uOEPUrtw3vKSV+uwRx5fh61NxocNt+ISF/6yeTgfibazB45TE5kdkDkz78sLb4Rog1Z0I5JDqoli43O4aCtO9LwezTRzxJ+94d4W9sRsZCOatJUO2ewDHLZpIGmWHrDjw18cz0eUwXCGoI8nHbFhLiVCpXZSgv6ABVOx0G2b2dldlh7ykZ5m6aWWs7ELFtSGyDIuAuMQeue64G94+UnCxQE1LpWMKAma8ZGY1ktA9naU+ow5cTOmno/Tyqqgi1/sjX/cwC5lm3BmoF0lHod/Qq/pEQGx97WMJ4IjBaZ71vqPeMUnf5tVMRW/xGxhgfutZTK+FA3vNFcT0NzIuUADlDiK7aUidWM2KnImd6mV1UO+NkmSTjf4PrfoQnCJlY1XVIbGN/bMerDECHCHAXXdH9+0YH4tPOU1ZHI9XnauacBT2/AN8RUAgIEHGoMjnsShg/oXtgmrS6O7OIaTzyMXdcQAjCNaB9ewRU+LmrXfyTivj6Qm8rng7s/ZUgAlVtPZd2/dkumi8+7Krw2uobTqER+ZGh6Gbthi2OcIjX7VGCPABg==",
						"vrfInput": {
							"output": "m8rQn8CgZA51x/OeXvKOA1oSYh2wtU6EVm0DhIxce9wiaFkkaGqcL/8cIzwHXGAG0TDqM+fXt7kmCLWmQ64UUg==",
							"proof": "7JFp5EqDQARMBnRBzTEFx3jVdl7lcVIAE86jOpI9N8tdbpi7x9EYxNwzlRg/1fmXWcgyK7YVUCS9N/wPSk6Aa6s+AGeyZBuiocg4lB/K5gs="
						},
						"blockHeight": 5,
						"slot": 7,
						"prevHash": "genesis",
						"issuerVk": "97211ca376ccc3a547fb579251879690f2f91911b84496e745f5c873438d653b",
						"issuerVrf": "syNCTuYbTRjsWXbdnEsdufWPNBzUl8rgpa9ydwxjh7M=",
						"blockSize": 2,
						"blockHash": "0268be9dbd0446eaa217e1dec8f399249305e551d7fc1437dd84521f74aa621c"
					},
					"headerHash": "ae85d245a3d00bfde01f59f3c4fe0b4bfae1cb37e9cf91929eadcea4985711de"
				}
			},
			"tip": {
				"slot": 795,
				"hash": "5230df78c87a3a1726f0083b3fdc0541e423e1a417913ebc7ff22037d1f3ca1c",
				"blockNo": 16367053
			}
		}
	}`

	var method1 CompatibleResultNextBlock
	err := json.Unmarshal([]byte(dataRequestNext), &method1)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	dataNextBlock := `{
        "direction": "forward",
        "block": {
            "type": "bft",
            "era": "byron",
            "id": "21d89eb76817a3ea1b9df033526f833cef18c4c762c356922eafb148bc14001e",
            "ancestor": "aad78a13b50a014a24633c7d44fd8f8d18f67bbb3fa9cbcedf834ac899759dcd",
            "height": 6574977808651210019,
            "slot": 3294403856197716808,
            "size": {
                "bytes": 921
            },
            "transactions": [],
            "operationalCertificates": [],
            "protocol": {
                "id": 100,
                "version": {
                    "major": 892,
                    "minor": 51768,
                    "patch": 35
                },
                "software": {
                    "appName": "sDSF7",
                    "number": 4117586641
                },
                "update": {
                    "proposal": {
                        "version": {
                            "major": 29680,
                            "minor": 55074,
                            "patch": 14
                        },
                        "software": {
                            "appName": "",
                            "number": 1115902224
                        },
                        "parameters": {
                            "scriptVersion": 25799,
                            "slotDuration": 66450280671243551,
                            "maxBlockBodySize": {
                                "bytes": 6332592073220457
                            },
                            "maxTransactionSize": {
                                "bytes": 81061672327968174
                            },
                            "maxUpdateProposalSize": {
                                "bytes": 172822501139406362
                            },
                            "multiPartyComputationThreshold": "288921139725109/500000000000000",
                            "heavyDelegationThreshold": "138029086789251/200000000000000",
                            "updateProposalThreshold": "294859239649641/500000000000000",
                            "updateProposalTimeToLive": 16087585201267149304,
                            "unlockStakeEpoch": 2481827810661368817,
                            "softForkInitThreshold": "187743902107349/200000000000000",
                            "softForkMinThreshold": "4036481357189/500000000000000",
                            "softForkDecrementThreshold": "13675322648881/40000000000000"
                        },
                        "metadata": {}
                    },
                    "votes": []
                }
            },
            "issuer": {
                "verificationKey": "8ece824656e8007745e655f080ced6a86012a26cd667232e0b810b1dd2d6555cb8a2405a55c46ce9fb0346bffb6cd7a8cd294d8950eeba7f6ba9c4fe26256b24"
            },
            "delegate": {
                "verificationKey": "8ece824656e8007745e655f080ced6a86012a26cd667232e0b810b1dd2d6555cb8a2405a55c46ce9fb0346bffb6cd7a8cd294d8950eeba7f6ba9c4fe26256b24"
            }
        },
        "tip": {
            "slot": 46734,
            "id": "4cbfd2af37df07a6bab5850b8122ed56708b4b57676cd23a6b564d33907a52e9",
            "height": 3642099
        }
	}`

	var method2 CompatibleResultNextBlock
	err = json.Unmarshal([]byte(dataNextBlock), &method2)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	dataFindIntersect := `{
        "IntersectionFound": {
            "point": "origin",
            "tip": {
                "slot": 62344,
                "hash": "2208e439244a1d0ef238352e3693098aba9de9dd0154f9056551636c8ed15dc1",
                "blockNo": 2
            }
        }
	}`

	var method3 CompatibleResultFindIntersection
	err = json.Unmarshal([]byte(dataFindIntersect), &method3)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

	dataFindIntersection := `{
        "intersection": {
            "slot": 71768,
            "id": "990aad500baaf649562318c3e23f96207a4fc41004902adf02716c9fa4b13827"
        },
        "tip": {
            "slot": 24383,
            "id": "aad78a13b50a014a24633c7d44fd8f8d18f67bbb3fa9cbcedf834ac899759dcd",
            "height": 1
        }
	}`

	var method4 CompatibleResultNextBlock
	err = json.Unmarshal([]byte(dataFindIntersection), &method4)
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
