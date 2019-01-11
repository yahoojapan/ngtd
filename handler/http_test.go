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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd/model"
)

func TestHTTP(t *testing.T) {
	t.Run("TestIndex", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Errorf("Unexpected error: TestIndex(%v)", err)
		}
		w := httptest.NewRecorder()
		Index(w, r)

		body, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("Unexpected error: TestIndex(%v)", err)
		}

		if "/" != string(body) {
			t.Errorf("TestIndex(): %v, wanted: /", string(body))
		}
	})

	t.Run("TestHTTPSearch", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := []struct {
			vector []float64
			want   model.SearchResult
		}{
			{[]float64{1, 0, 0, 0, 0, 0}, model.SearchResult{ID: "a", Distance: 0}},
			{[]float64{0, 1, 0, 0, 0, 0}, model.SearchResult{ID: "b", Distance: 0}},
			{[]float64{0, 0, 1, 0, 0, 0}, model.SearchResult{ID: "c", Distance: 0}},
			{[]float64{0, 0, 0, 1, 0, 0}, model.SearchResult{ID: "d", Distance: 0}},
			{[]float64{0, 0, 0, 0, 1, 0}, model.SearchResult{ID: "e", Distance: 0}},
			{[]float64{0, 0, 0, 0, 0, 1}, model.SearchResult{ID: "f", Distance: 0}},
		}

		for _, tt := range tests {
			req := model.SearchRequest{Vector: tt.vector, Size: 1, Epsilon: gongt.DefaultEpsilon}
			reqBody, err := json.Marshal(req)
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPSearch(%v)", err)
			}

			r, err := http.NewRequest(http.MethodPost, "/search", bytes.NewReader(reqBody))
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPSearch(%v)", err)
			}
			w := httptest.NewRecorder()
			Search(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("Unexpected error: TestHTTPSearch(%v)", err)
			}

			var res model.SearchResponse
			if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
				t.Errorf("Unexpected error: TestHTTPSearch(%v)", err)
			}

			if len(res.Result) != 1 || !reflect.DeepEqual(tt.want.ID, res.Result[0].ID) || tt.want.Distance != res.Result[0].Distance {
				t.Errorf("TestHTTPSearch(%v): %v, wanted: %v", tt.vector, res.Result[0], tt.want)
			}
		}
	})

	t.Run("TestSearchByID", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := []struct {
			id   string
			want model.SearchResult
		}{
			{"a", model.SearchResult{ID: "a", Distance: 0}},
			{"b", model.SearchResult{ID: "b", Distance: 0}},
			{"c", model.SearchResult{ID: "c", Distance: 0}},
			{"d", model.SearchResult{ID: "d", Distance: 0}},
			{"e", model.SearchResult{ID: "e", Distance: 0}},
			{"f", model.SearchResult{ID: "f", Distance: 0}},
		}

		for _, tt := range tests {
			req := model.SearchRequest{ID: tt.id, Size: 1, Epsilon: gongt.DefaultEpsilon}
			reqBody, err := json.Marshal(req)
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPSearchByID(%v)", err)
			}

			r, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(reqBody))
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPSearchByID(%v)", err)
			}
			w := httptest.NewRecorder()
			SearchByID(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("Unexpected error: TestHTTPSearchByID(%v)", err)
			}

			var res model.SearchResponse
			if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
				t.Errorf("Unexpected error: TestHTTPSearchByID(%v)", err)
			}

			if len(res.Result) != 1 || !reflect.DeepEqual(tt.want.ID, res.Result[0].ID) || tt.want.Distance != res.Result[0].Distance {
				t.Errorf("TestHTTPSearchByID(%v): %v, wanted: %v", tt.id, res.Result[0], tt.want)
			}
		}
	})

	t.Run("TestInsert", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := []struct {
			req model.InsertRequest
		}{
			{model.InsertRequest{ID: "g", Vector: []float64{1, 1, 0, 0, 0, 0}}},
			{model.InsertRequest{ID: "h", Vector: []float64{1, 0, 0, 0, 0, 0}}},
			{model.InsertRequest{ID: "i", Vector: []float64{1, 0, 1, 0, 0, 0}}},
			{model.InsertRequest{ID: "j", Vector: []float64{1, 0, 0, 1, 0, 0}}},
			{model.InsertRequest{ID: "k", Vector: []float64{1, 0, 0, 0, 1, 0}}},
			{model.InsertRequest{ID: "l", Vector: []float64{1, 0, 0, 0, 0, 1}}},
		}

		for _, tt := range tests {
			reqBody, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPInsert(%v)", err)
			}

			r, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(reqBody))
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPInsert(%v)", err)
			}
			w := httptest.NewRecorder()
			Insert(w, r)

			var res model.InsertResponse
			if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
				t.Errorf("Unexpected error: TestHTTPInsert(%v)", err)
			}

			if "Success" != res.Status {
				t.Errorf("TestHTTPInsert(%v): %v, wanted: Success", tt.req, res)
			}
		}
	})

	t.Run("TestMultiInsert", func(t *testing.T) {
		defer SetupWithTeardown(t)()

		tests := []struct {
			req model.MultiInsertRequest
		}{
			{
				model.MultiInsertRequest{
					InsertRequests: []model.InsertRequest{
						model.InsertRequest{ID: "g", Vector: []float64{1, 1, 0, 0, 0, 0}},
						model.InsertRequest{ID: "h", Vector: []float64{1, 0, 0, 0, 0, 0}},
						model.InsertRequest{ID: "i", Vector: []float64{1, 0, 1, 0, 0, 0}},
						model.InsertRequest{ID: "j", Vector: []float64{1, 0, 0, 1, 0, 0}},
						model.InsertRequest{ID: "k", Vector: []float64{1, 0, 0, 0, 1, 0}},
						model.InsertRequest{ID: "l", Vector: []float64{1, 0, 0, 0, 0, 1}},
					},
				},
			},
		}

		for _, tt := range tests {
			reqBody, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPMultiInsert(%v)", err)
			}

			r, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(reqBody))
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPInsert(%v)", err)
			}
			w := httptest.NewRecorder()
			MultiInsert(w, r)

			var res model.MultiInsertResponse
			if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
				t.Errorf("Unexpected error: TestHTTPMultiInsert(%v)", err)
			}

			if "Success" != res.Status {
				t.Errorf("TestHTTPMultiInsert(%v): %v, wanted: Success", tt.req, res)
			}
		}
	})

	t.Run("TestRemove", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := []struct {
			id []byte
		}{
			{[]byte("a")},
			{[]byte("b")},
			{[]byte("c")},
			{[]byte("d")},
			{[]byte("e")},
			{[]byte("f")},
		}

		m := mux.NewRouter()
		m.HandleFunc("/{id}", Remove)
		for _, tt := range tests {
			r, err := http.NewRequest(http.MethodGet, "/"+string(tt.id), bytes.NewReader(nil))
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPRemove(%v)", err)
			}
			r.Header.Add("Content-Type", "application/json; charset=utf-8")
			w := httptest.NewRecorder()
			m.ServeHTTP(w, r)

			var res model.RemoveResponse
			if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
				t.Errorf("Unexpected error: TestHTTPRemove(%v)", err)
			}

			if "Success" != res.Status {
				t.Errorf("TestHTTPRemove(%v): %v, wanted: Success", tt.id, res)
			}
		}
	})

	t.Run("TestMultiRemove", func(t *testing.T) {
		defer SetupWithTeardown(t)()

		tests := []struct {
			ids []string
		}{
			{
				[]string{
					"a", "b", "c", "d", "e", "f",
				},
			},
		}

		for _, tt := range tests {
			req := model.MultiRemoveRequest{IDs: tt.ids}
			reqBody, err := json.Marshal(req)
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPMultiRemove(%v)", err)
			}
			r, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(reqBody))
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPMultiRemove(%v)", err)
			}
			r.Header.Add("Content-Type", "application/json; charset=utf-8")
			w := httptest.NewRecorder()
			MultiRemove(w, r)

			var res model.MultiRemoveResponse
			if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
				t.Errorf("Unexpected error: TestHTTPRemove(%v)", err)
			}

			if "Success" != res.Status {
				t.Errorf("TestHTTPRemove(%v): %v, wanted: Success", tt.ids, res)
			}
		}
	})

	t.Run("TestGetObjects", func(t *testing.T) {
		defer SetupWithTeardown(t)()
		tests := []struct {
			ids  []string
			want []model.GetObjectResult
		}{
			{[]string{"a"}, []model.GetObjectResult{
				model.GetObjectResult{ID: "a", Vector: []float32{1, 0, 0, 0, 0, 0}}}},

			{[]string{"a", "b"}, []model.GetObjectResult{
				model.GetObjectResult{ID: "a", Vector: []float32{1, 0, 0, 0, 0, 0}},
				model.GetObjectResult{ID: "b", Vector: []float32{0, 1, 0, 0, 0, 0}}}},

			{[]string{"c", "d", "e", "f"}, []model.GetObjectResult{
				model.GetObjectResult{ID: "c", Vector: []float32{0, 0, 1, 0, 0, 0}},
				model.GetObjectResult{ID: "d", Vector: []float32{0, 0, 0, 1, 0, 0}},
				model.GetObjectResult{ID: "e", Vector: []float32{0, 0, 0, 0, 1, 0}},
				model.GetObjectResult{ID: "f", Vector: []float32{0, 0, 0, 0, 0, 1}}}},

			{[]string{"a", "g", "c"}, []model.GetObjectResult{
				model.GetObjectResult{ID: "a", Vector: []float32{1, 0, 0, 0, 0, 0}},
				model.GetObjectResult{ID: "c", Vector: []float32{0, 0, 1, 0, 0, 0}}}},
		}

		for _, tt := range tests {
			req := model.GetObjectsRequest{IDs: tt.ids}
			reqBody, err := json.Marshal(req)
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPGetObjects(%v)", err)
			}

			r, err := http.NewRequest(http.MethodPost, "/getobjects", bytes.NewReader(reqBody))
			if err != nil {
				t.Errorf("Unexpected error: TestHTTPGetObjects(%v)", err)
			}
			w := httptest.NewRecorder()
			GetObjects(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("Unexpected error: TestHTTPGetObjects(%v)", err)
			}
			var res model.GetObjectsResponse
			if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
				t.Errorf("Unexpected error: TestHTTPGetObjects(%v)", err)
			}

			if !reflect.DeepEqual(tt.want, res.Result) {
				t.Errorf("TestHTTPGetObjects(%v): %v, wanted: %v", tt.ids, res.Result, tt.want)
			}
		}
	})
}
