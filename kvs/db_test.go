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
	"fmt"
	"os"
	"reflect"
	"testing"
)

const dbpath = "test"

func TestToBytesAndToInt(t *testing.T) {
	t.Parallel()
	for i := uint(1); i <= 0xFFFFFFFF; i <<= 1 {		
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			var I uint
			I = ToInt(ToBytes(i-1))
			if i-1 != I {
				t.Errorf("TestToBytesAndToInt(%v): %v, wanted: %v", i-1, I, i-1)
			}

			I = ToInt(ToBytes(i))
			if i != I {
				t.Errorf("TestToBytesAndToInt(%v): %v, wanted: %v", i, I, i)
			}

			I = ToInt(ToBytes(i+1))
			if i+1 != I {
				t.Errorf("TestToBytesAndToInt(%v): %v, wanted: %v", i+1, I, i+1)
			}
		})
	}
	i := uint(1 << 32-1)
	t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
		var I uint
		I = ToInt(ToBytes(i))
		if i != I {
			t.Errorf("TestToBytesAndToInt(%v): %v, wanted: %v", i, I, i)
		}
	})
}

func SetupWithTeardown(db KVS, t *testing.T) func() {
	data := []struct{
		k []byte
		v uint
	}{
		{[]byte("foo"), 1},
		{[]byte("bar"), 2},
		{[]byte("hoge"), 3},
		{[]byte("huga"), 4},
	}
	for _, d := range data {
		db.Set(d.k, d.v)
	}
	return func() {
		db.Close()
		os.RemoveAll(dbpath)
	}
}

func GetKey(db KVS, t *testing.T) {
	tests := []struct {
		val uint
		key []byte
	}{
		{1, []byte("foo")},
		{2, []byte("bar")},
		{3, []byte("hoge")},
		{4, []byte("huga")},
	}
	for _, tt := range tests {
		key, err := db.GetKey(tt.val)
		if err != nil {
			t.Errorf("Unexpected error: TestGetKey(%v) %v", tt, err)
		}
		if !reflect.DeepEqual(tt.key, key) {
			t.Errorf("TestGetKey(%v): %v, wanted: %v", tt.val, key, tt.key)
		}
	}
}

func GetKeys(db KVS, t *testing.T) {
	tests := []struct {
		vals []uint
		keys [][]byte
	}{
		{[]uint{1}, [][]byte{[]byte("foo")}},
		{[]uint{1, 2}, [][]byte{[]byte("foo"), []byte("bar")}},
		{[]uint{3, 4}, [][]byte{[]byte("hoge"), []byte("huga")}},
	}
	for _, tt := range tests {
		keys, err := db.GetKeys(tt.vals)
		if err != nil {
			t.Errorf("Unexpected error: TestGetKeys(%v) %v", tt, err)
		}
		if !reflect.DeepEqual(tt.keys, keys) {
			t.Errorf("TestGetKeys(%v): %v, wanted: %v", tt.vals, keys, tt.keys)
		}
	}
}

func GetVal(db KVS, t *testing.T) {
	tests := []struct {
		key []byte
		val uint
	}{
		{[]byte("foo"), 1},
		{[]byte("bar"), 2},
		{[]byte("hoge"), 3},
		{[]byte("huga"), 4},
	}
	for _, tt := range tests {
		val, err := db.GetVal(tt.key)
		if err != nil {
			t.Errorf("Unexpected error: TestGetVal(%v) %v", tt, err)
		}
		if tt.val != val {
			t.Errorf("TestGetVal(%v): %v, wanted: %v", tt.key, val, tt.val)
		}
	}
}

func Set(db KVS, t *testing.T) {
	tests := []struct {
		key []byte
		val uint
	}{
		{[]byte("foo"), 1},
		{[]byte("bar"), 2},
		{[]byte("hoge"), 3},
		{[]byte("huga"), 4},
	}
	for _, tt := range tests {
		if err := db.Set(tt.key, tt.val); err != nil {
			t.Errorf("Unexpected error: TestSet(%v) %v", tt, err)
		}
	}
}

func Delete(db KVS, t *testing.T) {
	tests := []struct {
		key []byte
	}{
		{[]byte("foo")},
		{[]byte("bar")},
		{[]byte("hoge")},
		{[]byte("huga")},
	}
	for _, tt := range tests {
		if err := db.Delete(tt.key); err != nil {
			t.Errorf("Unexpected error: TestDelete(%v) %v", tt, err)
		}
	}
}

func Close(db KVS, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Errorf("Unexpected error: TestClose() %v", err)
	}
}
