// Copyright 2016 The Linux Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package refs implements generic name-based reference access.
package refs

import (
	"github.com/opencontainers/image-spec/specs-go"
	"golang.org/x/net/context"
)

// Engine represents a name-based reference storage engine.
type Engine interface {

	// Put adds a new reference to the store.  The action is idempotent;
	// a nil return means "that descriptor is stored at NAME" without
	// implying "because of your Put()".
	Put(ctx context.Context, name string, descriptor *specs.Descriptor) (err error)

	// Get returns a reference from the store.  Returns os.ErrNotExist
	// if the name is not found.
	Get(ctx context.Context, name string) (descriptor *specs.Descriptor, err error)

	// List returns available names from the store.
	//
	// Results are sorted alphabetically.
	//
	// Arguments:
	//
	// * prefix: limits the result set to names starting with that
	//   value.
	// * size: limits the length of the result set to the first 'size'
	//   matches.  A value of -1 means "all results".
	// * from: shifts the result set to start from the 'from'th match.
	// * done: a channel that can be closed to cancel the list
	//   goroutine(s).
	//
	// For example, a store with names like:
	//
	// * 123
	// * abcd
	// * abce
	// * abcf
	// * abcg
	//
	// will have the following call/result pairs:
	//
	// * List(ctx, "", -1, 0) -> "123", "abcd", "abce", "abcf", "abcg"
	// * List(ctx, "", 2, 0) -> "123", "abcd"
	// * List(ctx, "", 2, 1) -> "abcd", "abce"
	// * List(ctx,"abc", 2, 1) -> "abce", "abcf"
	List(ctx context.Context, prefix string, size int, from int) (names <-chan string, errc <-chan error)

	// Delete removes a reference from the store.  Returns
	// os.ErrNotExist if the name is not found.
	Delete(ctx context.Context, name string) (err error)

	// Close releases resources held by the engine.  Subsequent engine
	// method calls will fail.
	Close() (err error)
}
