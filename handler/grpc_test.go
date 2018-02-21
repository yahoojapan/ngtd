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

package handler

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"

	pb "github.com/yahoojapan/ngtd/proto"
	"github.com/yahoojapan/gongt"
)

func TestGRPC(t *testing.T) {
	t.Run("TestSearch", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		s := &GRPC{}
		tests := []struct{
			vector []float64
			want   pb.ObjectDistance
		}{
			{[]float64{1, 0, 0, 0, 0, 0}, pb.ObjectDistance{Id:[]byte("a"), Distance:0}},
			{[]float64{0, 1, 0, 0, 0, 0}, pb.ObjectDistance{Id:[]byte("b"), Distance:0}},
			{[]float64{0, 0, 1, 0, 0, 0}, pb.ObjectDistance{Id:[]byte("c"), Distance:0}},
			{[]float64{0, 0, 0, 1, 0, 0}, pb.ObjectDistance{Id:[]byte("d"), Distance:0}},
			{[]float64{0, 0, 0, 0, 1, 0}, pb.ObjectDistance{Id:[]byte("e"), Distance:0}},
			{[]float64{0, 0, 0, 0, 0, 1}, pb.ObjectDistance{Id:[]byte("f"), Distance:0}},
		}
		
		for _, tt := range tests {
			req := &pb.SearchRequest{Vector: tt.vector, Size_: 1, Epsilon: gongt.DefaultEpsilon}
			res, err := s.Search(context.Background(), req)
			if err != nil {
				t.Errorf("Unexpected error: TestSearch(%v)", err)
			}
			if !reflect.DeepEqual(res.Result[0].Id, tt.want.Id) || res.Result[0].Distance != tt.want.Distance {
				t.Errorf("TestSearch(%v): %v, wanted: %v", tt.vector, res.Result[0], tt.want)
			}
		}
	})

	t.Run("TestGRPCSearchByID", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		s := GRPC{}
		tests := []struct {
			ID   []byte
			want pb.ObjectDistance
		}{
			{[]byte("a"), pb.ObjectDistance{Id:[]byte("a"), Distance:0}},
			{[]byte("b"), pb.ObjectDistance{Id:[]byte("b"), Distance:0}},
			{[]byte("c"), pb.ObjectDistance{Id:[]byte("c"), Distance:0}},
			{[]byte("d"), pb.ObjectDistance{Id:[]byte("d"), Distance:0}},
			{[]byte("e"), pb.ObjectDistance{Id:[]byte("e"), Distance:0}},
			{[]byte("f"), pb.ObjectDistance{Id:[]byte("f"), Distance:0}},
		}

		for _, tt := range tests {
			req := &pb.SearchRequest{Id: tt.ID, Size_: 1, Epsilon: gongt.DefaultEpsilon}
			res, err := s.SearchByID(context.Background(), req)
			if err != nil {
				t.Errorf("Unexpected error: TestSearchByID(%v)", err)
			}
			if !reflect.DeepEqual(res.Result[0].Id, tt.want.Id) || res.Result[0].Distance != tt.want.Distance {
				t.Errorf("TestSearchByID(%v): %v, wanted: %v", tt.ID, res.Result[0], tt.want)
			}
		}
	})

	t.Run("TestInsert", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		s := GRPC{}
		tests := []struct {
			req pb.InsertRequest
		}{
			{pb.InsertRequest{Id: []byte("g"), Vector: []float64{1, 0, 0, 0, 0, 0}}},
		}

		for _, tt := range tests {
			res, err := s.Insert(context.Background(), &tt.req)
			if err != nil {
				t.Errorf("Unexpected error: TestInsert(%v)", err)
			}
			if res.Error != "" {
				t.Errorf("TestInsert(%v): %v, wanted: nil", tt.req, res.Error)
			}
		}
	})

	t.Run("TestRemove", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		s := GRPC{}
		tests := [][]byte{
			[]byte("a"),
			[]byte("b"),
			[]byte("c"),
			[]byte("d"),
			[]byte("e"),
			[]byte("f"),
		}

		for _, id := range tests {
			req := &pb.RemoveRequest{Id: id}
			_, err := s.Remove(context.Background(), req)
			if err != nil {
				t.Errorf("Unexpected error: TestRemove(%v)", err)
			}
		}
	})

	t.Run("TestGetDimension", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		g := GRPC{}
		const want = 6

		req := &pb.Empty{}
		res, err := g.GetDimension(context.Background(), req)
		if err != nil {
			t.Errorf("Unexpected error: TestRemove(%v)", err)
		}
		if res.Dimension != want {
			t.Errorf("TestGetDimension(): %v, wanted: %v", res.Dimension, want)
		}
	})
}
