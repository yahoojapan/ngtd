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
	"time"
)

func initRedis(t *testing.T) *Redis {
	r, err := NewRedis("127.0.0.1", "6379", "", 0, 1, time.Second, time.Second*300)
	if err != nil {
		t.Fatalf("Unexpected Error: initRedis(): %v", err)
	}
	return r
}

func TestRedis(t *testing.T) {
	t.Parallel()
	t.Run("TestGetKey", func(t *testing.T) {
		r := initRedis(t)
		defer SetupWithTeardown(r, t)()
		GetKey(r, t)
	})

	t.Run("TestGetKeys", func(t *testing.T) {
		r := initRedis(t)
		defer SetupWithTeardown(r, t)()
		GetKeys(r, t)
	})

	t.Run("TestGetVal", func(t *testing.T) {
		r := initRedis(t)
		defer SetupWithTeardown(r, t)()
		GetVal(r, t)
	})

	t.Run("TestSet", func(t *testing.T) {
		r := initRedis(t)
		Set(r, t)
		r.Close()
		os.RemoveAll("test")
	})

	t.Run("TestDelete", func(t *testing.T) {
		r := initRedis(t)
		defer SetupWithTeardown(r, t)()
		Delete(r, t)
	})

	t.Run("TestClose", func(t *testing.T) {
		r := initRedis(t)
		defer SetupWithTeardown(r, t)()
		Close(r, t)
	})
}
