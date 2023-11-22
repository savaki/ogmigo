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

package compatibility

import (
	"encoding/json"
	"testing"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/shared"
	"github.com/tj/assert"
)

// TODO:
//   - break this up into smaller tests;
//   - put strings in a test_data folder;
//   - check the data we get to make sure it round tripped correctly
//   - test dynamodb marshalling
func TestCompatibleResult(t *testing.T) {
	dataRequestNextForward := `{
        "RollForward": {
            "block": {
                "babbage": {
                    "body": [
                        {
                            "witness": {
                                "signatures": {
                                    "d7ca16cc9549ae4f4fad53ab776a2dd2af25fa88de648d5ad7a18b106ab0bafc": "1u/Zfip3DrebuM61lZIPaqtEW3qLRczzPx8yAWyl8HYg0nGJrJPCdfItftRSG9GR+BZsWXMAgdq16qfY5Q7g0A=="
                                },
                                "scripts": {},
                                "datums": {},
                                "redeemers": {
                                    "spend:3": {
                                        "redeemer": "a2d8668218d79f43f800b8d9051b9f04240121ffa142cbd204d866821903939f44abbff2b12105418e04ffd87e80ff9f9f2202ff9f4000ff9f04ffd8668218b29f0240ffff809f9f24234242c8ffd8668218d59f020520ffd8668219013c9f242340ffd866821901b89f22ffff",
                                        "executionUnits": {
                                            "memory": 3435122052805974510,
                                            "steps": 4325896289790197942
                                        }
                                    },
                                    "certificate:1": {
                                        "redeemer": "44735040e3",
                                        "executionUnits": {
                                            "memory": 5866473676942577327,
                                            "steps": 1473688137007832935
                                        }
                                    },
                                    "withdrawal:3": {
                                        "redeemer": "a1414ed866821902359fa143a3c99a20427b6cff",
                                        "executionUnits": {
                                            "memory": 3631115493909780867,
                                            "steps": 1043417232138852536
                                        }
                                    }
                                },
                                "bootstrap": [
                                    {
                                        "key": "7b5cb6d898502a692d256255efd9155053fa4ad247c3f74d152a67884e2c2424",
                                        "chainCode": "ac",
                                        "addressAttributes": "/Q==",
                                        "signature": "dlCn8VbOG+kDJJ7fkp9tdb/Gh/BhPTMZHVG9TresMvQeGJuYYN1J52Qjq9OTXbsCzppSVhQm/65zIBSfr+NeZQ=="
                                    },
                                    {
                                        "key": "16d35a2dffb176e89d67ac7e7fca602f1deb5059f2edf48673108d72478f62e5",
                                        "chainCode": "909d",
                                        "addressAttributes": null,
                                        "signature": "5i5QfbSJhuQnzN2kCxCN9tq4TOmg9nN3uSpx5f8F2wWqZCWo/9u/Xu92HowF8H2f7Y0jdvXHdRI/VfWV6z+rZg=="
                                    },
                                    {
                                        "key": "aa54eddec8234ca7a0205fd1262c0a7fc7e8ce56cb0edaba5ca2b66f58d8d8cf",
                                        "chainCode": "7c",
                                        "addressAttributes": "qsQ=",
                                        "signature": "CyCsYOrCpmGzWrq9VNFyLQwnAlNaU7UDCgjxK+LOoQlBdSlyYdn58bH87XntrcFHDk4TzHXV+RVyLWhW/KrGqQ=="
                                    }
                                ]
                            },
                            "raw": "hKgAgYJYIJXDADp4WF4NuMlJb23u9N4P8ACZS4U0zWbU/pa7Id3TAhKBglgg7hVazpxAKSB0y2r/jJzN0nPIFkj/EUnvNrzqbruKPiUAAYGjAFgdcDVCrLOmTYDCkwImDWLDuHp0KtFKv4VevGczCB4BggChWBy1rmY6rqjlABV730uq/W9boM5XWffNQQH8Ey9UoUEBBAKCAFgguzCkLB5i8K/aXwpOilYvehOiTOoA7oGRe4a4noATFKoCGQK1CBoAAXUlDoJYHA2U4XRzLvmq5z85WrRFB7+pg9ZQI8EalR8MMuRYHDVCrLOmTYDCkwImDWLDuHp0KtFKv4VevGczCB4JolgcNUKss6ZNgMKTAiYNYsO4enQq0Uq/hV68ZzMIHqFBAANYHOCnFDGYEsP3c7oE7F1rP/zVqthQBoBbBHsIJUGiQARBACIHWCDoi9dXrVub7fNy2NPwz2yWKkadthomX2QY4f/thtop7KMAgYJYINfKFsyVSa5PT61Tq3dqLdKvJfqI3mSNWtehixBqsLr8WEDW79l+KncOt5u4zrWVkg9qq0RbeotFzPM/HzIBbKXwdiDScYmsk8J18i1+1FIb0ZH4FmxZcwCB2rXqp9jlDuDQAoOEWCB7XLbYmFAqaS0lYlXv2RVQU/pK0kfD900VKmeITiwkJFhAdlCn8VbOG+kDJJ7fkp9tdb/Gh/BhPTMZHVG9TresMvQeGJuYYN1J52Qjq9OTXbsCzppSVhQm/65zIBSfr+NeZUGsQf2EWCAW01ot/7F26J1nrH5/ymAvHetQWfLt9IZzEI1yR49i5VhA5i5QfbSJhuQnzN2kCxCN9tq4TOmg9nN3uSpx5f8F2wWqZCWo/9u/Xu92HowF8H2f7Y0jdvXHdRI/VfWV6z+rZkKQnUCEWCCqVO3eyCNMp6AgX9EmLAp/x+jOVssO2rpcorZvWNjYz1hACyCsYOrCpmGzWrq9VNFyLQwnAlNaU7UDCgjxK+LOoQlBdSlyYdn58bH87XntrcFHDk4TzHXV+RVyLWhW/KrGqUF8QqrEBYOEAAOi2GaCGNefQ/gAuNkFG58EJAEh/6FCy9IE2GaCGQOTn0Srv/KxIQVBjgT/2H6A/5+fIgL/n0AA/58E/9hmghiynwJA//+An58kI0JCyP/YZoIY1Z8CBSD/2GaCGQE8nyQjQP/YZoIZAbifIv//ghsvrAFOlZxZ7hs8CKvBiCAAtoQCAURzUEDjghtRaea6LleOrxsUc5eNfA8fZ4QDA6FBTthmghkCNZ+hQ6PJmiBCe2z/ghsyZFBNoOjZgxsOevZxwihYuPX2",
                            "id": "4e1b832f07856a1baf99c35e4e5309f0f07a112200e933ff0082f5fc808bb8e0",
                            "body": {
                                "inputs": [
                                    {
                                        "txId": "95c3003a78585e0db8c9496f6deef4de0ff000994b8534cd66d4fe96bb21ddd3",
                                        "index": 2
                                    }
                                ],
                                "collaterals": [],
                                "references": [
                                    {
                                        "txId": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25",
                                        "index": 0
                                    }
                                ],
                                "collateralReturn": null,
                                "totalCollateral": null,
                                "outputs": [
                                    {
                                        "address": "addr_test1wq659t9n5excps5nqgnq6ckrhpa8g2k3f2lc2h4uvuess8s24hsvh",
                                        "value": {
                                            "coins": 0,
                                            "assets": {
                                                "b5ae663aaea8e500157bdf4baafd6f5ba0ce5759f7cd4101fc132f54.01": 4
                                            }
                                        },
                                        "datumHash": "bb30a42c1e62f0afda5f0a4e8a562f7a13a24cea00ee81917b86b89e801314aa",
                                        "datum": null,
                                        "script": null
                                    }
                                ],
                                "certificates": [],
                                "withdrawals": {},
                                "fee": 693,
                                "validityInterval": {
                                    "invalidBefore": 95525,
                                    "invalidHereafter": null
                                },
                                "update": null,
                                "mint": {
                                    "coins": 0,
                                    "assets": {
                                        "3542acb3a64d80c29302260d62c3b87a742ad14abf855ebc6733081e.00": 3,
                                        "e0a714319812c3f773ba04ec5d6b3ffcd5aad85006805b047b082541": 4,
                                        "e0a714319812c3f773ba04ec5d6b3ffcd5aad85006805b047b082541.00": -3
                                    }
                                },
                                "network": null,
                                "scriptIntegrityHash": null,
                                "requiredExtraSignatures": [
                                    "0d94e174732ef9aae73f395ab44507bfa983d65023c11a951f0c32e4",
                                    "3542acb3a64d80c29302260d62c3b87a742ad14abf855ebc6733081e"
                                ]
                            },
                            "metadata": null,
                            "inputSource": "inputs"
                        },
                        {
                            "witness": {
                                "signatures": {
                                    "f2d2feb39020be1892663bd42a80ca2792e467c438c0773fd1e7dac0dc15427e": "L9ZS6he+/GBUq/Dp0OB7CuvFN3Bxj6BVironQrvs3l/opV/qhkRpVKzPLyPGFojEn/PWSknrQApu9RUptxRaKw=="
                                },
                                "scripts": {
                                    "45c0ad94b0185b6fe2316ef22670a205a448844793ed947d5d6b6e17": {
                                        "native": {
                                            "any": [
                                                {
                                                    "startsAt": 71907
                                                }
                                            ]
                                        }
                                    },
                                    "e705a9fc3a483e27c688e383c7fafc10e2c1fe130ee556698b1f4092": {
                                        "native": {
                                            "any": [
                                                {
                                                    "any": [
                                                        {
                                                            "all": [
                                                                "4acf2773917c7b547c576a7ff110d2ba5733c1f1ca9cdc659aea3a56",
                                                                "b16b56f5ec064be6ac3cab6035efae86b366cc3dc4a0d571603d70e5",
                                                                "76e607db2a31c9a2c32761d2431a186a550cc321f79cd8d6a82b29b8"
                                                            ]
                                                        },
                                                        {
                                                            "any": [
                                                                "76e607db2a31c9a2c32761d2431a186a550cc321f79cd8d6a82b29b8"
                                                            ]
                                                        },
                                                        {
                                                            "expiresAt": 25392
                                                        },
                                                        {
                                                            "1": [
                                                                "e0a714319812c3f773ba04ec5d6b3ffcd5aad85006805b047b082541"
                                                            ]
                                                        }
                                                    ]
                                                },
                                                {
                                                    "startsAt": 33964
                                                },
                                                {
                                                    "startsAt": 20844
                                                }
                                            ]
                                        }
                                    }
                                },
                                "datums": {
                                    "587bd0ffa8cbb6b2543e748c9d3c33aad69cbd9ca6ad0318978b975aa516c9b2": "4343ef46",
                                    "e25d14c7f0f2e5cc7886a455f66a661867a1cc9228cbb2ea31ef464544882e43": "d866821902469fa12021d8668219035c9f43cc6e8a9f4223b1240003ffd866821901e79f0044a023eac743beb99e428fc1ffa142132a44f517b1d080ffff",
                                    "ffc127f938b482eae9caed1fb40e0753bc9064131fd8755256b6be486e0b7c41": "9f239f054480ec10439f4315e81c423d33ff447738c6b140ffd866821902479fd8668219015a9f0304ffd866821901a79f0144d1adf1544422f38b9d4303dcb223ffffd905449f4020ffff"
                                },
                                "redeemers": {
                                    "certificate:4": {
                                        "redeemer": "40",
                                        "executionUnits": {
                                            "memory": 3900831469749680301,
                                            "steps": 4884377380812661489
                                        }
                                    }
                                },
                                "bootstrap": [
                                    {
                                        "key": "63176c55566a62e8bf7f7d93d7d890623fac8b46852501fe9317082619e7ee0a",
                                        "chainCode": null,
                                        "addressAttributes": "2F0=",
                                        "signature": "QuARQ+rceR1IGwT+bDeGkuwBt+U79FMBsTiGrAXYI/BnJ2LJ2Fe3ZPvPGef7qoXbi930ZC59vt77seXiuRBXjQ=="
                                    }
                                ]
                            },
                            "raw": "hK0Ag4JYIJXDADp4WF4NuMlJb23u9N4P8ACZS4U0zWbU/pa7Id3TAIJYILswpCweYvCv2l8KTopWL3oTokzqAO6BkXuGuJ6AExSqAIJYIOiL11etW5vt83LY0/DPbJYqRp22GiZfZBjh/+2G2insAQ2BglggroXSRaPQC/3gH1nzxP4LS/rhyzfpz5GSnq3OpJhXEd4CEoGCWCC7MKQsHmLwr9pfCk6KVi96E6JM6gDugZF7hriegBMUqgEBgaMAWCJQDZThdHMu+arnPzlatEUHv6mD1lAjwRqVHwwy5ITEKQEBAYIEoVgcSs8nc5F8e1R8V2p/8RDSulczwfHKnNxlmuo6VqFAAQKCAFgg7hVazpxAKSB0y2r/jJzN0nPIFkj/EUnvNrzqbruKPiUCGQNjAxll8ASChAVYHErPJ3ORfHtUfFdqf/EQ0rpXM8HxypzcZZrqOlZYHOCnFDGYEsP3c7oE7F1rP/zVqthQBoBbBHsIJUFYIO4VWs6cQCkgdMtq/4yczdJzyBZI/xFJ7za86m67ij4lhAVYHKZGR0uPVDEmFQa2wnPTB8dWmk62yWtC3UopUgpYHA2U4XRzLvmq5z85WrRFB7+pg9ZQI8EalR8MMuRYIO4VWs6cQCkgdMtq/4yczdJzyBZI/xFJ7za86m67ij4lBaJYHfC1rmY6rqjlABV730uq/W9boM5XWffNQQH8Ey9UGQG/WB3xta5mOq6o5QAVe99Lqv1vW6DOV1n3zUEB/BMvVBkCaQaCoBlfJQgZMJoOgVgcta5mOq6o5QAVe99Lqv1vW6DOV1n3zUEB/BMvVAmhWBy1rmY6rqjlABV730uq/W9boM5XWffNQQH8Ey9UoUEABAtYIOiL11etW5vt83LY0/DPbJYqRp22GiZfZBjh/+2G2inspQCBglgg8tL+s5AgvhiSZjvUKoDKJ5LkZ8Q4wHc/0efawNwVQn5YQC/WUuoXvvxgVKvw6dDgewrrxTdwcY+gVYq6J0K77N5f6KVf6oZEaVSszy8jxhaIxJ/z1kpJ60AKbvUVKbcUWisCgYRYIGMXbFVWamLov399k9fYkGI/rItGhSUB/pMXCCYZ5+4KWEBC4BFD6tx5HUgbBP5sN4aS7AG35Tv0UwGxOIasBdgj8GcnYsnYV7dk+88Z5/uqhduL3fRkLn2+3vux5eK5EFeNQELYXQGCggKBggQaAAEY44ICg4IChIIBg4IAWBxKzydzkXx7VHxXan/xENK6VzPB8cqc3GWa6jpWggBYHLFrVvXsBkvmrDyrYDXvroazZsw9xKDVcWA9cOWCAFgcduYH2yoxyaLDJ2HSQxoYalUMwyH3nNjWqCspuIICgYIAWBx25gfbKjHJosMnYdJDGhhqVQzDIfec2NaoKym4ggUZYzCDAwGBggBYHOCnFDGYEsP3c7oE7F1rP/zVqthQBoBbBHsIJUGCBBmErIIEGVFsBINDQ+9G2GaCGQJGn6EgIdhmghkDXJ9DzG6Kn0IjsSQAA//YZoIZAeefAESgI+rHQ765nkKPwf+hQhMqRPUXsdCA//+fI58FRIDsEEOfQxXoHEI9M/9EdzjGsUD/2GaCGQJHn9hmghkBWp8DBP/YZoIZAaefAUTRrfFURCLzi51DA9yyI///2QVEn0Ag//8FgYQCBECCGzYiiY3jGDitG0PIy1lsSqLx9fY=",
                            "id": "4d9933098f4d53e009c8c222ba267cdfd3ecbdf9071140d1a58c493968d5dbd7",
                            "body": {
                                "inputs": [
                                    {
                                        "txId": "95c3003a78585e0db8c9496f6deef4de0ff000994b8534cd66d4fe96bb21ddd3",
                                        "index": 0
                                    },
                                    {
                                        "txId": "bb30a42c1e62f0afda5f0a4e8a562f7a13a24cea00ee81917b86b89e801314aa",
                                        "index": 0
                                    },
                                    {
                                        "txId": "e88bd757ad5b9bedf372d8d3f0cf6c962a469db61a265f6418e1ffed86da29ec",
                                        "index": 1
                                    }
                                ],
                                "collaterals": [
                                    {
                                        "txId": "ae85d245a3d00bfde01f59f3c4fe0b4bfae1cb37e9cf91929eadcea4985711de",
                                        "index": 2
                                    }
                                ],
                                "references": [
                                    {
                                        "txId": "bb30a42c1e62f0afda5f0a4e8a562f7a13a24cea00ee81917b86b89e801314aa",
                                        "index": 1
                                    }
                                ],
                                "collateralReturn": null,
                                "totalCollateral": null,
                                "outputs": [
                                    {
                                        "address": "addr_test12qxefct5wvh0n2h88uu44dz9q7l6nq7k2q3uzx54ruxr9eyycs5szqg64jrxw",
                                        "value": {
                                            "coins": 4,
                                            "assets": {
                                                "4acf2773917c7b547c576a7ff110d2ba5733c1f1ca9cdc659aea3a56": 1
                                            }
                                        },
                                        "datumHash": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25",
                                        "datum": null,
                                        "script": null
                                    }
                                ],
                                "certificates": [
                                    {
                                        "genesisDelegation": {
                                            "verificationKeyHash": "4acf2773917c7b547c576a7ff110d2ba5733c1f1ca9cdc659aea3a56",
                                            "delegateKeyHash": "e0a714319812c3f773ba04ec5d6b3ffcd5aad85006805b047b082541",
                                            "vrfVerificationKeyHash": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25"
                                        }
                                    },
                                    {
                                        "genesisDelegation": {
                                            "verificationKeyHash": "a646474b8f5431261506b6c273d307c7569a4eb6c96b42dd4a29520a",
                                            "delegateKeyHash": "0d94e174732ef9aae73f395ab44507bfa983d65023c11a951f0c32e4",
                                            "vrfVerificationKeyHash": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25"
                                        }
                                    }
                                ],
                                "withdrawals": {
                                    "stake_test17z66ue36465w2qq40005h2hadad6pnjht8mu6sgplsfj74qvpedfl": 447,
                                    "stake17x66ue36465w2qq40005h2hadad6pnjht8mu6sgplsfj74qttn0dz": 617
                                },
                                "fee": 867,
                                "validityInterval": {
                                    "invalidBefore": 12442,
                                    "invalidHereafter": 26096
                                },
                                "update": {
                                    "proposal": {},
                                    "epoch": 24357
                                },
                                "mint": {
                                    "coins": 0,
                                    "assets": {
                                        "b5ae663aaea8e500157bdf4baafd6f5ba0ce5759f7cd4101fc132f54.00": 4
                                    }
                                },
                                "network": null,
                                "scriptIntegrityHash": "e88bd757ad5b9bedf372d8d3f0cf6c962a469db61a265f6418e1ffed86da29ec",
                                "requiredExtraSignatures": [
                                    "b5ae663aaea8e500157bdf4baafd6f5ba0ce5759f7cd4101fc132f54"
                                ]
                            },
                            "metadata": null,
                            "inputSource": "inputs"
                        },
                        {
                            "witness": {
                                "signatures": {
                                    "c0cd23b758159f90363767d37b8b6bfc48bc592d2c550eb0ded82f3fac10a141": "7GImyTno0VLk1fxx1Nkkhgk4OMpoM6HpiKKtKS8QCiE9ksL4zjru+7eX+J4PLTqgAhZj1GRTNjhN8PdoQU1zvA==",
                                    "b76e9f28fb73ca59a75b6a23232288aa371841fa8df62020ea5f5e15e65d08f1": "rPJR4X/pvkkum/3gYKgEwRnddW20Cs+ldxKNJDNpKmkewvMvLgb8tbjegdl37KaD01vY/oIgs6LrDbdChkZ+Dg==",
                                    "749d0dadac7e7f97e1f9589cec1c89f3f0865b67ac6f7f2ee6534f8a0e881c81": "Vj31T0BjI8vmyYjP5Em/AQx6GkKVTsxqwNEOd5iT9qMVkq9CUHVhv2j/+8HUVYRrkBmpal8gf+XkrElnL9w8zg=="
                                },
                                "scripts": {
                                    "a45d3dfb8f96e7aa7e6c0ce6b3caf2533b4ba54d475bd005a59409b1": {
                                        "native": "e0a714319812c3f773ba04ec5d6b3ffcd5aad85006805b047b082541"
                                    }
                                },
                                "datums": {},
                                "redeemers": {},
                                "bootstrap": [
                                    {
                                        "key": "440a92ff267209984b3c6afbee9c7d2390c91aa8ca92317e9113d6bba7ad8c4f",
                                        "chainCode": null,
                                        "addressAttributes": "QbU=",
                                        "signature": "XEhDW/qGmUZaV7ot/iPhwo4XxaKiMHhFj5bIq040hNLnHFhhHcfLUATAj/SKw+O65RMU4Txsj9tPI6U6wvy19g=="
                                    }
                                ]
                            },
                            "raw": "hKoAgBKDglggAmi+nb0ERuqiF+HeyPOZJJMF5VHX/BQ33YRSH3SqYhwBglgguzCkLB5i8K/aXwpOilYvehOiTOoA7oGRe4a4noATFKoBglgg7hVazpxAKSB0y2r/jJzN0nPIFkj/EUnvNrzqbruKPiUCAYOkAFgdcbFrVvXsBkvmrDyrYDXvroazZsw9xKDVcWA9cOUBggKhWByxa1b17AZL5qw8q2A1766Gs2bMPcSg1XFgPXDloUADAoIB2BhCQXED2BhYIoIAggBYHLWuZjquqOUAFXvfS6r9b1ugzldZ981BAfwTL1SkAFgdcQ2U4XRzLvmq5z85WrRFB7+pg9ZQI8EalR8MMuQBAQKCAFgg6IvXV61bm+3zctjT8M9slipGnbYaJl9kGOH/7YbaKewD2BhYIoIAggBYHOCnFDGYEsP3c7oE7F1rP/zVqthQBoBbBHsIJUGkAFgiQOCnFDGYEsP3c7oE7F1rP/zVqthQBoBbBHsIJUGEuFgBAgECAoIB2BhBQAPYGEqCAUdGAQAAIgARERgbAhkBNwMaAAFZtQWhWB3hNUKss6ZNgMKTAiYNYsO4enQq0Uq/hV68ZzMIHhkBRAaCoBnxfgmhWBxKzydzkXx7VHxXan/xENK6VzPB8cqc3GWa6jpWoUICAiIPAaMAg4JYIMDNI7dYFZ+QNjdn03uLa/xIvFktLFUOsN7YLz+sEKFBWEDsYibJOejRUuTV/HHU2SSGCTg4ymgzoemIoq0pLxAKIT2SwvjOOu77t5f4ng8tOqACFmPUZFM2OE3w92hBTXO8glggt26fKPtzylmnW2ojIyKIqjcYQfqN9iAg6l9eFeZdCPFYQKzyUeF/6b5JLpv94GCoBMEZ3XVttArPpXcSjSQzaSppHsLzLy4G/LW43oHZd+ymg9Nb2P6CILOi6w23QoZGfg6CWCB0nQ2trH5/l+H5WJzsHInz8IZbZ6xvfy7mU0+KDogcgVhAVj31T0BjI8vmyYjP5Em/AQx6GkKVTsxqwNEOd5iT9qMVkq9CUHVhv2j/+8HUVYRrkBmpal8gf+XkrElnL9w8zgKBhFggRAqS/yZyCZhLPGr77px9I5DJGqjKkjF+kRPWu6etjE9YQFxIQ1v6hplGWle6Lf4j4cKOF8WiojB4RY+WyKtONITS5xxYYR3Hy1AEwI/0isPjuuUTFOE8bI/bTyOlOsL8tfZAQkG1AYGCAFgc4KcUMZgSw/dzugTsXWs//NWq2FAGgFsEewglQfX2",
                            "id": "3e6230bf0d2ead922e2296386e20a10f23d14796076262f152fce9e15d857bb9",
                            "body": {
                                "inputs": [],
                                "collaterals": [],
                                "references": [
                                    {
                                        "txId": "0268be9dbd0446eaa217e1dec8f399249305e551d7fc1437dd84521f74aa621c",
                                        "index": 1
                                    },
                                    {
                                        "txId": "bb30a42c1e62f0afda5f0a4e8a562f7a13a24cea00ee81917b86b89e801314aa",
                                        "index": 1
                                    },
                                    {
                                        "txId": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25",
                                        "index": 2
                                    }
                                ],
                                "collateralReturn": null,
                                "totalCollateral": 27,
                                "outputs": [
                                    {
                                        "address": "addr1wxckk4h4asryhe4v8j4kqd0046rtxekv8hz2p4t3vq7hpegtxpwnn",
                                        "value": {
                                            "coins": 2,
                                            "assets": {
                                                "b16b56f5ec064be6ac3cab6035efae86b366cc3dc4a0d571603d70e5": 3
                                            }
                                        },
                                        "datumHash": null,
                                        "datum": "4171",
                                        "script": {
                                            "native": "b5ae663aaea8e500157bdf4baafd6f5ba0ce5759f7cd4101fc132f54"
                                        }
                                    },
                                    {
                                        "address": "addr1wyxefct5wvh0n2h88uu44dz9q7l6nq7k2q3uzx54ruxr9eqj6wm5c",
                                        "value": {
                                            "coins": 1,
                                            "assets": {}
                                        },
                                        "datumHash": "e88bd757ad5b9bedf372d8d3f0cf6c962a469db61a265f6418e1ffed86da29ec",
                                        "datum": null,
                                        "script": {
                                            "native": "e0a714319812c3f773ba04ec5d6b3ffcd5aad85006805b047b082541"
                                        }
                                    },
                                    {
                                        "address": "addr_test1grs2w9p3nqfv8amnhgzwchtt8l7dt2kc2qrgqkcy0vyz2svyhpvqzqsdmflg9",
                                        "value": {
                                            "coins": 2,
                                            "assets": {}
                                        },
                                        "datumHash": null,
                                        "datum": "40",
                                        "script": {
                                            "plutus:v1": "46010000220011"
                                        }
                                    }
                                ],
                                "certificates": [],
                                "withdrawals": {
                                    "stake1uy659t9n5excps5nqgnq6ckrhpa8g2k3f2lc2h4uvuess8syll2gq": 324
                                },
                                "fee": 311,
                                "validityInterval": {
                                    "invalidBefore": null,
                                    "invalidHereafter": 88501
                                },
                                "update": {
                                    "proposal": {},
                                    "epoch": 61822
                                },
                                "mint": {
                                    "coins": 0,
                                    "assets": {
                                        "4acf2773917c7b547c576a7ff110d2ba5733c1f1ca9cdc659aea3a56.0202": -3
                                    }
                                },
                                "network": "mainnet",
                                "scriptIntegrityHash": null,
                                "requiredExtraSignatures": []
                            },
                            "metadata": null,
                            "inputSource": "inputs"
                        }
                    ],
                    "header": {
                        "protocolVersion": {
                            "major": 703,
                            "minor": 285
                        },
                        "opCert": {
                            "hotVk": "wEDZgoCDpAbdSCGY+Kj2HbOLDXk46rWoE0TDahZCwfc=",
                            "count": 2,
                            "kesPeriod": 0,
                            "sigma": "qJTV04y0xWBaJ4O9O0krwxWPTSOvlAlRzcbYFGZo9ngukyJxTxa7xssxpHoz2VREh9Y6xw+9lXRe2pNYDQJIBQ=="
                        },
                        "signature": "BmRc5r2b1gIcsUfgZo1fId7HeCVyRIun7tbKd/5ybeDZl9EpW2Wg6cq7BoLJi4cX/3eeAsaTONNo4xMvNv+PBIRTIy00kYa3QM3teXLSOyrut3RgG6Fvr//kO1Txzw2HskKhbgkUqr9Br2Fz2K8FNoxmzUxilikM4rLj+oQxmbjMhSgyNI79gmEi9qg//rxXyRVIpnV2/6BnUv8A9FYg04ZjbUhGG3H7wnpPFmDfZChlGSBMMMFMdwouzSwlrzdrejZUlXE21h4ku9xdjwegXg59oVzlqqbhWsksEeES4kvYpPbk3U83/HrWV4YYcqgqU4qciU0xwFFot3NxvJqUe1adOoXodOEH55Dk0Rl3eIpbXb9R3PsWOXVls9fImEEYGdG1uayHYi9vhZp9IIqC5pbal5jhivHoOPHrMgh+v+iN+mYT6riqnsZvBvdzSzvnr3sT/l5b8+sxMsu+5JQZzjF6TKojrgBXt4uDM1//CDOb3rt0RdigU0R5deabXlR3ffemHs1/qOEGVNyDc9yQgYThUbeBSyJocKL7P/mvGk/UBoT2/cBgpNyKvLuE6lMMQkmwkuvhND5ABdC7B1eHGA==",
                        "vrfInput": {
                            "output": "diIhAko4yqhZ7pF5O6MpyImxuJ5vFBshesI3nDvHvTaZViOqK7+9cFmd4qitXnnEt9mcnq2W7tTruL1B5xnyIg==",
                            "proof": "n0qb7NDz3SIMisn4ZH8ZyhrWRI8oA5THqIlVhmGN+SYOEflrnn0mQYgNSC4AjqQIaJMMNKFVlgkaNO8vEPzxY2iIZO7QPZH4sVKFqH6DKgA="
                        },
                        "blockHeight": 1,
                        "slot": 5,
                        "prevHash": "genesis",
                        "issuerVk": "53c834690d5b1ecc76fe5e9879d73393c28d682a84c832fc1ca77ed1ce0ba4b1",
                        "issuerVrf": "hV68Yig5mNA9qEZFtibjJEp/joLTWTdVLmx3c8BRXtU=",
                        "blockSize": 2,
                        "blockHash": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25"
                    },
                    "headerHash": "ee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25"
                }
            },
            "tip": {
                "slot": 51204,
                "hash": "aad78a13b50a014a24633c7d44fd8f8d18f67bbb3fa9cbcedf834ac899759dcd",
                "blockNo": 3518078
            }
        }
    }`

	var method1 CompatibleResultNextBlock
	err := json.Unmarshal([]byte(dataRequestNextForward), &method1)
	assert.Nil(t, err)

	dataRequestNextBackward := `{
        "RollBackward": {
            "point": {
                "slot": 92267,
                "hash": "6487fa2e6f0e85ef6e887931381057146060bfd2ed7324f7829c369c3628dc16"
            },
            "tip": {
                "slot": 4744,
                "hash": "92cdf578c47085a5992256f0dcf97d0b19f1f1c9de4d5fe30c3ace6191b6e5db",
                "blockNo": 3015040
            }
        }
	}`

	var method2 chainsync.ResultNextBlockPraos
	err = json.Unmarshal([]byte(dataRequestNextBackward), &method2)
	assert.Nil(t, err)

	dataNextBlockForward := `{
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
 	}`

	var method3 chainsync.ResultNextBlockPraos
	err = json.Unmarshal([]byte(dataNextBlockForward), &method3)
	assert.Nil(t, err)

	dataNextBlockBackward := `{
        "direction": "backward",
        "point": {
            "slot": 92267,
            "id": "b5f556f2ff67952ca1237ccb40dbf33f213ce15b769597d12188e9ab8fbd7bdf"
        },
        "tip": {
            "slot": 4744,
            "id": "669ea8a00b7307b08fae2a6eb3a59d42f7d32e80860320208938019789bef05f",
            "height": 3015040
        }
	}`

	var method4 chainsync.ResultNextBlockPraos
	err = json.Unmarshal([]byte(dataNextBlockBackward), &method4)
	assert.Nil(t, err)

	dataFindIntersectFound := `{
        "IntersectionFound": {
            "point": "origin",
            "tip": {
                "slot": 62344,
                "hash": "2208e439244a1d0ef238352e3693098aba9de9dd0154f9056551636c8ed15dc1",
                "blockNo": 2
            }
        }
	}`

	var method5 CompatibleResultFindIntersection
	err = json.Unmarshal([]byte(dataFindIntersectFound), &method5)
	assert.Nil(t, err)

	dataFindIntersectNotFound := `{
        "IntersectionNotFound": {
            "tip": {
                "slot": 36991,
                "hash": "f63498b4ae65be466e4a71878971b9c524458996450b0ff8262cddf3f0d99229",
                "blockNo": 6
            }
        }
	}`

	var method6 CompatibleResultFindIntersection
	err = json.Unmarshal([]byte(dataFindIntersectNotFound), &method6)
	assert.Nil(t, err)

	dataFindIntersectionFound := `{
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

	var method7 chainsync.ResultNextBlockPraos
	err = json.Unmarshal([]byte(dataFindIntersectionFound), &method7)
	assert.Nil(t, err)

	dataFindIntersectionNotFound := `{
		"error": {
			"code": 1000,
			"message": "No intersection found.",
			"data": {
				"tip": {
					"slot": 68040,
					"id": "ce8c0f3211d39ae0db42fb105d076f91132449d0a96e8ff092b01488af3ae12b",
					"height": 13
				}
			}
		}
	}`

	var method8 chainsync.ResultNextBlockPraos
	err = json.Unmarshal([]byte(dataFindIntersectionNotFound), &method8)
	assert.Nil(t, err)

	dataV5Result := `{"IntersectionFound":{"Point":{"hash":"fe896cf51a9a2309fea7c3527b992c3c4b94073b744b6f29bb844cf0c5491e72","slot":33951595},"Tip":{"blockNo":1470914,"hash":"fe896cf51a9a2309fea7c3527b992c3c4b94073b744b6f29bb844cf0c5491e72","slot":33951595}}}`
	var method9 CompatibleResult
	err = json.Unmarshal([]byte(dataV5Result), &method9)
	assert.Nil(t, err)

	dataV5Result = `{"RollForward":{"block":{"babbage":{"body":[{"id":"03a67c5103c2284ddcb09c60fd79c5c2554ca600bc0c7ae4b55268038ba3af35","inputSource":"inputs","body":{"collaterals":[{"txId":"8eb845bfe9cc4f83ef8e91f71d544d4df47ccaa65811100a0afebf8cc931ff1e","index":0}],"fee":607466,"inputs":[{"txId":"8eb845bfe9cc4f83ef8e91f71d544d4df47ccaa65811100a0afebf8cc931ff1e","index":0},{"txId":"b41dadb5ae3299abb70c16505531847c8e146abdda6fc7f1b7ee084fb7f7f945","index":0},{"txId":"fecc767fd92b535912013a13282481ce2e4084d9b6de75f87bd8e342757eb9ea","index":1}],"mint":{"coins":0},"network":"testnet","outputs":[{"address":"addr_test1wp4y7k88m8mekc0hlcd6z0jgwn8fssy3axhapz74xgeu5hsdzcpdu","datum":"9f0000000000ff","value":{"coins":3000000,"assets":{"5241761209d3d8e8da3210a41d47f43e1e011c4172e1f58c5c9b6363":1}},"script":null},{"address":"addr_test1wr8s48slcgnz44uv2yvu6am5y0zvk93lffw80gnl9ns64tsdqedsl","datum":"9f1b0000013d0d8e41a51a35097ffa1b00006a0c800691b71b000000e38f67bffc1b000000014a3000661b003c59e1d78147bd9fc24994019d7c5fd6832c6dc24c0221e7260207000000000000ff1b0000018bf41a40381b0000018bf420b4909f1b00000221a51c68591b00006a0c800691b7ff1a002dc6c0ff","value":{"coins":3000000,"assets":{"f178d70eb9b63b8782337b6f3d3f4ab6472ec77a9d03de91872b7180":1}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a035097ffff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867161}},"script":null},{"address":"addr_test1wperr0p4kc9w7mwlf23037xq5wz3y8hmfntptv6jy56fmvq50espl","datum":"9f9f0000000000ff1a03509809ff","value":{"coins":3000000,"assets":{"2145ad11bb8e369f83927f5aa706d9cfb87af4c39802cb0a8c28d35e":1,"919d4c2c9455016289341b1a14dedf697687af31751170d56a31466e.7455534443":85163867176}},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767006},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767006},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767006},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767006},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767006},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767006},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767006},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null},{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":337767005},"script":null}],"requiredExtraSignatures":["e4aec08ecfcdcce8eef88f5468c19259a2b265e8ed3820cc5fd0ae72"],"scriptIntegrityHash":"a2b135d4e160f64e4137b0994c8b67df7ca70264d8c184d31e5149e872f3f497","update":null,"validityInterval":{"invalidBefore":33951566,"invalidHereafter":33951866},"collateralReturn":{"address":"addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860","value":{"coins":6077413563},"script":null},"totalCollateral":3000000,"references":[{"txId":"24979f72318121aa1895f310f73a18e46331f5e65c529559592772261a09a42f","index":2},{"txId":"66a5925f4aa180bc8110d063f9dd7790531e92d3d95d13b132e7677af5abec90","index":0},{"txId":"75fd83e6f18fc1a544961bcd71cc93dd8229a73277121665f73a9fdd817381d9","index":0}]},"witness":{"redeemers":{"spend:1":{"redeemer":"1a002dc6c0","executionUnits":{"memory":3094905,"steps":1350385904}},"spend:2":{"redeemer":"d87980","executionUnits":{"memory":32194,"steps":12867753}}},"scripts":{},"signatures":{"af2d2d704143318aed707c71edbf9f13ed00d9517c9133dd69bcf43954f65e8e":"rSIocfWkJ8sw4znmWEjYNsOL1deBWCxXGatLaiK4fCvZ2So5j/WsJsdfKsc/QMR7Yr+AfFnlm2O6Q9K1mmqUBg=="}},"metadata":null,"raw":"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"}],"header":{"blockHash":"aaf7ba6358420d89b7eabdc77260a2b5486e7b5f094273f88696b86446994e5a","blockHeight":1470915,"blockSize":3739,"issuerVK":"9a4872f4228075308877ec3e61728505eef47910d8332a715993386064dccd73","issuerVrf":"e/mBmZkGd6BY6aENOltWlR41dFZaCNS9jfLlykyrcUU=","opCert":{"count":4,"hotVk":"BmZfTprSCXrbNMQk3C/8t0Qyz75eA0U3llNRWMbUVr4=","kesPeriod":244,"sigma":"8rMMQzkQnpJOXJW5xSeq7ac/QPenj0543u2YiYZRjQLywYspHsGhrStM49IULOmVvoVH8MpJ6NTk59sMuUHZBA=="},"prevHash":"fe896cf51a9a2309fea7c3527b992c3c4b94073b744b6f29bb844cf0c5491e72","protocolVersion":{"major":8,"minor":0},"signature":"pR8ZWurWUNhdKXN3xHZAqu0Fjv8ndH5Qh8gZCGDBONyBSrDIdvXLjeGRnh42T0CKgKmotUub4N08wRhosqPAAT0BTE/Ko1t7Vf1M8jg+iQrVO9oI6WdZozZuE+ZxOkMKIiJeIoohHfTca9Tyw68KaVpVLY/P1iuX43kgBDhA8TNUPsfh6ZGJHXXyRok7Z4QeCg07bkZCTm77JH8Mv4yOsfv36Wh7apX1Ap+DuuDpf8mfLWIjDXPNrj2haEGXnHq498dteoZVFsj9RwfhUs7anuakJmx/tfVcG3MlAtk2dAcVXJ6wl76Gtd+Wji46Lrmqlr8j57S6pGbyQdFZStqThSh0T3MV/BwmTDbXbNopNB272uzs08BhY0mEOVf0I0ilXR3C2VuKtFErd50zhGjstvLC7thYq/zR2H7w54UT3WJQypE2H2ggZLf3kxrBpp6YLSued92c5mKkluYTwDcY/31IYbmUobTGVIEP7cuG4ANrVlPvOwKOtZJPEP84sKM+yQBhaqDK8tkcaVLXlJaApqyqnQzfBTmn5v+6NPYS8QNmq0sDPyxYEfZEwt6Mnni0x/Ol3YC8m+mtce0k6svRjQ==","slot":33951638},"headerHash":"91f1c4ab70f5c2ea40672db921666e54162a5358e0a983c18a04b76bfcae85f4"}},"tip":{"blockNo":1470915,"hash":"91f1c4ab70f5c2ea40672db921666e54162a5358e0a983c18a04b76bfcae85f4","slot":33951638}}}`
	var method10 CompatibleResult
	err = json.Unmarshal([]byte(dataV5Result), &method10)
	assert.Nil(t, err)
}

func TestTasteTestBlock(t *testing.T) {
	example := `{
        "RollForward": {
            "block": {
                "babbage": {
                    "body": [
                        {
                            "id": "d55e007d49f88acc77695162ee641fc7e18e1510343781f0ec49d64e1838dac6",
                            "inputSource": "inputs",
                            "body": {
                                "fee": 181341,
                                "inputs": [
                                    {
                                        "txId": "4bc0d7d3330c1fb371d3971dcf4dabfecd653d521970fddc3502cbbe2c44d066",
                                        "index": 18
                                    },
                                    {
                                        "txId": "7f1e5193ac64b5ca50d04570e9f2a9be15d91901734946a363355ec7eb3e6df5",
                                        "index": 7
                                    },
                                    {
                                        "txId": "7f1e5193ac64b5ca50d04570e9f2a9be15d91901734946a363355ec7eb3e6df5",
                                        "index": 8
                                    },
                                    {
                                        "txId": "cac81008b25eec4652807529ff7fd051dbb9871d3ea57facefd234478e1f74b7",
                                        "index": 18
                                    },
                                    {
                                        "txId": "cac81008b25eec4652807529ff7fd051dbb9871d3ea57facefd234478e1f74b7",
                                        "index": 23
                                    },
                                    {
                                        "txId": "cac81008b25eec4652807529ff7fd051dbb9871d3ea57facefd234478e1f74b7",
                                        "index": 25
                                    },
                                    {
                                        "txId": "cac81008b25eec4652807529ff7fd051dbb9871d3ea57facefd234478e1f74b7",
                                        "index": 27
                                    },
                                    {
                                        "txId": "cac81008b25eec4652807529ff7fd051dbb9871d3ea57facefd234478e1f74b7",
                                        "index": 31
                                    }
                                ],
                                "mint": {
                                    "coins": 0
                                },
                                "network": "testnet",
                                "outputs": [
                                    {
                                        "address": "addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860",
                                        "value": {
                                            "coins": 193807538
                                        },
                                        "script": null
                                    },
                                    {
                                        "address": "addr_test1vrj2asywelxue68wlz84g6xpjfv69vn9arknsgxvtlg2uusqey860",
                                        "value": {
                                            "coins": 1925867757
                                        },
                                        "script": null
                                    }
                                ],
                                "update": null,
                                "validityInterval": {}
                            },
                            "witness": {
                                "redeemers": {},
                                "scripts": {},
                                "signatures": {
                                    "af2d2d704143318aed707c71edbf9f13ed00d9517c9133dd69bcf43954f65e8e": "DZZk+OYM4VyUe3e+CNsoUDjeobIqjJtzjzHga5gNx8HCCOGgimLgxxTgKoz1E5Fjx6GOFWQrRXZOBU+/hC0aCw=="
                                }
                            },
                            "metadata": null,
                            "raw": "hKQAiIJYIEvA19MzDB+zcdOXHc9Nq/7NZT1SGXD93DUCy74sRNBmEoJYIH8eUZOsZLXKUNBFcOnyqb4V2RkBc0lGo2M1XsfrPm31B4JYIH8eUZOsZLXKUNBFcOnyqb4V2RkBc0lGo2M1XsfrPm31CIJYIMrIEAiyXuxGUoB1Kf9/0FHbuYcdPqV/rO/SNEeOH3S3EoJYIMrIEAiyXuxGUoB1Kf9/0FHbuYcdPqV/rO/SNEeOH3S3F4JYIMrIEAiyXuxGUoB1Kf9/0FHbuYcdPqV/rO/SNEeOH3S3GBmCWCDKyBAIsl7sRlKAdSn/f9BR27mHHT6lf6zv0jRHjh90txgbglggysgQCLJe7EZSgHUp/3/QUdu5hx0+pX+s79I0R44fdLcYHwGCglgdYOSuwI7Pzczo7viPVGjBklmismXo7TggzF/QrnIaC41EsoJYHWDkrsCOz83M6O74j1RowZJZorJl6O04IMxf0K5yGnLKaO0CGgACxF0PAKMAgYJYIK8tLXBBQzGK7XB8ce2/nxPtANlRfJEz3Wm89DlU9l6OWEANlmT45gzhXJR7d74I2yhQON6hsiqMm3OPMeBrmA3HwcII4aCKYuDHFOAqjPUTkWPHoY4VZCtFdk4FT7+ELRoLBIAFgPX2"
                        },
                        {
                            "id": "3b8fa78469b7dc9f481b7773d871254b58fb60fe395c1dae6b22b0eed64176ce",
                            "inputSource": "inputs",
                            "body": {
                                "collaterals": [
                                    {
                                        "txId": "09f227731bfae6b423ba2f760286dea1120e071528f71137dd92c37b37ac22c4",
                                        "index": 1
                                    }
                                ],
                                "fee": 652685,
                                "inputs": [
                                    {
                                        "txId": "09f227731bfae6b423ba2f760286dea1120e071528f71137dd92c37b37ac22c4",
                                        "index": 1
                                    },
                                    {
                                        "txId": "4be5dffcd6636e0869c3aa05a96314443f2c75717894ee0312fe81cfe22224eb",
                                        "index": 0
                                    },
                                    {
                                        "txId": "a52092e8f6444f2105af2e69fc2c7b190ecf200fc3ec29bc7f509c3b572bb56c",
                                        "index": 0
                                    }
                                ],
                                "mint": {
                                    "coins": 0,
                                    "assets": {
                                        "541529653c02d7155a7c0d68f1854e132d129538bb85256eb79db127.46534e121fd22e0b57ac206fefc763f8bfa0771919f5218b40691eea4514d0": -1
                                    }
                                },
                                "network": null,
                                "outputs": [
                                    {
                                        "address": "addr_test1wql9ncpnj7pffkjcgpwr3t6ws2g0hh5p7368u6lpgura0ssh79m9w",
                                        "datum": "d8799fd8799f581c0d33957c07acdddecc9882457da22f05e0d189f7fc95b1972e6d5105ffd8799f581c26977346f8c25a12f6101e0a06385abad18d65530d420203a8560b71ffff",
                                        "value": {
                                            "coins": 604000000,
                                            "assets": {
                                                "541529653c02d7155a7c0d68f1854e132d129538bb85256eb79db127.46534e0d33957c07acdddecc9882457da22f05e0d189f7fc95b1972e6d5105": 1
                                            }
                                        },
                                        "script": null
                                    },
                                    {
                                        "address": "addr_test1qrp8nglm8d8x9w783c5g0qa4spzaft5z5xyx0kp495p8wksjrlfzuz6h4ssxlm78v0utlgrhryvl2gvtgp53a6j9zngqtjfk6s",
                                        "value": {
                                            "coins": 18030776540
                                        },
                                        "script": null
                                    }
                                ],
                                "requiredExtraSignatures": [
                                    "121fd22e0b57ac206fefc763f8bfa0771919f5218b40691eea4514d0"
                                ],
                                "scriptIntegrityHash": "34e728f19feb9e9ffae260a392847a27f963cee480de00545dc8f8e1fdda5e60",
                                "update": null,
                                "validityInterval": {
                                    "invalidBefore": 34011737,
                                    "invalidHereafter": 34012937
                                },
                                "collateralReturn": {
                                    "address": "addr_test1qrp8nglm8d8x9w783c5g0qa4spzaft5z5xyx0kp495p8wksjrlfzuz6h4ssxlm78v0utlgrhryvl2gvtgp53a6j9zngqtjfk6s",
                                    "value": {
                                        "coins": 16043005753
                                    },
                                    "script": null
                                },
                                "totalCollateral": 979028
                            },
                            "witness": {
                                "redeemers": {
                                    "spend:1": {
                                        "redeemer": "d87980",
                                        "executionUnits": {
                                            "memory": 98848,
                                            "steps": 36557395
                                        }
                                    },
                                    "spend:2": {
                                        "redeemer": "d87980",
                                        "executionUnits": {
                                            "memory": 107634,
                                            "steps": 40658473
                                        }
                                    },
                                    "mint:0": {
                                        "redeemer": "d87c9f581c121fd22e0b57ac206fefc763f8bfa0771919f5218b40691eea4514d0d8799fd8799f581c0d33957c07acdddecc9882457da22f05e0d189f7fc95b1972e6d5105ffd8799f581c26977346f8c25a12f6101e0a06385abad18d65530d420203a8560b71ffffff",
                                        "executionUnits": {
                                            "memory": 518776,
                                            "steps": 193907336
                                        }
                                    }
                                },
                                "scripts": {
                                    "3e59e033978294da58405c38af4e8290fbde81f4747e6be14707d7c2": {
                                        "plutus:v2": "5909d1010000322223232325333573466e1d20000021323232323232533357346644664646460044660040040024600446600400400244a666aae7c0045288a9991199ab9a00200114a060066ae840044c008d5d1000919919180111980100100091801119801001000912999aab9f00114a02a666ae68cdd79aba100100314a2260046ae88004cc8c8c8c0088cc0080080048c0088cc008008004894ccd55cf8008a5eb804cd5d018019aba1001300235744002aae7400c004dd59aba100333223233232323002233002002001230022330020020012253335573e002297ae0133574060066ae84004c008d5d1000aab9d33232323002233002002001230022330020020012253335573e002297adef6c60132533357346008002266ae80004c00cd5d1001098019aba2002357420024664646460044660040040024600446600400400244a666aae7c0045280a99919ab9a00114a260066ae840044c008d5d1000918019bae35573a0026eacd55cf000801919b8f33371890001b8d00200100237566ae84d5d1000a4410346534e00149858c8d55cf1baa0013232325333573466e1d200200213322332323002233002002001230022330020020012253335573e0022c264a666ae68cc88cdd79ba73235573c6ea8004008dd3991aab9e37540020020086ae840044d5d09aba20011300335744004646aae78dd50009aba1001002004357420022c6aae78008d55ce8009baa357426ae88014dd61aba1003357446ae88004d5d11aba20013235573c6ea8004d5d0800991aab9e37540020082a666ae68cdc3a4008004266664664446646460044660040040024600446600400400244a666aae7c004489400454ccd5cd19baf35573a6ae840040104c014d5d0800898011aba2001001232223002003375a6aae78004004d5d0991aba235744002646aae78dd5000803991bab35742646ae88d5d11aba235744646ae88d5d1000800991aab9e37540020026ae84c8d55cf1baa0010042498584c8c8c8c8c8c8c8c8c94ccd5cd19b87480080084c8c8c8c94ccd5cd19b874801000854ccd5cd199119baf374e646aae78dd50008011ba73235573c6ea8004004c8c8c8c8c80154ccd5cd19b87480000084c8c8c8c8c8c8c8c8c8c8c8c8c92653335573e0022930b1aba20065333573466e1d20000021323232324994ccd55cf8008a4c2c6ae8800cdd70009aba100115333573466e1d20020021324994ccd55cf8008a4c2c2c6aae78008d55ce8009baa001357420026ae880194ccd5cd19b87480000084c8c8c8c92653335573e0022930b1aba2003375c0026ae8400454ccd5cd19b87480080084c92653335573e0022930b0b1aab9e00235573a0026ea8004d5d08008b1aab9e00235573a0026ea8004d5d08008098a999ab9a33223375e6e98008dd3000992999aab9f00110011333573466ebcd55ce9aba10013752910100357440020026eacd5d09aba20083253335573e00220022666ae68cdd79aab9d357420026ea522100357440020026eacd5d08020a999ab9a3371064a666aae7c004520001333573466ebcd55ce9aba1001375291100375a6aae78d5d09bab35573c6ae84005200037566ae84d5d100419b803253335573e002290000999ab9a3375e6aae74d5d08009ba9488100375a6aae78d5d09bab35573c6ae84005200037566ae840112080dac4091533357346644664646460044660040040024600446600400400244a666aae7c0045288a9991199ab9a00200114a060066ae840044c008d5d1000918019bab35573c002002644664646460044660040040024600446600400400244a666aae7c0045288a9991199ab9a00200114a060066ae840044c008d5d1000918019bad35573c002002466e1c0052000332233323222333323232300223300200200123002233002002001322323223300100300222253335573e002266ae8000c0084c8c8c94ccd5cd19baf002001133574066ec0008cc024d55cf0031aab9e00333300822002005357440082a666ae68cdc81bae002375c002266ae80018cccc0208800400cd5d1002002899aba0003333300822001006005357440086aae74008d55ce8021aba10012253335573e004200226666006440026ae84008d5d100100080080191001001000911ba63300337560046eac00488ccc888cccc8c8c8c0088cc0080080048c0088cc008008004c88c8c88cc00400c0088894ccd55cf800899aba000300213232325333573466ebc0080044cd5d019bb00023300935573c00c6aae7800cccc02088008014d5d10020a999ab9a337206eb8008dd7000899aba00063333008220010033574400800a266ae8000ccccc02088004018014d5d10021aab9d00235573a0086ae84004894ccd55cf801080089999801910009aba10023574400400200200644004004002446ea0cdc09bad002375a0020040020040026eacd5d080525eb7bdb18054ccd5cd199119192999ab9a323370e6aae74dd5000a40046ae84d5d100109919192999ab9a3370e90000010a5015333573466e1d2002002132337100020106eb4d5d08008a5135573c0046aae74004dd500089919192999ab9a3370e90000010a99919ab9a00114a29405280a999ab9a3370e90010010a99919ab9a00114a26466e20004020dd69aba10011323370e0100026eb4d5d08008a99919ab9a00114a29445281aab9e00235573a0026ea8004d5d0800991aab9e37540026ae84d5d1191aab9e375400200266e04dd69aba13235573c6ea80040512080c60a35742646ae88d5d10009aba200a1498585858585858d55cf0011aab9d00137546ae84d5d10009aba23235573c6ea8004cc88cc8c8c0088cc0080080048c0088cc008008004894ccd55cf8008b09919192999ab9a3370e90010010a999ab9a3371e00c6eb8d5d080089aba100413005357440082600a6ae88010d55cf0011aab9d0013754646ae84c8d55cf1baa0010013235742646aae78dd50008009aba100100237586ae8401cdd71aba10011635573c0046aae74004dd5191aba13235573c6ea8004004d5d0800991aab9e3754002646464a666ae68cdc3a4004004266446646460044660040040024600446600400400244a666aae7c004584c94ccd5cd199119baf374e646aae78dd50008011ba73235573c6ea8004004010d5d080089aba135744002260066ae88008c8d55cf1baa001357420020040086ae8400458d55cf0011aab9d00137546ae84d5d10029bac357420066ae88d5d10009aba235744002646aae78dd50009aba10013235573c6ea8004010d55cf0011aab9d001375400498183d8799f1b0000018c54d934e0d8799fd8799f581cfda1f7894b02beed50bd45a5839e85e1863d83869cd4d7767b6dfcbeffd8799fd8799fd8799f581cb4a78f7f93ef6bf8d862cb650af20857d36221705d53c4b7eca3f8d4ffffffffd8799fd87a9f581c8b5a501a438741ea4c75bc14e54662970977afd02781c3bc4911e70fffffff0001"
                                    },
                                    "541529653c02d7155a7c0d68f1854e132d129538bb85256eb79db127": {
                                        "plutus:v2": "5918ef0100003222323232323232323232323232323253335734a666aae7c0085288991991199ab9a0020014a0664646460044660040040024600446600400400244a666aae7c0045288a9991199ab9a00200114a060066ae840044c008d5d10008009aba200333232323002233002002001230022330020020012253335573e002294454cc88ccd5cd0010008a50300335742002260046ae88004004008c88cc88cdd79ba73235573c6ea8004008dd3991aab9e37540020020046ae84c8d55cf1baa0010013235742646aae78dd50008009aba1002132323232325333573466e1d20000021533357346644664646460044660040040024600446600400400244a666aae7c0045280a99919ab9a00114a260066ae840044c008d5d10009199119baf374e646aae78dd50008011ba73235573c6ea800400400cd5d0991aab9e37540020020026ae84c8d55cf1baa00101637586ae8404854ccd5cd19198009125014a200a266644666464600446600400400246004466004004002446600244a666ae68c01800848c88cc004014008cc88cc8c8c0088cc0080080048c0088cc00800800488cc00488cc8888cc00801000c008c010004400c00800400848cc014008cc8888cc00801000c00c004580048940048c00488c94cc88ccd5cd0010008a50323232325333573466e1d200000214a02944d55cf0011aab9d00137540026ae840044c8c8c8c94ccd5cd19b87480000085280a5135573c0046aae74004dd50009aba135744002646aae78dd5000800802112999ab9a3233001224a0294400454ccd5cd19999111199991991119919180111980100100091801119801001000912999aab9f001122500115333573466ebcd55ce9aba10010041300535742002260046ae880040048c888c00800cdd59aab9e00137520020080024664466ebcdd30011ba60013322332233574066ec00080052f5bded8c06ea4008dd40008020018008b1bae00f48810346534e00480080385261616161615333573466e1d200200213332233323230022330020020012300223300200200122330012253335734600c004246446600200a0046644664646004466004004002460044660040040024466002446644446600400800600460080022006004002004246600a004664444660040080060060022c00244a002460024464a6644666ae68008004528191919192999ab9a3370e90000010a5014a26aae78008d55ce8009baa00135742002264646464a666ae68cdc3a400000429405289aab9e00235573a0026ea8004d5d09aba20013235573c6ea8004004014894ccd5cd19198009125014a20022a666ae68c8cc00489280a510061533357346666444466664664446646460044660040040024600446600400400244a666aae7c004489400454ccd5cd19baf35573a6ae840040104c014d5d0800898011aba200100123222300200337566aae78004dd48008020009199119baf374c0046e98004cc88cc88cd5d019bb00020014bd6f7b6301ba900237500020080060022c6eb803d22010346534e0032337029000000a400401c2930b0b0b0a999ab9a3370e90020010a999ab9a33223335734004002940cc88cc88c94ccd5cd1919b8735573a6ea80052002357426ae880084c8c8c94ccd5cd19b87480000085280a999ab9a3370e900100109919b88001007375a6ae840045289aab9e00235573a0026ea80044c8c8c94ccd5cd19b874800000854cc8cd5cd0008a514a0294054ccd5cd19b874800800854cc8cd5cd0008a51323371000200e6eb4d5d080089919b87007001375a6ae8400454cc8cd5cd0008a514a22940d55cf0011aab9d00137540026ae84004008c8d55cf1baa001357426ae88c8d55cf1baa001001375a6ae84d5d1191aab9e375400202c6ae84038cc8c8c8c0088cc0080080048c0088cc008008004894ccd55cf8008a5015333573466ebcd5d08008018a5113002357440026ae84004dd61aba13574401c2664464a666ae68cc88c94cc88ccd5cd0010008a503232325333573466e1d200200214a2266e40dd71aba100100535573c0046aae74004dd51aba100113232325333573466e1d200000213372000a6eb8d5d08008a5135573c0046aae74004dd51aba135744002646aae78dd50008010010008999911999191801119801001000918011198010010009119800912999ab9a3006002123223300100500233223323230022330020020012300223300200200122330012233222233002004003002300400110030020010021233005002332222330020040030030011600122500123001223375e00a00201044a666ae68c8cc00489280a510011300222323253335734a6644666ae680080045281991919180111980100100091801119801001000911980091299919ab9a00114a2600a0042600800229408c004894cc88ccd5cd0010008a5033223375e6e98008dd3000801003899baf00500100d13001332233223374a900019aba000233574000297ae03374a900019aba0375200497ae0357426ae88c8d55cf1baa0010010070081533357346666444466664664446646460044660040040024600446600400400244a666aae7c004489400454ccd5cd19baf35573a6ae840040104c014d5d0800898011aba200100123222300200337566aae78004dd48008020009199119baf374c0046e98004cc88cc88cd5d019bb00020014bd6f7b6301ba900237500020080060022c6eb8058c8cdc524410346534e000010074800805452616162332323230022330020020012300223300200200122330012253323357340022944c0140084c010004528118009119baf00400100d332233223374a900019aba000233574000297ae035742646aae78dd500080119ba548000cd5d01ba90014bd700030028b0b1bae002357420026ae84d5d10008b0992999ab9a323232533357346466e1cd55ce9baa00148008d5d09aba200213232325333573466e1d200000214a02a666ae68cdc3a400800429444cdc40039bad357420026aae78008d55ce8009baa00113232325333573466e1d2000002153323357340022945280a5015333573466e1d2002002153323357340022944cdc40039bad3574200226466e1c020004dd69aba1001153323357340022945288a5035573c0046aae74004dd50009aba10013235573c6ea8004d5d0991aab9e37540020026ae8403c4c8c8ccc02088cc00489840085888c94ccd5cd19b8733223333232222332323002233002002001230022330020020012253335573e002200a2a666ae68cdd79aab9d3574200200c260086aae78d5d0800898011aba2001001375200200490001199991911119919180111980100100091801119801001000912999aab9f001100515333573466ebcd55ce9aba10010061300435573c6ae840044c008d5d10008009ba900100248001d69bab001005375c02600290010a999ab9a333322223333233222332323002233002002001230022330020020012253335573e002244a0022a666ae68cdd79aab9d357420020082600a6ae840044c008d5d10008009191118010019bab35573c0026ea40040100048cc88cdd79ba6002374c0026644664466ae80cdd8001000a5eb7bdb180dd48011ba800100400300116375c0260026466e0520000014800804854ccd5cd1991919180111980100100091801119801001000912999aab9f00114a02a666ae68cdd79aba100100314a2260046ae88004014dd61aba1357440282a666ae68cdc4a400c6466646464600446600400400246004466004004002444a666aae7c00440084cc00ccc010008d5d08009aba2001223370000464a666aae7c004520001332332322253335573e002200426600666e0000920023574400246600400400246444a666aae7c00440084cc00ccdc0001240046ae880048cc008008004004cdc02400090011aba200137566aae780052000001003149858585858c8cdc52450346534e00001003375c0026ae8400854ccd5cd19911991192999ab9a323370e6aae74dd5000a40046ae84d5d100109919192999ab9a3370e90000010a5015333573466e1d20020021323371000200e6eb4d5d08008a5135573c0046aae74004dd500089919192999ab9a3370e90000010a99919ab9a00114a29405280a999ab9a3370e90010010a99919ab9a00114a26466e2000401cdd69aba10011323370e00e0026eb4d5d08008a99919ab9a00114a29445281aab9e00235573a0026ea8004d5d0800801191aab9e37540026ae84d5d1191aab9e37540020020026ae8403c4cc88c94ccd5cd19911929991199ab9a00200114a0646464a666ae68cdc3a400400429444cdc81bae3574200200a6aae78008d55ce8009baa357420022646464a666ae68cdc3a4000004266e40014dd71aba100114a26aae78008d55ce8009baa357426ae88004c8d55cf1baa001002002001132333001332233223374a900019aba000233574000297ae035742646aae78dd500080119ba548000cd5d01ba90014bd70001801005111998019991199119ba548000cd5d000119aba00014bd7019ba548000cd5d01ba90024bd701aba135744646aae78dd5000800802002800912999ab9a3233001224a029440044c8c01488c014894ccd5cd1991919180111980100100091801119801001000911980091299919ab9a00114a2600a0042600800229408c004894cc88ccd5cd0010008a503375e01e0022664466ebcdd30011ba60010070020121533357346666444466664664446646460044660040040024600446600400400244a666aae7c004489400454ccd5cd19baf35573a6ae840040104c014d5d0800898011aba200100123222300200337566aae78004dd48008020009199119baf374c0046e98004cc88cc88cd5d019bb00020014bd6f7b6301ba900237500020080060022c6eb806c014c8cdc0a4000002900100d0a999ab9a33232323002233002002001230022330020020012253335573e002294054ccd5cd19baf3574200200629444c008d5d10008069bac357426ae880704c94ccd5cd2999ab9a3322332232533357346466e1cd55ce9baa00148008d5d09aba200213232325333573466e1d200000214a02a666ae68cdc3a400400426466e2000401cdd69aba100114a26aae78008d55ce8009baa00113232325333573466e1d2000002153323357340022945280a5015333573466e1d2002002153323357340022944c8cdc40008039bad3574200226466e1c01c004dd69aba1001153323357340022945288a5035573c0046aae74004dd50009aba10010023235573c6ea8004d5d09aba23235573c6ea8004004cdc09bad3574200290407859291aba101d14a22664646460044660040040024600446600400400244a666aae7c0045280a99919ab9a00114a260066ae840044c008d5d1000929991199ab9a00200114a066ebcd5d0991aab9e37540020026ae84d5d1001099b8933223370066e0c0080054ccd5cd19b8848000cdc30010008a400429000192999aab9f0011480004ccd5cd19baf35573a6ae84004dd4a4500375a6aae78d5d09bab35573c6ae84005200000448020c94ccd55cf8008a40002666ae68cdd79aab9d357420026ea522100375a6aae78d5d09bab35573c6ae84005200037566ae84d5d1191aab9e37540020026eb0d5d080f8a4c2c6ae88c8d55cf1baa00102416161632337149110346534e000010061622332233323230022330020020012300223300200200122330012253335734600c004246446600200a0046644664646004466004004002460044660040040024466002446644446600400800600460080022006004002004246600a004664444660040080060060022c00244a002460024466ebc01400400458dd70011aba1002357426ae8800858dd69aba135744646aae78dd500080b1aab9e00235573a0026ea8048cc8c8c8c0088cc0080080048c0088cc00800800488cc00488cc8888cc00801000cc014008c01000448940048c8c8c8ccc8c8c0088cc0080080048c0088cc008008004894ccd55cf8008b0a999ab9a3375e6aae74d5d0800807899191198008018011bad35573c64a666aae7c004584c94ccd55cf80089aba100216357440026eacd55cf1aba1002375c6aae74c94ccd55cf8008b0992999aab9f0011357420042c6ae88004dd59aab9e35742002260046ae8800400488c8c8c94ccd5cd19b87480100084c94ccd5cd19baf3253335573e0022c264a666aae7c0044d5d08010b1aba20013233232323002233002002001230022330020020012253335573e002297ae0133574060066ae84004c008d5d1000aab9d33232323002233002002001230022330020020012253335573e002297adef6c60132533357346008002266ae80004c00cd5d1001098019aba2002357420024664646460044660040040024600446600400400244a666aae7c0045280a99919ab9a00114a260066ae840044c008d5d1000918019bae35573a0026eacd55cf000804119b8f33371890001b8d4890346534e000014890346534e0001315333573466e1cc8ccc8c8c8c0088cc0080080048c0088cc0080080048894ccd55cf8008801099801998020011aba1001357440024466e00008c94ccd55cf8008a400026646646444a666aae7c00440084cc00ccdc0001240046ae880048cc0080080048c8894ccd55cf800880109980199b8000248008d5d100091980100100080099b80480012002357440026eacd55cf000a400000200e90020a999ab9a3370e00a90010a999ab9a323232325333573466e1d200000213232325333573466e1d20000021337206eb8d5d08021bae357420022944d55cf0011aab9d00137546ae84d5d10020a5135573c0046aae74004dd51aba10013235573c6ea800400454ccd5cd199119801119801119b8f00200114a02660024940528a999ab9a3371e666e312000371a9110346534e000064890346534e0015333573466e212006371a00c264446004006666e312006337026e3401920060061225001163232325333573466e1d2000002132223002003375c6ae840044894004d55cf0011aab9d00137546ae84c8d55cf1baa00100115333573466e252080a4e8033253335573e002290000999ab9a3375e6aae74d5d08009ba9488100375a6aae78d5d09bab35573c6ae840052000007123300100800216161616161632323232320055333573466e1d20000021323232323232323232323232324994ccd55cf8008a4c2c6ae880194ccd5cd19b87480000084c8c8c8c92653335573e0022930b1aba2003375c0026ae8400454ccd5cd19b87480080084c92653335573e0022930b0b1aab9e00235573a0026ea8004d5d08009aba20065333573466e1d20000021323232324994ccd55cf8008a4c2c6ae8800cdd70009aba100115333573466e1d20020021324994ccd55cf8008a4c2c2c6aae78008d55ce8009baa001357420022c6aae78008d55ce8009baa001357420022c6aae78008d55ce8009baa357426ae88010dd59aba100135744002646aae78dd5000800992999aab9f00112250011332222330020040033574200264664644a666aae7c00448940044cc8888cc00801000cd5d080098011aba2001233002002001232253335573e002244a0022664444660040080066ae84004c008d5d10009198010010008009aba2001002332323230022330020020012300223300200200122330012233222233002004003300500230040011225001232323223300100200337566ae84008c8c8c8c94ccd5cd19b87480100084d5d08008b1aab9e00235573a0026ea8004d5d09aba200135744646aae78dd5000800992999aab9f00112250011332222330020040033574200264664644a666aae7c00448940044cc8888cc00801000cd5d080098011aba2001233002002001232253335573e002244a0022664444660040080066ae84004c008d5d10009198010010008009aba200100216300237586ae84028c004008c8c8c8c0088cc0080080048c0088cc008008004894ccd55cf8008a5eb804c94ccd5cd1802000899aba0001300335744004260066ae88008d5d0800918019bab357426ae88c8d55cf1baa001001300237586ae840248cc8c8c8c0088cc0080080048c0088cc008008004894ccd55cf8008a50153323357340022944c00cd5d0800898011aba200123375e6aae74004014004c8c8c8c0088cc0080080048c0088cc008008004894ccd55cf8008a5eb804cd5d018019aba100130023574400246ae84d5d1191aab9e3754002002646646446646460044660040040024600446600400400244a666aae7c00452f5bded8c026466600a6aae78d5d0801119aba0337606aae74d5d0801800801080098011aba2001001233300237560024644460040066e9800448940048c94ccd5cd1aba3001122500112230020033322332323002233002002001230022330020020012253335573e002297adef6c6013233300535573c6ae840088cd5d019bb035573a6ae8400c0040084004c008d5d100080092999ab9a3375e0026ea120001225001122300200300100137566ae8400cc8c8c94ccd5cd19b87480000084d5d08008b1aab9e00235573a0026ea8d5d09aba200632357446ae88004d5d10009aba2357440026ae88004d5d1000991aab9e37540026ae84004c8d55cf1baa0010014c0188d8799fd8799fd8799f5820341f7811d34af5b7e5c7abb1bf17f772da4ad5094f06b4a2bc015f283f5e6fbeff00ff1b0000018c54d934e0d8799fd8799f581cfda1f7894b02beed50bd45a5839e85e1863d83869cd4d7767b6dfcbeffd8799fd8799fd8799f581cb4a78f7f93ef6bf8d862cb650af20857d36221705d53c4b7eca3f8d4ffffffffff0001"
                                    }
                                },
                                "signatures": {
                                    "52b92d51dc638d085f8663103d5509f0da29bbee418d75f1f2dc7025d69c9643": "LjurZ4qMyrH6SDKn5PdZyO668Oeh3hseQ3MAQf0pmQeWjtzSQs7ThOswE3NS1eqOqdYOLaUkIlLy5/UJ+EEMCw==",
                                    "ada7452fdc47bae69310c44022a0624b2b835c42a92d9eb0353adb5b363ad2d8": "9QSgzE/H7LeJxkv/jneN8r68CYDtAbPbYgG80STsu4jSKArT1cjN8zqNsTHR+UgI/o+qYLrKK+c/iZJGjyIRDw=="
                                }
                            },
                            "metadata": null,
                            "raw": "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VMzVzRm4dIAAAITIyMjIyMjIyMjIyMjJJlMzVXPgAikwsauiAGUzNXNGbh0gAAAhMjIyMkmUzNVc+ACKTCxq6IAM3XAAmroQARUzNXNGbh0gAgAhMkmUzNVc+ACKTCwsaq54AI1VzoAJuqABNXQgAmrogBlMzVzRm4dIAAAITIyMjJJlMzVXPgAikwsauiADN1wAJq6EAEVMzVzRm4dIAIAITJJlMzVXPgAikwsLGqueACNVc6ACbqgATV0IAIsaq54AI1VzoAJuqABNXQgAgJipmauaMyIzdebpgAjdMACZKZmqufABEAETM1c0ZuvNVc6auhABN1KRAQA1dEACACbqzV0Jq6IAgyUzNVc+ACIAImZq5ozdeaq501dCACbqUiEANXRAAgAm6s1dCAIKmZq5ozcQZKZmqufABFIAATM1c0ZuvNVc6auhABN1KREAN1pqrnjV0JurNVc8auhABSAAN1Zq6E1dEAQZuAMlMzVXPgAikAAJmauaM3XmqudNXQgAm6lIgQA3WmqueNXQm6s1Vzxq6EAFIAA3VmroQBEggNrECRUzNXNGZEZkZGRgBEZgBABAAkYARGYAQAQAJEpmaq58AEUoipmRGZq5oAIAEUoGAGauhABEwAjV0QAJGAGbqzVXPAAgAmRGZGRkYARGYAQAQAJGAERmAEAEACRKZmqufABFKIqZkRmauaACABFKBgBmroQARMAI1dEACRgBm601VzwAIAJGbhwAUgADMiMzIyIjMzIyMjACIzACACABIwAiMwAgAgATIjIyIzABADACIiUzNVc+ACJmroAAwAhMjIyUzNXNGbrwAgARM1dAZuwACMwCTVXPADGqueADMzAIIgAgBTV0QAgqZmrmjNyBuuACN1wAImaugAGMzMAgiABADNXRACACiZq6AAMzMwCCIAEAYAU1dEAIaq50AI1VzoAhq6EAEiUzNVc+AEIAImZmAGRAAmroQAjV0QAQAIAIAZEAEAEACRG6YzADN1YARurABIjMyIjMzIyMjACIzACACABIwAiMwAgAgATIjIyIzABADACIiUzNVc+ACJmroAAwAhMjIyUzNXNGbrwAgARM1dAZuwACMwCTVXPADGqueADMzAIIgAgBTV0QAgqZmrmjNyBuuACN1wAImaugAGMzMAgiABADNXRACACiZq6AAMzMwCCIAEAYAU1dEAIaq50AI1VzoAhq6EAEiUzNVc+AEIAImZmAGRAAmroQAjV0QAQAIAIAZEAEAEACRG6gzcCbrQAjdaACAEACAEACbqzV0IBSXre9sYBUzNXNGZEZGSmZq5oyM3DmqudN1QAKQARq6E1dEAEJkZGSmZq5ozcOkAAAEKUBUzNXNGbh0gAgAhMjNxAAIBButNXQgAilE1VzwARqrnQATdUACJkZGSmZq5ozcOkAAAEKmZGauaABFKKUBSgKmZq5ozcOkAEAEKmZGauaABFKJkZuIABAIN1pq6EAETIzcOAQACbrTV0IAIqZkZq5oAEUopRFKBqrngAjVXOgAm6oAE1dCACZGqueN1QAJq6E1dEZGqueN1QAIAJm4E3WmroTI1VzxuqABAUSCAxgo1dCZGrojV0QAJq6IAoUmFhYWFhYWNVc8AEaq50AE3VGroTV0QAJq6IyNVc8bqgATMiMyMjACIzACACABIwAiMwAgAgASJTM1Vz4AIsJkZGSmZq5ozcOkAEAEKmZq5ozceAMbrjV0IAImroQBBMAU1dEAIJgCmrogBDVXPABGqudABN1RkauhMjVXPG6oAEAEyNXQmRqrnjdUACACauhABACN1hq6EAc3XGroQARY1VzwARqrnQATdUZGroTI1VzxuqABABNXQgAmRqrnjdUACZGRkpmauaM3DpABABCZkRmRkYARGYAQAQAJGAERmAEAEACRKZmqufABFhMlMzVzRmRGbrzdOZGqueN1QAIARunMjVXPG6oAEAEAQ1dCACJq6E1dEACJgBmrogAjI1VzxuqABNXQgAgBACGroQARY1VzwARqrnQATdUauhNXRACm6w1dCAGauiNXRAAmrojV0QAJkaq543VAAmroQATI1VzxuqABAENVc8AEaq50AE3VABJgYPYeZ8bAAABjFTZNODYeZ/YeZ9YHP2h94lLAr7tUL1FpYOeheGGPYOGnNTXdntt/L7/2Hmf2Hmf2HmfWBy0p49/k+9r+Nhiy2UK8ghX02IhcF1TxLfso/jU/////9h5n9h6n1gci1pQGkOHQepMdbwU5UZilwl3r9AngcO8SRHnD////wABWRjyWRjvAQAAMiIyMjIyMjIyMjIyMjIyMlMzVzSmZqrnwAhSiJkZkRmauaACABSgZkZGRgBEZgBABAAkYARGYAQAQAJEpmaq58AEUoipmRGZq5oAIAEUoGAGauhABEwAjV0QAIAJq6IAMzIyMjACIzACACABIwAiMwAgAgASJTM1Vz4AIpRFTMiMzVzQAQAIpQMAM1dCACJgBGrogAQAQAjIjMiM3Xm6cyNVc8bqgAQAjdOZGqueN1QAIAIARq6EyNVc8bqgAQATI1dCZGqueN1QAIAJq6EAITIyMjIyUzNXNGbh0gAAAhUzNXNGZEZkZGRgBEZgBABAAkYARGYAQAQAJEpmaq58AEUoCpmRmrmgARSiYAZq6EAETACNXRAAkZkRm683TmRqrnjdUACAEbpzI1VzxuqABABADNXQmRqrnjdUACACACauhMjVXPG6oAEBY3WGroQEhUzNXNGRmACRJQFKIAomZkRmZGRgBEZgBABAAkYARGYAQAQAJEZgAkSmZq5owBgAhIyIzABAFACMyIzIyMAIjMAIAIAEjACIzACACABIjMAEiMyIiMwAgBAAwAjAEABEAMAIAEAISMwBQAjMiIjMAIAQAMAMAEWABIlABIwASIyUzIjM1c0AEACKUDIyMjJTM1c0ZuHSAAACFKApRNVc8AEaq50AE3VAAmroQARMjIyMlMzVzRm4dIAAAIUoClE1VzwARqrnQATdUACauhNXRAAmRqrnjdUACACAIRKZmrmjIzABIkoClEAEVMzVzRmZkREZmZGZERmRkYARGYAQAQAJGAERmAEAEACRKZmqufABEiUAEVMzVzRm681Vzpq6EAEAQTAFNXQgAiYARq6IAEAEjIiMAIAM3VmqueABN1IAIAgAJGZEZuvN0wARumABMyIzIjNXQGbsAAgAUvW97YwG6kAI3UAAgCABgAixuuAPSIEDRlNOAEgAgDhSYWFhYWFTM1c0ZuHSACACEzMiMzIyMAIjMAIAIAEjACIzACACABIjMAEiUzNXNGAMAEJGRGYAIAoARmRGZGRgBEZgBABAAkYARGYAQAQAJEZgAkRmRERmAEAIAGAEYAgAIgBgBAAgBCRmAKAEZkREZgBACABgBgAiwAJEoAJGACRGSmZEZmrmgAgARSgZGRkZKZmrmjNw6QAAAQpQFKJqrngAjVXOgAm6oAE1dCACJkZGRkpmauaM3DpAAABClAUomqueACNVc6ACbqgATV0Jq6IAEyNVc8bqgAQAQBSJTM1c0ZGYAJElAUogAipmauaMjMAEiSgKUQBhUzNXNGZmRERmZkZkRGZGRgBEZgBABAAkYARGYAQAQAJEpmaq58AESJQARUzNXNGbrzVXOmroQAQBBMAU1dCACJgBGrogAQASMiIwAgAzdWaq54AE3UgAgCAAkZkRm683TABG6YAEzIjMiM1dAZuwACABS9b3tjAbqQAjdQACAIAGACLG64A9IgEDRlNOADIzcCkAAACkAEAcKTCwsLCpmauaM3DpACABCpmauaMyIzNXNABAApQMyIzIjJTM1c0ZGbhzVXOm6oAFIAI1dCauiACEyMjJTM1c0ZuHSAAACFKAqZmrmjNw6QAQAQmRm4gAEAc3WmroQARSiaq54AI1VzoAJuqABEyMjJTM1c0ZuHSAAACFTMjNXNAAilFKApQFTM1c0ZuHSACACFTMjNXNAAilEyM3EAAgDm601dCACJkZuHAHABN1pq6EAEVMyM1c0ACKUUoilA1VzwARqrnQATdUACauhABACMjVXPG6oAE1dCauiMjVXPG6oAEAE3WmroTV0Rkaq543VAAgLGroQDjMjIyMAIjMAIAIAEjACIzACACABIlMzVXPgAilAVMzVzRm681dCACAGKURMAI1dEACauhABN1hq6E1dEAcJmRGSmZq5ozIjJTMiMzVzQAQAIpQMjIyUzNXNGbh0gAgAhSiJm5A3XGroQAQBTVXPABGqudABN1Rq6EAETIyMlMzVzRm4dIAAAITNyAApuuNXQgAilE1VzwARqrnQATdUauhNXRAAmRqrnjdUACAEAEACJmZEZmRkYARGYAQAQAJGAERmAEAEACRGYAJEpmauaMAYAISMiMwAQBQAjMiMyMjACIzACACABIwAiMwAgAgASIzABIjMiIjMAIAQAMAIwBAARADACABACEjMAUAIzIiIzACAEADADABFgASJQASMAEiM3XgCgAgEESmZq5oyMwASJKApRABEwAiIyMlMzVzSmZEZmrmgAgARSgZkZGRgBEZgBABAAkYARGYAQAQAJEZgAkSmZGauaABFKJgCgBCYAgAIpQIwASJTMiMzVzQAQAIpQMyIzdebpgAjdMACAEAOJm68AUAEA0TABMyIzIjN0qQABmroAAjNXQAApeuAzdKkAAZq6A3UgBJeuA1dCauiMjVXPG6oAEAEAcAgVMzVzRmZkREZmZGZERmRkYARGYAQAQAJGAERmAEAEACRKZmqufABEiUAEVMzVzRm681Vzpq6EAEAQTAFNXQgAiYARq6IAEAEjIiMAIAM3VmqueABN1IAIAgAJGZEZuvN0wARumABMyIzIjNXQGbsAAgAUvW97YwG6kAI3UAAgCABgAixuuAWMjNxSRBA0ZTTgAAEAdIAIBUUmFhYjMjIyMAIjMAIAIAEjACIzACACABIjMAEiUzIzVzQAIpRMAUAITAEABFKBGACRGbrwBAAQDTMiMyIzdKkAAZq6AAIzV0AAKXrgNXQmRqrnjdUACAEZulSAAM1dAbqQAUvXAAMAKLCxuuACNXQgAmroTV0QAIsJkpmauaMjIyUzNXNGRm4c1VzpuqABSACNXQmrogAhMjIyUzNXNGbh0gAAAhSgKmZq5ozcOkAIAEKURM3EADm601dCACaq54AI1VzoAJuqABEyMjJTM1c0ZuHSAAACFTMjNXNAAilFKApQFTM1c0ZuHSACACFTMjNXNAAilEzcQAObrTV0IAImRm4cAgAE3WmroQARUzIzVzQAIpRSiKUDVXPABGqudABN1QAJq6EAEyNVc8bqgATV0Jkaq543VAAgAmroQDxMjIzMAgiMwASJhACFiIyUzNXNGbhzMiMzMjIiIzIyMAIjMAIAIAEjACIzACACABIlMzVXPgAiAKKmZq5ozdeaq501dCACAMJgCGqueNXQgAiYARq6IAEAE3UgAgBJAAEZmZGRERmRkYARGYAQAQAJGAERmAEAEACRKZmqufABEAUVMzVzRm681Vzpq6EAEAYTAENVc8auhABEwAjV0QAIAJupABACSAAdabqwAQBTdcAmACkAEKmZq5ozMyIiMzMjMiIzIyMAIjMAIAIAEjACIzACACABIlMzVXPgAiRKACKmZq5ozdeaq501dCACAIJgCmroQARMAI1dEACACRkRGAEAGbqzVXPAAm6kAEAQAEjMiM3Xm6YAI3TAAmZEZkRmroDN2AAQAKXre9sYDdSAEbqAAQBAAwARY3XAJgAmRm4FIAAAFIAIBIVMzVzRmRkZGAERmAEAEACRgBEZgBABAAkSmZqrnwARSgKmZq5ozdeauhABADFKImAEauiABAFN1hq6E1dEAoKmZq5ozcSkAMZGZkZGRgBEZgBABAAkYARGYAQAQAJESmZqrnwARACEzADMwBAAjV0IAJq6IAEiM3AABGSmZqrnwARSAAEzIzIyIlMzVXPgAiAEJmAGZuAACSACNXRAAkZgBABAAkZESmZqrnwARACEzADM3AABJABGrogASMwAgAgAQATNwCQACQARq6IAE3VmqueABSAAABADFJhYWFhYyM3FJFA0ZTTgAAEAM3XAAmroQAhUzNXNGZEZkRkpmauaMjNw5qrnTdUACkAEauhNXRABCZGRkpmauaM3DpAAABClAVMzVzRm4dIAIAITIzcQACAObrTV0IAIpRNVc8AEaq50AE3VAAiZGRkpmauaM3DpAAABCpmRmrmgARSilAUoCpmauaM3DpABABCpmRmrmgARSiZGbiAAQBzdaauhABEyM3DgDgAm601dCACKmZGauaABFKKURSgaq54AI1VzoAJuqABNXQgAgBGRqrnjdUACauhNXRGRqrnjdUACACACauhAPEzIjJTM1c0ZkRkpmRGZq5oAIAEUoGRkZKZmrmjNw6QAQAQpREzcgbrjV0IAIApqrngAjVXOgAm6o1dCACJkZGSmZq5ozcOkAAAEJm5AAU3XGroQARSiaq54AI1VzoAJuqNXQmrogATI1VzxuqABACACABEyMzABMyIzIjN0qQABmroAAjNXQAApeuA1dCZGqueN1QAIARm6VIAAzV0BupABS9cAAYAQBREZmAGZkRmRGbpUgADNXQABGaugABS9cBm6VIAAzV0BupACS9cBq6E1dEZGqueN1QAIAIAgAoAJEpmauaMjMAEiSgKUQARMjAFIjAFIlMzVzRmRkZGAERmAEAEACRgBEZgBABAAkRmACRKZkZq5oAEUomAKAEJgCAAilAjABIlMyIzNXNABAAilAzdeAeACJmRGbrzdMAEbpgAQBwAgEhUzNXNGZmRERmZkZkRGZGRgBEZgBABAAkYARGYAQAQAJEpmaq58AESJQARUzNXNGbrzVXOmroQAQBBMAU1dCACJgBGrogAQASMiIwAgAzdWaq54AE3UgAgCAAkZkRm683TABG6YAEzIjMiM1dAZuwACABS9b3tjAbqQAjdQACAIAGACLG64BsAUyM3ApAAAApABANCpmauaMyMjIwAiMwAgAgASMAIjMAIAIAEiUzNVc+ACKUBUzNXNGbrzV0IAIAYpREwAjV0QAIBpusNXQmrogHBMlMzVzSmZq5ozIjMiMlMzVzRkZuHNVc6bqgAUgAjV0Jq6IAITIyMlMzVzRm4dIAAAIUoCpmauaM3DpABABCZGbiAAQBzdaauhABFKJqrngAjVXOgAm6oAETIyMlMzVzRm4dIAAAIVMyM1c0ACKUUoClAVMzVzRm4dIAIAIVMyM1c0ACKUTIzcQACAObrTV0IAImRm4cAcAE3WmroQARUzIzVzQAIpRSiKUDVXPABGqudABN1QAJq6EAEAIyNVc8bqgATV0Jq6IyNVc8bqgAQATNwJutNXQgApBAeFkpGroQHRSiJmRkZGAERmAEAEACRgBEZgBABAAkSmZqrnwARSgKmZGauaABFKJgBmroQARMAI1dEACSmZEZmrmgAgARSgZuvNXQmRqrnjdUACACauhNXRABCZuJMyIzcAZuDACABUzNXNGbiEgADNwwAQAIpABCkAAZKZmqufABFIAATM1c0ZuvNVc6auhABN1KRQA3WmqueNXQm6s1Vzxq6EAFIAAARIAgyUzNVc+ACKQAAmZq5ozdeaq501dCACbqUiEAN1pqrnjV0JurNVc8auhABSAAN1Zq6E1dEZGqueN1QAIAJusNXQgPikwsauiMjVXPG6oAECQWFhYyM3FJEQNGU04AABAGFiIzIjMyMjACIzACACABIwAiMwAgAgASIzABIlMzVzRgDABCRkRmACAKAEZkRmRkYARGYAQAQAJGAERmAEAEACRGYAJEZkREZgBACABgBGAIACIAYAQAIAQkZgCgBGZERGYAQAgAYAYAIsACRKACRgAkRm68AUAEAEWN1wARq6EAI1dCauiACFjdaauhNXRGRqrnjdUACAsaq54AI1VzoAJuqASMyMjIwAiMwAgAgASMAIjMAIAIAEiMwASIzIiIzACAEADMAUAIwBAARIlABIyMjIzMjIwAiMwAgAgASMAIjMAIAIAEiUzNVc+ACLCpmauaM3XmqudNXQgAgHiZGRGYAIAYARutNVc8ZKZmqufABFhMlMzVXPgAiauhACFjV0QAJurNVc8auhACN1xqrnTJTM1Vz4AIsJkpmaq58AETV0IAQsauiABN1ZqrnjV0IAImAEauiABABIjIyMlMzVzRm4dIAQAITJTM1c0ZuvMlMzVXPgAiwmSmZqrnwARNXQgBCxq6IAEyMyMjIwAiMwAgAgASMAIjMAIAIAEiUzNVc+ACKXrgEzV0BgBmroQATACNXRAAqrnTMjIyMAIjMAIAIAEjACIzACACABIlMzVXPgAil63vbGATJTM1c0YAgAImaugABMAM1dEAEJgBmrogAjV0IAJGZGRkYARGYAQAQAJGAERmAEAEACRKZmqufABFKAqZkZq5oAEUomAGauhABEwAjV0QAJGAGbrjVXOgAm6s1VzwAIBBGbjzM3GJAAG41IkDRlNOAAAUiQNGU04AATFTM1c0ZuHMjMyMjIwAiMwAgAgASMAIjMAIAIAEiJTM1Vz4AIgBCZgBmYAgARq6EAE1dEACRGbgAAjJTM1Vz4AIpAACZkZkZESmZqrnwARACEzADM3AABJABGrogASMwAgAgASMiJTM1Vz4AIgBCZgBmbgAAkgAjV0QAJGYAQAQAIAJm4BIABIAI1dEACbqzVXPAApAAAAgDpACCpmauaM3DgCpABCpmauaMjIyMlMzVzRm4dIAAAITIyMlMzVzRm4dIAAAITNyBuuNXQgCG641dCACKUTVXPABGqudABN1Rq6E1dEAIKUTVXPABGqudABN1Rq6EAEyNVc8bqgAQARUzNXNGZEZgBEZgBEZuPACABFKAmYAJJQFKKmZq5ozceZm4xIAA3GpEQNGU04AAGSJA0ZTTgAVMzVzRm4hIAY3GgDCZERgBABmZuMSAGM3Am40AZIAYAYSJQARYyMjJTM1c0ZuHSAAACEyIjACADN1xq6EAESJQATVXPABGqudABN1Rq6EyNVc8bqgAQARUzNXNGbiUggKToAzJTM1Vz4AIpAACZmrmjN15qrnTV0IAJupSIEAN1pqrnjV0JurNVc8auhABSAAAHEjMAEAgAIWFhYWFhYyMjIyMgBVMzVzRm4dIAAAITIyMjIyMjIyMjIyMjJJlMzVXPgAikwsauiAGUzNXNGbh0gAAAhMjIyMkmUzNVc+ACKTCxq6IAM3XAAmroQARUzNXNGbh0gAgAhMkmUzNVc+ACKTCwsaq54AI1VzoAJuqABNXQgAmrogBlMzVzRm4dIAAAITIyMjJJlMzVXPgAikwsauiADN1wAJq6EAEVMzVzRm4dIAIAITJJlMzVXPgAikwsLGqueACNVc6ACbqgATV0IAIsaq54AI1VzoAJuqABNXQgAixqrngAjVXOgAm6o1dCauiAEN1Zq6EAE1dEACZGqueN1QAIAJkpmaq58AESJQARMyIiMwAgBAAzV0IAJkZkZEpmaq58AESJQARMyIiMwAgBAAzV0IAJgBGrogASMwAgAgASMiUzNVc+ACJEoAImZERGYAQAgAZq6EAEwAjV0QAJGYAQAQAIAJq6IAEAIzIyMjACIzACACABIwAiMwAgAgASIzABIjMiIjMAIAQAMwBQAjAEABEiUAEjIyMiMwAQAgAzdWauhACMjIyMlMzVzRm4dIAQAITV0IAIsaq54AI1VzoAJuqABNXQmrogATV0Rkaq543VAAgAmSmZqrnwARIlABEzIiIzACAEADNXQgAmRmRkSmZqrnwARIlABEzIiIzACAEADNXQgAmAEauiABIzACACABIyJTM1Vz4AIkSgAiZkREZgBACABmroQATACNXRAAkZgBABAAgAmrogAQAhYwAjdYauhAKMAEAIyMjIwAiMwAgAgASMAIjMAIAIAEiUzNVc+ACKXrgEyUzNXNGAIACJmroAATADNXRABCYAZq6IAI1dCACRgBm6s1dCauiMjVXPG6oAEAEwAjdYauhAJIzIyMjACIzACACABIwAiMwAgAgASJTM1Vz4AIpQFTMjNXNAAilEwAzV0IAImAEauiABIzdeaq50AEAUAEyMjIwAiMwAgAgASMAIjMAIAIAEiUzNVc+ACKXrgEzV0BgBmroQATACNXRAAkauhNXRGRqrnjdUACACZGZGRGZGRgBEZgBABAAkYARGYAQAQAJEpmaq58AEUvW97YwCZGZgCmqueNXQgBEZq6Azdgaq501dCAGACAEIAJgBGrogAQASMzACN1YAJGREYAQAZumABEiUAEjJTM1c0aujABEiUAESIwAgAzMiMyMjACIzACACABIwAiMwAgAgASJTM1Vz4AIpet72xgEyMzAFNVc8auhACIzV0Bm7A1Vzpq6EAMAEAIQATACNXRAAgAkpmauaM3XgAm6hIAASJQARIjACADABABN1Zq6EAMyMjJTM1c0ZuHSAAACE1dCACLGqueACNVc6ACbqjV0Jq6IAYyNXRGrogATV0QAJq6I1dEACauiABNXRAAmRqrnjdUACauhABMjVXPG6oAEAFMAYjYeZ/YeZ/YeZ9YIDQfeBHTSvW35cersb8X93LaStUJTwa0orwBXyg/Xm++/wD/GwAAAYxU2TTg2Hmf2HmfWBz9ofeJSwK+7VC9RaWDnoXhhj2DhpzU13Z7bfy+/9h5n9h5n9h5n1gctKePf5Pva/jYYstlCvIIV9NiIXBdU8S37KP41P//////AAH19g=="
                        }
                    ],
                    "header": {
                        "blockHash": "c3165b4227f606e5413d3de42adde15545d6abc69cd369837645ec626045a6a0",
                        "blockHeight": 1473439,
                        "blockSize": 10394,
                        "issuerVK": "5e05cc4b4d924ba37c5c8024260119857ccec07e9591799110c31e3678a942c9",
                        "issuerVrf": "KdRNQkWCBedfFT8zaNvVFtQ+f5z/56HdHsFf5HTXSrY=",
                        "opCert": {
                            "count": 5,
                            "hotVk": "acqmfMxyrtitd6te4Y18rt8dK+bfUEwPLW5XQmX+fiU=",
                            "kesPeriod": 259,
                            "sigma": "1fSWIGKTvybiYWUjNw5ppYollqt7VqsFvuemKlvV06VAw+4KNP1Y28TojGoi7iXXkSqk58wqQ2O+0XK6kLI/CQ=="
                        },
                        "prevHash": "0e32b314d5244e020fb97f38a97b63368bf5c4bc614f738d0b6939888ae26853",
                        "protocolVersion": {
                            "major": 9,
                            "minor": 0
                        },
                        "signature": "fsWZuxW1ayN9NP+PJSBM+fqciGSHZbakd4QVRyhmGe1zo1f8UAXE9c2QXJjiSPLkpgkywqZCu0WbX95sAd1HDOnRIqfJ4Mv//eMUpG+7K3cBo0Pgcz7J0OoWkw5Uc6YscagXUr2MyOsew00O7KnE2+if7JRmB2E3/8+C2ej5LYFUaRRPJiA/K9y7ZgtZQ8SQQTv823tizqPQyrinfe893rQYe27IbFIBlVun4LWrA2PcFx+FII6AdT5DqwUXuWUdZqKlMmpg4yhuqqViV4pMCCtlx7GnP/uCvSS+6auzAks9TzjdVr1xwsalnedEeNcRGioMzNxTN9mbwC5iSRNxmHBPChhHig0zaOfAh6zsA1Elt8z0X2+gli7FUnM8KxRS0dpEc640hTgM/27ZQsXiWoYjeC0xaHHIZLznMC/pTwqO94nTL6QpHM8aZT+yoZqWYlbUGTl9rJ/N5X2jUioJ7VhrnH7xv4XrJbGrAqgX/RFhLedbmIlp2LBSL17NbO3gi1SC/j1QgMY1rWK3Yy8fcj0Cw/DWFP2Ia3+05FnauZujtVZHDfTY9qriAchR5Gah02hArsvqxeNrxmmED+emZA==",
                        "slot": 34012382
                    },
                    "headerHash": "96ae468b77f51d1d8b0086d92debfd89904974f23db932e72dcedbc8edc6b011"
                }
            },
            "tip": {
                "blockNo": 1473439,
                "hash": "96ae468b77f51d1d8b0086d92debfd89904974f23db932e72dcedbc8edc6b011",
                "slot": 34012382
            }
        }
    }`
	var result CompatibleResult
	err := json.Unmarshal([]byte(example), &result)
	assert.Nil(t, err)
}

func ValueChecks(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		equal1 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		equal2 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		if !shared.Equal(equal1, equal2) {
			t.Fatalf("%v and %v are not equal", equal1, equal2)
		}

		val1 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000001)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		val2 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		val3 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(12345678900)},
		)
		val4 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		if !shared.GreaterThan(val1, val2) {
			t.Fatalf("%v is not greater than %v", val1, val2)
		}
		if !shared.LessThan(val2, val1) {
			t.Fatalf("%v is not less than %v", val1, val2)
		}
		if !shared.GreaterThan(val3, val4) {
			t.Fatalf("%v is not greater than %v", val3, val4)
		}
		if !shared.LessThan(val4, val3) {
			t.Fatalf("%v is not less than %v", val4, val3)
		}
		if ok, err := shared.Enough(val3, val4); !ok {
			t.Fatalf("%v does not have enough assets for %v: %v", val4, val3, err)
		}
		if shared.Equal(val1, val2) {
			t.Fatalf("%v and %v are equal", val1, val2)
		}
		if shared.Equal(val3, val4) {
			t.Fatalf("%v and %v are equal", val3, val4)
		}

		val5 := shared.Add(val1, val2)
		val6 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(2000001)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(2469135780)},
		)
		if !shared.Equal(val5, val6) {
			t.Fatalf("%v is not the expected value (%v)", val5, val6)
		}

		val7 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(600000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(2345678900)},
		)
		val8 := shared.Subtract(val3, val7)
		val9 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(400000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1000000000)},
		)
		if !shared.Equal(val8, val9) {
			t.Fatalf("%v is not the expected value (%v)", val8, val9)
		}
	})
}
