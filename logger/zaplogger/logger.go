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

package zaplogger

import (
	"go.uber.org/zap"

	"github.com/SundaeSwap-finance/ogmigo"
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
