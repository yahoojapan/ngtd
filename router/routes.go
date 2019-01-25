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

package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/yahoojapan/ngtd/handler"
)

//Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var (
	routes = []Route{
		Route{
			"Index",
			http.MethodGet,
			"/",
			handler.Index,
		},
		Route{
			"Search",
			http.MethodPost,
			"/search",
			handler.Search,
		},
		Route{
			"SearchByID",
			http.MethodPost,
			"/searchbyid",
			handler.SearchByID,
		},
		Route{
			"Insert",
			http.MethodPost,
			"/insert",
			handler.Insert,
		},
		Route{
			"MultiInsert",
			http.MethodPost,
			"/multiinsert",
			handler.MultiInsert,
		},
		Route{
			"Remove",
			http.MethodGet,
			"/remove/{id}",
			handler.Remove,
		},
		Route{
			"MultiRemove",
			http.MethodPost,
			"/multiremove",
			handler.MultiRemove,
		},
		Route{
			"CreateIndex",
			http.MethodGet,
			"/index/create/{pool_size}",
			handler.CreateIndex,
		},
		Route{
			"SaveIndex",
			http.MethodGet,
			"/index/save",
			handler.SaveIndex,
		},
		Route{
			"GetErrors",
			http.MethodGet,
			"/errors",
			handler.GetErrors,
		},
		Route{
			"GetDim",
			http.MethodGet,
			"/dimension",
			handler.GetDimension,
		},
		Route{
			"GetObjects",
			http.MethodPost,
			"/getobjects",
			handler.GetObjects,
		},
	}

	profiles = []Route{
		Route{
			"Debug pprof",
			http.MethodGet,
			"/debug/pprof/",
			pprof.Index,
		},
		Route{
			"Debug cmdline",
			http.MethodGet,
			"/debug/pprof/cmdline",
			pprof.Cmdline,
		},
		Route{
			"Debug profile",
			http.MethodGet,
			"/debug/pprof/profile",
			pprof.Profile,
		},
		Route{
			"Debug symbol profile",
			http.MethodGet,
			"/debug/pprof/symbol",
			pprof.Symbol,
		},
		Route{
			"Debug trace profile",
			http.MethodGet,
			"/debug/pprof/trace",
			pprof.Trace,
		},
		Route{
			"Debug heap profile",
			http.MethodGet,
			"/debug/pprof/heap",
			pprof.Handler("heap").ServeHTTP,
		},
		Route{
			"Debug goroutine profile",
			http.MethodGet,
			"/debug/pprof/goroutine",
			pprof.Handler("goroutine").ServeHTTP,
		},
		Route{
			"Debug thread profile",
			http.MethodGet,
			"/debug/pprof/threadcreate",
			pprof.Handler("threadcreate").ServeHTTP,
		},
		Route{
			"Debug block profile",
			http.MethodGet,
			"/debug/pprof/block",
			pprof.Handler("block").ServeHTTP,
		},
	}
)
