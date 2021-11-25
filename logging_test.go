package ogmigo

import (
	"testing"
)

func Test_DefaultLogger_print(t *testing.T) {
	DefaultLogger.Info(nil, "test", KV("key", "value"))
}

func Test_NopLogger_print(t *testing.T) {
	NopLogger.Info(nil, "test", KV("key", "value"))
}
