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
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fxamacker/cbor"
	"github.com/nsf/jsondiff"
)

func TestUnmarshal(t *testing.T) {
	err := filepath.Walk("../../ext/ogmios/server/test/vectors/ChainSync/Response/RequestNext", assertStructMatchesSchema(t))
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
		err = decoder.Decode(&Response{})
		if err != nil {
			t.Fatalf("got %v; want nil: %v", err, fmt.Sprintf("struct did not match schema for file, %v", path))
		}

		return nil
	}
}

func TestDynamodbSerialize(t *testing.T) {
	err := filepath.Walk("../../ext/ogmios/server/test/vectors/ChainSync/Response/RequestNext", assertDynamoDBSerialize(t))
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
}

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

		var want Response
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

		var got Response
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
		item, err := cbor.Marshal(want.Point(), encOptions)
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
			Hash:    "hash",
			Slot:    456,
		}
		item, err := cbor.Marshal(want.Point(), encOptions)
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
			Hash:    "hash",
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
			Hash:    "hash",
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
