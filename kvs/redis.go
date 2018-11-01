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

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
	kv     int
	vk     int
}

func NewRedis(host, port, pass string, kv, vk int) (*Redis, error) {
	if kv == vk {
		return nil, fmt.Errorf("kv and vk must be defferent. (%d, %d)", kv, vk)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
		kv: kv,
		vk: vk,
	}, nil
}

func (r *Redis) GetKey(val uint) ([]byte, error) {
	pipe := r.client.TxPipeline()
	pipe.Select(r.vk)
	key := pipe.Get(fmt.Sprint(val))
	if _, err := pipe.Exec(); err != nil {
		return nil, err
	}
	return key.Bytes()
}

func (r *Redis) GetVal(key []byte) (uint, error) {
	pipe := r.client.TxPipeline()
	pipe.Select(r.kv)
	val := pipe.Get(string(key))
	if _, err := pipe.Exec(); err != nil {
		return 0, err
	}
	v, err := val.Uint64()
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (r *Redis) Set(key []byte, val uint) error {
	pipe := r.client.TxPipeline()
	pipe.Select(r.kv)
	kv := pipe.Set(string(key), val, 0)
	pipe.Select(r.vk)
	vk := pipe.Set(fmt.Sprint(val), key, 0)
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	if err := kv.Err(); err != nil {
		return err
	}
	return vk.Err()
}

func (r *Redis) Delete(key []byte) error {
	val, err := r.GetVal(key)
	if err != nil {
		return err
	}
	pipe := r.client.TxPipeline()
	pipe.Select(r.kv)
	kv := pipe.Del(string(key))
	pipe.Select(r.vk)
	vk := pipe.Del(fmt.Sprint(val))
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	if err := kv.Err(); err != nil {
		return err
	}
	return vk.Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
