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
	"strconv"

	"github.com/go-redis/redis"
)

const (
	base = 36
)

func toRedisVal(v uint) string {
	return strconv.FormatUint(uint64(v), base)
}

func fromRedisVal(v string) (uint, error) {
	val, err := strconv.ParseUint(v, base, 32)
	return uint(val), err
}

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
		kv:     kv,
		vk:     vk,
	}, nil
}

func (r *Redis) GetKey(val uint) ([]byte, error) {
	pipe := r.client.TxPipeline()
	pipe.Select(r.vk)
	key := pipe.Get(toRedisVal(val))
	if _, err := pipe.Exec(); err != nil {
		return nil, err
	}
	return key.Bytes()
}

func (r *Redis) GetKeys(vals []uint) ([][]byte, error) {
	strVals := make([]string, len(vals))
	for i, v := range vals {
		strVals[i] = toRedisVal(v)
	}

	pipe := r.client.TxPipeline()
	pipe.Select(r.vk)
	keys := pipe.MGet(strVals...)
	if _, err := pipe.Exec(); err != nil {
		return nil, err
	}
	response := keys.Val()
	byteKeys := make([][]byte, len(response))
	for i, k := range response {
		if xi, ok := k.(string); ok {
			byteKeys[i] = []byte(xi)
		}
	}
	return byteKeys, nil
}

func (r *Redis) GetVal(key []byte) (uint, error) {
	pipe := r.client.TxPipeline()
	pipe.Select(r.kv)
	val := pipe.Get(string(key))
	if _, err := pipe.Exec(); err != nil {
		return 0, err
	}
	return fromRedisVal(val.Val())
}

func (r *Redis) Set(key []byte, val uint) error {
	v := toRedisVal(val)
	pipe := r.client.TxPipeline()
	pipe.Select(r.kv)
	kv := pipe.Set(string(key), v, 0)
	pipe.Select(r.vk)
	vk := pipe.Set(v, key, 0)
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
	vk := pipe.Del(toRedisVal(val))
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
