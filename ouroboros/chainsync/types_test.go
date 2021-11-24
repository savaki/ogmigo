package chainsync

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fxamacker/cbor"
	"github.com/nsf/jsondiff"
	"github.com/tj/assert"
)

func TestUnmarshal(t *testing.T) {
	err := filepath.Walk("../../ext/ogmios/server/test/vectors/ChainSync/Response/RequestNext", assertStructMatchesSchema(t))
	assert.Nil(t, err)
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
		assert.Nil(t, err)
		defer f.Close()

		decoder := json.NewDecoder(f)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&Response{})
		assert.Nil(t, err, fmt.Sprintf("struct did not match schema for file, %v", path))

		return nil
	}
}

func TestDynamodbSerialize(t *testing.T) {
	err := filepath.Walk("../../ext/ogmios/server/test/vectors/ChainSync/Response/RequestNext", assertDynamoDBSerialize(t))
	assert.Nil(t, err)
	decoder := json.NewDecoder(nil)
	decoder.DisallowUnknownFields()
}

func assertDynamoDBSerialize(t *testing.T) filepath.WalkFunc {
	return func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		path, _ = filepath.Abs(path)
		f, err := os.Open(path)
		assert.Nil(t, err)
		defer f.Close()

		var want Response
		err = json.NewDecoder(f).Decode(&want)
		assert.Nil(t, err)

		item, err := dynamodbattribute.Marshal(want)
		assert.Nil(t, err)

		var got Response
		err = dynamodbattribute.Unmarshal(item, &got)
		assert.Nil(t, err)

		w, err := json.Marshal(want)
		assert.Nil(t, err)

		g, err := json.Marshal(got)
		assert.Nil(t, err)

		opts := jsondiff.DefaultConsoleOptions()
		diff, s := jsondiff.Compare(w, g, &opts)
		if diff == jsondiff.FullMatch {
			return nil
		}

		fmt.Println(s)
		assert.Equal(t, jsondiff.FullMatch, diff)

		return nil
	}
}

func TestPoint_CBOR(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := PointString("origin")
		item, err := cbor.Marshal(want.Point(), encOptions)
		assert.Nil(t, err)

		var point Point
		err = cbor.Unmarshal(item, &point)
		assert.Nil(t, err)
		assert.Equal(t, PointTypeString, point.PointType())

		got, ok := point.PointString()
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})

	t.Run("struct", func(t *testing.T) {
		want := &PointStruct{
			BlockNo: 123,
			Hash:    "hash",
			Slot:    456,
		}
		item, err := cbor.Marshal(want.Point(), encOptions)
		assert.Nil(t, err)

		var point Point
		err = cbor.Unmarshal(item, &point)
		assert.Nil(t, err)
		assert.Equal(t, PointTypeStruct, point.PointType())

		got, ok := point.PointStruct()
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})
}

func TestPoint_DynamoDB(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := PointString("origin")
		item, err := dynamodbattribute.Marshal(want.Point())
		assert.Nil(t, err)

		var point Point
		err = dynamodbattribute.Unmarshal(item, &point)
		assert.Nil(t, err)
		assert.Equal(t, PointTypeString, point.PointType())

		got, ok := point.PointString()
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})

	t.Run("struct", func(t *testing.T) {
		want := &PointStruct{
			BlockNo: 123,
			Hash:    "hash",
			Slot:    456,
		}
		item, err := dynamodbattribute.Marshal(want.Point())
		assert.Nil(t, err)

		var point Point
		err = dynamodbattribute.Unmarshal(item, &point)
		assert.Nil(t, err)
		assert.Equal(t, PointTypeStruct, point.PointType())

		got, ok := point.PointStruct()
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})
}

func TestPoint_JSON(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		want := PointString("origin")
		data, err := json.Marshal(want.Point())
		assert.Nil(t, err)

		var point Point
		err = json.Unmarshal(data, &point)
		assert.Nil(t, err)
		assert.Equal(t, PointTypeString, point.PointType())

		got, ok := point.PointString()
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})

	t.Run("struct", func(t *testing.T) {
		want := &PointStruct{
			BlockNo: 123,
			Hash:    "hash",
			Slot:    456,
		}
		data, err := json.Marshal(want.Point())
		assert.Nil(t, err)

		var point Point
		err = json.Unmarshal(data, &point)
		assert.Nil(t, err)
		assert.Equal(t, PointTypeStruct, point.PointType())

		got, ok := point.PointStruct()
		assert.True(t, ok)
		assert.Equal(t, want, got)
	})
}

func TestTxID_Index(t *testing.T) {
	assert.Equal(t, 3, TxID("a#3").Index())
}

func TestTxID_TxHash(t *testing.T) {
	assert.Equal(t, "a", TxID("a#3").TxHash())
}