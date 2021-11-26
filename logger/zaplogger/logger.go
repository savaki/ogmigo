package zaplogger

import (
	"github.com/savaki/ogmigo"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func Wrap(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(message string, kvs ...ogmigo.KeyValue) {
	l.logger.Debug(
		message, getFields(kvs)...)
}

func (l *Logger) Info(message string, kvs ...ogmigo.KeyValue) {
	l.logger.Info(message, getFields(kvs)...)
}

func (l *Logger) With(kvs ...ogmigo.KeyValue) ogmigo.Logger {
	return &Logger{
		logger: l.logger.With(getFields(kvs)...),
	}
}

func getFields(kvs []ogmigo.KeyValue) []zap.Field {
	var fields []zap.Field
	for _, kv := range kvs {
		fields = append(fields, zap.String(kv.Key, kv.Value))
	}
	return fields
}
