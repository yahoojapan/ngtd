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

package model

type SearchRequest struct {
	Vector  []float64 `json:"vector"`
	ID      string    `json:"id"`
	Size    int       `json:"size"`
	Epsilon float32   `json:"epsilon"`
}

type SearchResult struct {
	ID       string  `json:"id"`
	Distance float32 `json:"distance"`
}

type SearchResponse struct {
	Result []SearchResult `json:"result"`
	Errors []error        `json:"errors"`
}

type InsertRequest struct {
	Vector []float64 `json:"vector"`
	ID     string    `json:"id"`
}

type InsertResponse struct {
	Status string `json:"status"`
}

type MultiInsertRequest struct {
	InsertRequests []InsertRequest `json:"insert_requests"`
}

type MultiInsertResponse struct {
	Status string  `json:"status"`
	Errors []error `json:"errors"`
}

type RemoveResponse struct {
	Status string `json:"status"`
}

type MultiRemoveRequest struct {
	IDs []string `json:"ids"`
}

type MultiRemoveResponse struct {
	Status string  `json:"status"`
	Errors []error `json:"errors"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   error  `json:"error"`
	Message string `json:"message"`
}

type DefaultResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Error   error  `json:"error"`
}
