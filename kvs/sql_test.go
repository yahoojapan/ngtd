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

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func initSQL(t *testing.T) *SQL {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Unexpected Error: initSQL(): %v", err)
	}
	s, err := NewSQL(db)
	if err != nil {
		t.Fatalf("Unexpected Error: initSQL(): %v", err)
	}
	return s
}

func TestSQL(t *testing.T) {
	t.Parallel()
	t.Run("TestGetKey", func(t *testing.T) {
		s := initSQL(t)
		defer SetupWithTeardown(s, t)()
		GetKey(s, t)
	})

	t.Run("TestGetKeys", func(t *testing.T) {
		s := initSQL(t)
		defer SetupWithTeardown(s, t)()
		GetKeys(s, t)
	})

	t.Run("TestGetVal", func(t *testing.T) {
		s := initSQL(t)
		defer SetupWithTeardown(s, t)()
		GetVal(s, t)
	})

	t.Run("TestSet", func(t *testing.T) {
		s := initSQL(t)
		Set(s, t)
		s.Close()
		os.RemoveAll("test")
	})

	t.Run("TestDelete", func(t *testing.T) {
		s := initSQL(t)
		defer SetupWithTeardown(s, t)()
		Delete(s, t)
	})

	t.Run("TestClose", func(t *testing.T) {
		s := initSQL(t)
		defer SetupWithTeardown(s, t)()
		Close(s, t)
	})
}
