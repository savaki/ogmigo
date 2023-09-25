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
	"bufio" // REMOVE LATER?
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshalV6(t *testing.T) {
	err := filepath.Walk("../../ext/ogmios/server/test/vectors/NextBlockResponse", assertStructMatchesSchemaV6(t))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	decoder := json.NewDecoder(nil)
	decoder.DisallowUnknownFields()
}

func find(s string, d map[string]interface{}) string {
	r := ""
	for k, v := range d {
		if k == s {
			r += fmt.Sprintf("%s", v) + ";"
		}
		if c, ok := v.(map[string]interface{}); ok {
			r += find(s, c)
		}
	}
	return r
}

func assertStructMatchesSchemaV6(t *testing.T) filepath.WalkFunc {
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

		d := map[string]interface{}{}
		stat, err := f.Stat()
		bs := make([]byte, stat.Size())
		_, err = bufio.NewReader(f).Read(bs)
		json.Unmarshal(bs, &d)
		blockEra := find("era", d)
		blockType := find("type", d)
		t.Logf("DEBUG - %s: %s\n", "era", blockEra)
		t.Logf("DEBUG - %s: %s\n", "type", blockType)
		f.Seek(0, io.SeekStart)
		decoder := json.NewDecoder(f)
		decoder.DisallowUnknownFields()
		if blockEra == "Byron" && blockType == "EBB" {
			err = decoder.Decode(&ResponseByronEBBV6{})
		} else if blockEra == "Byron" && blockType == "BFT" {
			err = decoder.Decode(&ResponseByronBFTV6{})
		} else {
			err = decoder.Decode(&ResponsePraosV6{})
		}

		if err != nil {
			t.Fatalf("got %v; want nil: %v", err, fmt.Sprintf("struct did not match schema for file, %v", path))
		}

		return nil
	}
}

func TestByronBFT(t *testing.T) {
	data := `{
		"jsonrpc": "2.0",
		"method": "nextBlock",
		"result": {
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
		},
		"id": null
	}
`
	var response ResponseByronBFTV6
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
}

func TestByronEBB(t *testing.T) {
	data := `{
		"jsonrpc": "2.0",
		"method": "nextBlock",
		"result": {
			"direction": "forward",
			"block": {
				"type": "ebb",
				"era": "byron",
				"height": 13521870305663481883,
				"id": "88cb33510fb3f8b878f44523163b00b619d82957d28726230f7214a36b49e2e4",
				"ancestor": "aad78a13b50a014a24633c7d44fd8f8d18f67bbb3fa9cbcedf834ac899759dcd"
			},
			"tip": {
				"slot": 58655,
				"id": "d15959c51e4e83b791a56fac47ecb0b7dfc7a65e7c062d8ea3ae97f6635db1dc",
				"height": 1038970
			}
		},
		"id": null
	}
`
	var response ResponseByronEBBV6
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
}

func TestPraos(t *testing.T) {
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
	var response ResponsePraosV6
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
}
