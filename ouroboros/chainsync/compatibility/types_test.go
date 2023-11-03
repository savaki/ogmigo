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
)

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}

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
	if err != nil {
		t.Fatalf("error unmarshalling response: %v", err)
	}
}
