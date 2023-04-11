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
	"fmt"
)

// Error encapsulates errors from ogmios
type Error struct {
	Type        string `json:"type,omitempty"`
	Version     string `json:"version,omitempty"`
	ServiceName string `json:"servicename,omitempty"`
	Fault       Fault  `json:"fault,omitempty"`
}

// Error implements error interface
func (e Error) Error() string { return fmt.Sprintf("%v: %v", e.Fault.Code, e.Fault.String) }

// Fault provides additional context for ogmios errors
type Fault struct {
	Code   string `json:"code,omitempty"`   // Code identifies error
	String string `json:"string,omitempty"` // String provides human readable description
}

type WrappedReadMessageError struct {
	message     string
	originalErr error
}

func (e *WrappedReadMessageError) Error() string {
	return fmt.Sprintf("%s: %v", e.message, e.originalErr)
}

func (e *WrappedReadMessageError) Temporary() bool {
	return true
}

func NewWrappedReadMessageError(message string, originalErr error) error {
	return &WrappedReadMessageError{
		message:     message,
		originalErr: originalErr,
	}
}
