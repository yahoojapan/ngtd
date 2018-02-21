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
	"path"

	"github.com/syndtr/goleveldb/leveldb"
)

type GoLevel struct {
	kv *leveldb.DB
	vk *leveldb.DB
}

func NewGoLevel(p string) (*GoLevel, error) {
	kv, err := leveldb.OpenFile(path.Join(p, "kv"), nil)
	if err != nil {
		return nil, err
	}
	vk, err := leveldb.OpenFile(path.Join(p, "vk"), nil)
	if err != nil {
		return nil, err
	}

	return &GoLevel{
		kv: kv,
		vk: vk,
	}, err
}

func (g *GoLevel) GetKey(val uint) ([]byte, error) {
	return g.vk.Get(ToBytes(val), nil)
}

func (g *GoLevel) GetVal(key []byte) (uint, error) {
	val, err := g.kv.Get(key, nil)
	if err != nil {
		return 0, err
	}
	return ToInt(val), nil
}

func (g *GoLevel) Set(key []byte, val uint) error {
	v := ToBytes(val)
	if err := g.kv.Put(key, v, nil); err != nil {
		return err
	}
	return g.vk.Put(v, key, nil)
}

func (g *GoLevel) Delete(key []byte) error {
	val, err := g.GetVal(key)
	if err != nil {
		return err
	}
	if err := g.kv.Delete(ToBytes(val), nil); err != nil {
		return err
	}
	return g.vk.Delete(key, nil)
}

func (g *GoLevel) Close() error {
	if err := g.kv.Close(); err != nil {
		return err
	}
	return g.vk.Close()
}
