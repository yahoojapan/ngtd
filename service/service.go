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

package service

import (
	"sync"

	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd/kvs"
)

type Service struct {
	db kvs.KVS
}

type SearchResult struct {
	Id       []byte
	Distance float32
	Error    error
}

var (
	once = &sync.Once{}
	s    *Service
)

func init() {
	Get()
}

func Get() *Service {
	once.Do(func() {
		s = NewService(nil)
	})
	return s
}

func (s *Service) SetDB(db kvs.KVS) {
	s.db = db
}

func SetDB(db kvs.KVS) {
	s.db = db
}

func NewService(db kvs.KVS) *Service {
	return &Service{
		db: db,
	}
}

func Search(vector []float64, size int, epsilon float32) ([]SearchResult, error) {
	return s.Search(vector, size, epsilon)
}

func (s *Service) Search(vector []float64, size int, epsilon float32) ([]SearchResult, error) {
	result, err := gongt.StrictSearch(vector, size, epsilon)
	if err != nil {
		return nil, err
	}

	vals := make([]uint, len(result))
	for i, v := range result {
		vals[i] = uint(v.ID)
	}
	ids, err := s.db.GetKeys(vals)
	if err != nil {
		return nil, err
	}
	ret := make([]SearchResult, len(result))
	for i, id := range ids {
		ret[i] = SearchResult{
			Id:       id,
			Distance: result[i].Distance,
			Error:    nil,
		}
	}
	return ret, nil
}

func SearchByID(id []byte, size int, epsilon float32) ([]SearchResult, error) {
	return s.SearchByID(id, size, epsilon)
}

func (s *Service) SearchByID(id []byte, size int, epsilon float32) ([]SearchResult, error) {
	in, err := s.db.GetVal(id)
	if err != nil {
		return nil, err
	}
	vector, err := gongt.GetStrictVector(in)
	if err != nil {
		return nil, err
	}
	v := make([]float64, len(vector))
	for i, e := range vector {
		v[i] = float64(e)
	}
	return s.Search(v, size, epsilon)
}

func Insert(vector []float64, id []byte) error {
	return s.Insert(vector, id)
}

func (s *Service) Insert(vector []float64, id []byte) error {
	in, err := gongt.StrictInsert(vector)
	if err != nil {
		return err
	}

	return s.db.Set(id, in)
}

func Remove(id []byte) error {
	return s.Remove(id)
}

func (s *Service) Remove(id []byte) error {
	in, err := s.db.GetVal(id)
	if err != nil {
		return err
	}

	if err := gongt.StrictRemove(in); err != nil {
		return err
	}

	return s.db.Delete(id)
}
