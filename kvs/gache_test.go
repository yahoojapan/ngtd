//
// Copyright (C) 2018 Yahoo Japan Corporation
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
//

package kvs

import (
	"os"
	"testing"
)

func initMemory(t *testing.T) KVS {
	return NewMemory()
}

func TestMemory(t *testing.T) {
	t.Parallel()
	t.Run("TestGetKey", func(t *testing.T) {
		g := initMemory(t)
		defer SetupWithTeardown(g, t)()
		GetKey(g, t)
	})

	t.Run("TestGetKeys", func(t *testing.T) {
		g := initMemory(t)
		defer SetupWithTeardown(g, t)()
		GetKeys(g, t)
	})

	t.Run("TestGetVal", func(t *testing.T) {
		g := initMemory(t)
		defer SetupWithTeardown(g, t)()
		GetVal(g, t)
	})

	t.Run("TestSet", func(t *testing.T) {
		g := initMemory(t)
		Set(g, t)
		g.Close()
		os.RemoveAll("test")
	})

	t.Run("TestDelete", func(t *testing.T) {
		g := initMemory(t)
		defer SetupWithTeardown(g, t)()
		Delete(g, t)
	})

	t.Run("TestClose", func(t *testing.T) {
		g := initMemory(t)
		defer SetupWithTeardown(g, t)()
		Close(g, t)
	})
}
