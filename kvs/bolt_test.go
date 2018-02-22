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

func initBolt(t *testing.T) *BoltDB {
	b, err := NewBoltDB(dbpath)
	if err != nil {
		t.Fatalf("Unexpected Error: initBolt(): %v", err)
	}
	return b
}

func TestBolt(t *testing.T) {
	t.Parallel()
	t.Run("TestGetKey", func(t *testing.T) {
		b := initBolt(t)
		defer SetupWithTeardown(b, t)()
		GetKey(b, t)
	})

	t.Run("TestGetVal", func(t *testing.T) {
		b := initBolt(t)
		defer SetupWithTeardown(b, t)()
		GetVal(b, t)
	})

	t.Run("TestSet", func(t *testing.T) {
		b := initBolt(t)
		Set(b, t)
		b.Close()
		os.RemoveAll("test")
	})

	t.Run("TestDelete", func(t *testing.T) {
		b := initBolt(t)
		defer SetupWithTeardown(b, t)()
		Delete(b, t)
	})

	t.Run("TestClose", func(t *testing.T) {
		b := initBolt(t)
		defer SetupWithTeardown(b, t)()
		Close(b, t)
	})
}
