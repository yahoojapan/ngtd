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

	"fmt"
	"github.com/boltdb/bolt"
)

var (
	kvBoltBucketName = []byte("kv")
	vkBoltBucketName = []byte("vk")
)

// BoltDB is one implementation of KVS
type BoltDB struct {
	db *bolt.DB
}

// NewBoltDB returns BoltDB instance
func NewBoltDB(p string) (*BoltDB, error) {
	db, err := bolt.Open(p, 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(kvBoltBucketName); err != nil {
			return errors.New("cannot create bucket")
		}
		if _, err := tx.CreateBucketIfNotExists(vkBoltBucketName); err != nil {
			return errors.New("cannot create bucket")
		}
		return nil
	})
	return &BoltDB{
		db: db,
	}, nil
}

func (b *BoltDB) get(boltBucketName []byte, key []byte) ([]byte, error) {
	var value []byte
	if err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBucketName)
		if bucket == nil {
			return errors.New("BoltDB Bucket NotFound")
		}
		value = bucket.Get(key)
		return nil
	}); err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BoltDB) GetKey(val uint) ([]byte, error) {
	return b.get(vkBoltBucketName, ToBytes(val))
}

func (b *BoltDB) GetVal(key []byte) (uint, error) {
	val, err := b.get(kvBoltBucketName, key)
	if len(val) != 4 {
		return 0, fmt.Errorf("key not found: %v", key)
	}
	if err != nil {
		return 0, err
	}

	return ToInt(val), nil
}

func (b *BoltDB) set(boltBucketName, key, val []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(boltBucketName)
		if err != nil {
			return errors.New("Create BoltDB Bucket Failed")
		}
		return bucket.Put(key, val)
	})
}

func (b *BoltDB) Set(key []byte, val uint) error {
	v := ToBytes(val)
	if err := b.set(kvBoltBucketName, key, v); err != nil {
		return err
	}
	return b.set(vkBoltBucketName, v, key)
}

func (b *BoltDB) del(boltBucketName, key []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBucketName)
		if bucket == nil {
			return errors.New("BoltDB Bucket NotFound")
		}
		return bucket.Delete(key)
	})
}

func (b *BoltDB) Delete(key []byte) error {
	val, err := b.GetVal(key)
	if err != nil {
		return err
	}
	if err := b.del(vkBoltBucketName, ToBytes(val)); err != nil {
		return err
	}
	return b.del(kvBoltBucketName, key)
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}
