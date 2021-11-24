package num

import (
	"fmt"
	"math/big"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Int big.Int

func Int64(v int64) Int {
	bi := big.NewInt(v)
	return Int(*bi)
}

func New(s string) (Int, bool) {
	bi, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return Int{}, false
	}

	return Int(*bi), true
}

func (i Int) Add(that Int) Int {
	sum := big.NewInt(0).Add(i.BigInt(), that.BigInt())
	return Int(*sum)
}

func (i Int) BigInt() *big.Int {
	bi := big.Int(i)
	return &bi
}

func (i Int) Int() int {
	return int(i.BigInt().Int64())
}

func (i Int) Int64() int64 {
	return i.BigInt().Int64()
}

func (i Int) MarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	item.N = aws.String(i.BigInt().String())
	return nil
}

func (i Int) MarshalJSON() ([]byte, error) {
	s := i.BigInt().String()
	return []byte(s), nil
}

func (i Int) String() string {
	return i.BigInt().String()
}

func (i Int) Sub(that Int) Int {
	sum := big.NewInt(0).Sub(i.BigInt(), that.BigInt())
	return Int(*sum)
}

func (i *Int) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if aws.BoolValue(item.NULL) {
		return nil
	}
	if item.N == nil {
		return fmt.Errorf("unable to unmarshal invalid Int: N not set")
	}

	s := aws.StringValue(item.N)
	v, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return fmt.Errorf("failed to parse number, %v", s)
	}

	number := Int(*v)
	*i = number

	return nil
}

func (i *Int) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}

	s := string(data)
	v, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return fmt.Errorf("failed to parse number, %v", s)
	}

	number := Int(*v)
	*i = number

	return nil
}
