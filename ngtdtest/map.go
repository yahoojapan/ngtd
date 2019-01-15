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

package ngtdtest

type Map struct {
	kv map[string]uint
	vk map[uint][]byte
}

func NewMap() *Map {
	return &Map{
		kv: make(map[string]uint),
		vk: make(map[uint][]byte),
	}
}

func (m *Map) GetVal(key []byte) (uint, error) {
	return m.kv[string(key)], nil
}

func (m *Map) GetKey(val uint) ([]byte, error) {
	return []byte(m.vk[val]), nil
}

func (m *Map) GetKeys(vals []uint) ([][]byte, error) {
	result := make([][]byte, len(vals))
	for i, val := range vals {
		result[i] = []byte(m.vk[val])
	}
	return result, nil
}

func (m *Map) Set(key []byte, val uint) error {
	m.kv[string(key)] = val
	m.vk[val] = key
	return nil
}

func (m *Map) Delete(key []byte) error {
	k := string(key)
	v := m.kv[k]
	delete(m.kv, k)
	delete(m.vk, v)
	return nil
}

func (m *Map) Close() error {
	return nil
}
