package zaplogger

import (
	"github.com/savaki/ogmigo"
	"go.uber.org/zap"
	"testing"
)

func TestLogger_Debug(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	l := Wrap(logger)
	l.Debug("debug", ogmigo.KV("foo", "bar"))
	l.Info("info", ogmigo.KV("hello", "world"))
}
