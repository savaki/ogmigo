package ogmigo

import (
	"bytes"
	"log"
)

type KeyValue struct {
	Key   string
	Value string
}

func KV(key, value string) KeyValue {
	return KeyValue{
		Key:   key,
		Value: value,
	}
}

type Logger interface {
	Debug(message string, kvs ...KeyValue)
	Info(message string, kvs ...KeyValue)
	With(kvs ...KeyValue) Logger
}

// DefaultLogger logs via the log package
var DefaultLogger = defaultLogger{}

type defaultLogger struct {
	kvs []KeyValue
}

func (d defaultLogger) print(message string, kvs ...KeyValue) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(message)
	if len(kvs) > 0 {
		buf.WriteString(":")
	}
	for _, kv := range kvs {
		buf.WriteString(" ")
		buf.WriteString(kv.Key)
		buf.WriteString("=")
		buf.WriteString(kv.Value)
	}
	log.Println(buf)
}

func (d defaultLogger) Debug(message string, kvs ...KeyValue) {
	d.print(message, kvs...)
}

func (d defaultLogger) Info(message string, kvs ...KeyValue) {
	d.print(message, kvs...)
}

func (d defaultLogger) With(kvs ...KeyValue) Logger {
	return defaultLogger{
		kvs: append(d.kvs, kvs...),
	}
}

// NopLogger logs nothing
var NopLogger = nopLogger{}

type nopLogger struct {
}

func (n nopLogger) Debug(string, ...KeyValue) {}
func (n nopLogger) Info(string, ...KeyValue)  {}
func (n nopLogger) With(...KeyValue) Logger   { return n }
