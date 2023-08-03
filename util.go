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

// Map provides a simple type alias
type Map map[string]interface{}

type circular struct {
	index int
	data  [][]byte
}

func newCircular(cap int) *circular {
	return &circular{
		data: make([][]byte, cap),
	}
}

func (c *circular) add(data []byte) {
	c.data[c.index] = data
	c.index = (c.index + 1) % len(c.data)
}

func (c *circular) list() (data [][]byte) {
	for i := 0; i < len(c.data); i++ {
		offset := (c.index + i) % len(c.data)
		if v := c.data[offset]; len(v) > 0 {
			data = append(data, v)
		}
	}
	return data
}

func (c *circular) prefix(data ...[]byte) [][]byte {
	return append(data, c.list()...)
}

func makePayload(methodName string, args Map) Map {
	return Map{
		"type":        "jsonwsp/request",
		"version":     "1.0",
		"servicename": "ogmios",
		"methodname":  methodName,
		"args":        args,
	}
}

func makePayloadV6(method string, params Map) Map {
	return Map{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}
}
