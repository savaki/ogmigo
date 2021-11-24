package ogmios

import (
	"testing"
)

func Test_defaultLogger_print(t *testing.T) {
	DefaultLogger.Info(nil, "test", KV("key", "value"))
}
