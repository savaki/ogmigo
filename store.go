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
	"context"

	"github.com/savaki/ogmigo/ouroboros/chainsync"
)

// Store allows points to be saved and retrieved to allow graceful recovery
// after shutdown
type Store interface {
	// Save the point; save will be called multiple times and should only
	// keep track of the most recent points
	Save(ctx context.Context, point chainsync.Point) error
	// Load saved points
	Load(ctx context.Context) (chainsync.Points, error)
}
