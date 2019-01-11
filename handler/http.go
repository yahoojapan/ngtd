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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"unsafe"

	"github.com/gorilla/mux"
	"github.com/kpango/glg"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd/model"
	"github.com/yahoojapan/ngtd/service"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.URL.String())
}

func Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqBody model.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w,
			http.StatusBadRequest,
			"Invalid JSON Format",
			err)
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	result, err := service.Search(reqBody.Vector, reqBody.Size, reqBody.Epsilon)
	if err != nil {
		ErrorResponse(w,
			http.StatusInternalServerError,
			"Search Error",
			err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SearchResponse{
		Result: toModelSearchResult(result),
	})
}

func SearchByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqBody model.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w,
			http.StatusBadRequest,
			"Invalid JSON Format",
			err)
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	result, err := service.SearchByID(*(*[]byte)(unsafe.Pointer(&reqBody.ID)), reqBody.Size, reqBody.Epsilon)
	if err != nil {
		ErrorResponse(w,
			http.StatusInternalServerError,
			"Search Error",
			err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SearchResponse{
		Result: toModelSearchResult(result),
	})
}

func Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqBody model.InsertRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w,
			http.StatusBadRequest,
			"Invalid JSON Format",
			err)
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	err = service.Insert(reqBody.Vector, *(*[]byte)(unsafe.Pointer(&reqBody.ID)))
	if err != nil {
		ErrorResponse(w,
			http.StatusInternalServerError,
			"Insert Failed",
			err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.InsertResponse{
		Status: "Success",
	})
}

func MultiInsert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqBody model.MultiInsertRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w,
			http.StatusBadRequest,
			"Invalid JSON Format",
			err)
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	errs := make([]error, 0, len(reqBody.InsertRequests))
	for _, insertRequest := range reqBody.InsertRequests {
		err := service.Insert(insertRequest.Vector, *(*[]byte)(unsafe.Pointer(&insertRequest.ID)))
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.MultiInsertResponse{
			Status: "Failed",
			Errors: errs,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.MultiInsertResponse{
		Status: "Success",
		Errors: nil,
	})
}

func Remove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id := mux.Vars(r)["id"]
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	err := service.Remove([]byte(id))
	if err != nil {
		ErrorResponse(w,
			http.StatusInternalServerError,
			"Remove Failed",
			err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.RemoveResponse{
		Status: "Success",
	})
}

func MultiRemove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqBody model.MultiRemoveRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w,
			http.StatusBadRequest,
			"Invalid JSON Format",
			err)
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	errs := make([]error, 0, len(reqBody.IDs))
	for _, id := range reqBody.IDs {
		err := service.Remove(*(*[]byte)(unsafe.Pointer(&id)))
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.MultiRemoveResponse{
			Status: "Failed",
			Errors: errs,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.MultiRemoveResponse{
		Status: "Success",
		Errors: nil,
	})
}

func CreateIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	poolSize, err := strconv.Atoi(mux.Vars(r)["pool_size"])
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	if err != nil {
		ErrorResponse(w,
			http.StatusBadRequest,
			"Bad Request",
			err)
		return
	}

	err = gongt.CreateIndex(poolSize)
	if err != nil {
		ErrorResponse(w,
			http.StatusInternalServerError,
			"CreateIndex Failed",
			err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.DefaultResponse{
		Code:    http.StatusOK,
		Message: "Index Successfully Created",
	})
}

func SaveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	err := gongt.SaveIndex()
	if err != nil {
		ErrorResponse(w,
			http.StatusInternalServerError,
			"SaveIndex Failed",
			err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.DefaultResponse{
		Code:    http.StatusOK,
		Message: "Index Successfully Saved",
	})
}

func GetErrors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	json.NewEncoder(w).Encode(struct {
		Errors []error `json:"errors"`
	}{
		Errors: gongt.GetErrors(),
	})
}

func GetDimension(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	json.NewEncoder(w).Encode(struct {
		Dimension int `json:"dimension"`
	}{
		Dimension: gongt.GetDim(),
	})
}

// GetObjects returns vectors.
func GetObjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqBody model.GetObjectsRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		ErrorResponse(w,
			http.StatusBadRequest,
			"Invalid JSON Format",
			err)
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()

	results := make([]model.GetObjectResult, 0, len(reqBody.IDs))
	errs := make([]string, 0, len(reqBody.IDs))
	for _, id := range reqBody.IDs {
		result, err := service.GetObject(*(*[]byte)(unsafe.Pointer(&id)))
		if err != nil {
			errs = append(errs, fmt.Sprintf("Error: GetObject(%s) caused %s", id, err.Error()))
		} else {
			results = append(results, model.GetObjectResult{
				ID:     *(*string)(unsafe.Pointer(&result.Id)),
				Vector: result.Vector,
			})
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.GetObjectsResponse{
		Result: results,
		Errors: errs,
	})
}

func ErrorResponse(w http.ResponseWriter, code int, message string, err error) {
	glg.Error(err)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(model.DefaultResponse{
		Code:    code,
		Error:   err,
		Message: message,
	})
}

func toModelSearchResult(s []service.SearchResult) []model.SearchResult {
	ret := make([]model.SearchResult, len(s))
	for i, r := range s {
		ret[i] = model.SearchResult{
			ID:       *(*string)(unsafe.Pointer(&r.Id)),
			Distance: r.Distance,
		}
	}
	return ret
}
