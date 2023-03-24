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
	Error(error error, message string, kvs ...KeyValue)
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

func (d defaultLogger) Error(error error, message string, kvs ...KeyValue) {
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

func (n nopLogger) Debug(string, ...KeyValue)        {}
func (n nopLogger) Info(string, ...KeyValue)         {}
func (n nopLogger) Error(error, string, ...KeyValue) {}
func (n nopLogger) With(...KeyValue) Logger          { return n }
