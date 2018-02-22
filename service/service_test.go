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
	"reflect"
	"testing"

	"github.com/yahoojapan/ngtd/ngtdtest"
	"github.com/yahoojapan/gongt"
)

func SetupWithTeardown(t *testing.T) func() {
	Get().SetDB(ngtdtest.CreateDB(t))
	return func() {
		ngtdtest.DeleteDB(t)
	}
}

func TestService(t *testing.T) {
	t.Run("TestSearch", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := []struct{
			vector []float64
			want   SearchResult
		}{
			{[]float64{1, 0, 0, 0, 0, 0}, SearchResult{Id:[]byte("a"), Distance:0}},
			{[]float64{0, 1, 0, 0, 0, 0}, SearchResult{Id:[]byte("b"), Distance:0}},
			{[]float64{0, 0, 1, 0, 0, 0}, SearchResult{Id:[]byte("c"), Distance:0}},
			{[]float64{0, 0, 0, 1, 0, 0}, SearchResult{Id:[]byte("d"), Distance:0}},
			{[]float64{0, 0, 0, 0, 1, 0}, SearchResult{Id:[]byte("e"), Distance:0}},
			{[]float64{0, 0, 0, 0, 0, 1}, SearchResult{Id:[]byte("f"), Distance:0}},
		}
		
		for _, tt := range tests {
			res, err := Search(tt.vector, 1, gongt.DefaultEpsilon)
			if err != nil {
				t.Errorf("Unexpected error: TestSearch(%v)", err)
			}
			if !reflect.DeepEqual(res[0].Id, tt.want.Id) || res[0].Distance != tt.want.Distance {
				t.Errorf("TestSearch(%v): %v, wanted: %v", tt.vector, res[0], tt.want)
			}
		}
	})

	t.Run("TestSearchByID", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := []struct {
			id   []byte
			want SearchResult
		}{
			{[]byte("a"), SearchResult{Id:[]byte("a"), Distance:0}},
			{[]byte("b"), SearchResult{Id:[]byte("b"), Distance:0}},
			{[]byte("c"), SearchResult{Id:[]byte("c"), Distance:0}},
			{[]byte("d"), SearchResult{Id:[]byte("d"), Distance:0}},
			{[]byte("e"), SearchResult{Id:[]byte("e"), Distance:0}},
			{[]byte("f"), SearchResult{Id:[]byte("f"), Distance:0}},
		}
		
		for _, tt := range tests {
			res, err := SearchByID(tt.id, 1, gongt.DefaultEpsilon)
			if err != nil {
				t.Errorf("Unexpected error: TestSearchByID(%v)", err)
			}
			if !reflect.DeepEqual(res[0].Id, tt.want.Id) || res[0].Distance != tt.want.Distance {
				t.Errorf("TestSearchByID(%v): %v, wanted: %v", tt.id, res[0], tt.want)
			}
		}
	})

  t.Run("TestInsert", func(t *testing.T) {
		gongt.Get().SetObjectType(gongt.Uint8).SetDimension(6)
		defer SetupWithTeardown(t)()
		tests := []struct {
			id     []byte
			vector []float64
		}{
			{id: []byte("g"), vector: []float64{1, 0, 0, 0, 0, 0}},
		}
		
		for _, tt := range tests {
			if err := Insert(tt.vector, tt.id); err != nil {
				t.Errorf("Unexpected error: TestInsert(%v)", err)
			}
		}
	})

  t.Run("TestRemove", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := [][]byte{
			[]byte("a"),
			[]byte("b"),
			[]byte("c"),
			[]byte("d"),
			[]byte("e"),
			[]byte("f"),
		}
		
		for _, id := range tests {
			if err := Remove(id); err != nil {
				t.Errorf("Unexpected error: TestRemove(%v)", err)
			}
		}
	})
}
