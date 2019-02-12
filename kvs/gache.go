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
	"errors"
	"strconv"
	"unsafe"

	"github.com/kpango/gache"
)

// BoltDB is one implementation of KVS
type Memory struct {
	vk gache.Gache
	kv gache.Gache
}

// NewBoltDB returns BoltDB instance
func NewMemory() *Memory {
	return &Memory{
		kv: gache.New().SetDefaultExpire(0),
		vk: gache.New().SetDefaultExpire(0),
	}
}

func (m *Memory) GetKey(val uint) ([]byte, error) {
	b, ok := m.vk.Get(strconv.FormatUint(*(*uint64)(unsafe.Pointer(&val)), 10))
	if !ok {
		return nil, errors.New("not found")
	}
	return b.([]byte), nil
}

// GetKeys wraps multiple calls GetKey
func (m *Memory) GetKeys(vals []uint) ([][]byte, error) {
	ret := make([][]byte, len(vals))
	for i, val := range vals {
		k, err := m.GetKey(val)
		if err != nil {
			return nil, err
		}
		ret[i] = k
	}
	return ret, nil
}

func (m *Memory) GetVal(key []byte) (uint, error) {
	b, ok := m.kv.Get(*(*string)(unsafe.Pointer(&key)))
	if !ok {
		return 0, errors.New("not found")
	}
	return b.(uint), nil
}

func (m *Memory) Set(key []byte, val uint) error {
	m.kv.Set(*(*string)(unsafe.Pointer(&key)), val)
	m.vk.Set(strconv.FormatUint(*(*uint64)(unsafe.Pointer(&val)), 10), key)
	return nil
}

func (m *Memory) Delete(key []byte) error {
	val, err := m.GetVal(key)
	if err != nil {
		return err
	}
	m.kv.Delete(*(*string)(unsafe.Pointer(&key)))
	m.vk.Delete(strconv.FormatUint(*(*uint64)(unsafe.Pointer(&val)), 10))
	return nil
}

func (m *Memory) Close() error {
	m.kv.Clear()
	m.vk.Clear()
	return nil
}
