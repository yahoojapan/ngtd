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
	"encoding/binary"
)

type KVS interface {
	GetKey(uint) ([]byte, error)
	GetVal([]byte) (uint, error)
	Set([]byte, uint) error
	Delete([]byte) error
	Close() error
}

var (
	byteOrder = binary.LittleEndian
)

// ToBytes convert integer to byte array
func ToBytes(i uint) []byte {
	key := make([]byte, 4)
	byteOrder.PutUint32(key, uint32(i))
	return key
}

// ToInt converts byte array to integer
func ToInt(bs []byte) uint {
	return uint(byteOrder.Uint32(bs))
}
