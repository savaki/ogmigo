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

package num

import (
	"encoding/json"
	"testing"

	"github.com/tj/assert"
)

type Value struct {
	Coins Int `json:"coins,omitempty"`
}

func TestDeserializeValue(t *testing.T) {
	want := Value{Coins: Int64(123)}
	data, err := json.Marshal(want)
	assert.Nil(t, err)

	var got Value
	err = json.Unmarshal(data, &got)
	assert.Equal(t, want, got)
}

func TestMath(t *testing.T) {
	a := Int64(100)
	b := Int64(25)

	assert.Equal(t, 125, a.Add(b).Int())
	assert.Equal(t, 75, a.Sub(b).Int())
}

func TestNew(t *testing.T) {
	s, ok := New("123")
	assert.True(t, ok)
	assert.Equal(t, "123", s.String())
	assert.Equal(t, int64(123), s.Int64())
}
