# NGTD: Serving [NGT](https://github.com/yahoojapan/NGT) over HTTP or gRPC. [![License: Apache](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0) [![release](https://img.shields.io/github/release/yahoojapan/ngtd.svg?style=flat-square)](https://github.com/yahoojapan/ngtd/releases/latest) [![CircleCI](https://circleci.com/gh/yahoojapan/ngtd.svg)](https://circleci.com/gh/yahoojapan/ngtd) [![codecov](https://codecov.io/gh/yahoojapan/ngtd/branch/master/graph/badge.svg)](https://codecov.io/gh/yahoojapan/ngtd) [![Go Report Card](https://goreportcard.com/badge/github.com/yahoojapan/ngtd)](https://goreportcard.com/report/github.com/yahoojapan/ngtd) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/b03d543ee4a9448ba6d25f94f4989ba4)](https://www.codacy.com/app/i.can.feel.gravity/ngtd?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=yahoojapan/ngtd&amp;utm_campaign=Badge_Grade) [![GoDoc](http://godoc.org/github.com/yahoojapan/ngtd?status.svg)](http://godoc.org/github.com/yahoojapan/ngtd) [![DepShield Badge](https://depshield.sonatype.org/badges/yahoojapan/ngtd/depshield.svg)](https://depshield.github.io)

Description
-----------
NGTD provides serving function for [NGT](https://github.com/yahoojapan/NGT).

NGTD supports gRPC and HTTP protocol, so you can implement applications with your favorite programming language.

You can set any labels for each vectors, and enable to search with the label.

Install
-------
You must install [NGT](https://github.com/yahoojapan/NGT) before installing ngtd.

The easiest way to install is the following command:
```
$ go get -u github.com/yahoojapan/ngtd/cmd/ngtd
```

## Docker
You can get [ngtd docker image](https://hub.docker.com/r/yahoojapan/ngtd/) with the following command.

```
$ docker pull yahoojapan/ngtd
```

The image include only NGTD single binary, you enable to run the docker container.

Usage
-----
```
$ ngtd --help
NAME:
   ngtd - NGT Daemonize

USAGE:
   ngtd [global options] command [command options] [arguments...]

VERSION:
   0.0.1-first

COMMANDS:
     http, H   serve ngtd index by http
     grpc, g   serve ngtd index by grpc
     build, b  build ngtd index
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Serve
ngtd supports http and gRPC protocol.
### HTTP
```
$ ngtd http --help
NAME:
   ngtd http - serve ngtd index by http

USAGE:
   ngtd http [command options] [arguments...]

OPTIONS:
   --index value, -i value                 path to index (default: "/usr/share/ngtd/index")
   --dimension value, -d value             vector dimension size.(Must set if create new index) (default: -1)
   --database-type value, -t value         ngtd inner kvs type(redis, golevel, bolt or sqlite)
   --database-path value, -p value         ngtd inner kvs path(for golevel, bolt and sqlite) (default: "/usr/share/ngtd/db/kvs.db")
   --redis-host value                      redis running host (default: "localhost")
   --redis-port value                      redis running port (default: "6379")
   --redis-password value                  redis password
   --redis-database-index value, -I value  list up 2 redis database indexes (default: 0, 1)
   --port value, -P value                  listening port (default: 8200)
```

#### Request
```
$ curl -H 'Content-Type: application/json' -X POST http://localhost:8200/search -d '{"vector":[...], "size": 10, "epsilon": 0.01}'
$ curl -H 'Content-Type: application/json' -X POST http://localhost:8200/searchbyid -d '{"id":"<id>", "size": 10, "epsilon": 0.01}'
```
If you want more information, please read [model.go](model/model.go)

### gRPC
```
$ ngtd grpc --help
NAME:
   ngtd grpc - serve ngtd index by grpc

USAGE:
   ngtd grpc [command options] [arguments...]

OPTIONS:
   --index value, -i value                 path to index (default: "/usr/share/ngtd/index")
   --dimension value, -d value             vector dimension size.(Must set if create new index) (default: -1)
   --database-type value, -t value         ngtd inner kvs type(redis, golevel, bolt or sqlite)
   --database-path value, -p value         ngtd inner kvs path(for golevel, bolt and sqlite) (default: "/usr/share/ngtd/db/kvs.db")
   --redis-host value                      redis running host (default: "localhost")
   --redis-port value                      redis running port (default: "6379")
   --redis-password value                  redis password
   --redis-database-index value, -I value  list up 2 redis database indexes (default: 0, 1)
   --port value, -P value                  listening port (default: 8200)
```

#### Client
If you use language except golang, compile [proto file](proto/ngtd.proto) for the language.

Go examples are in [example/](example/).

## Build
Build will construct the NGT database by importing the Vector data, please check the Input Format section at the bottom of the Build chapter for the file format to import.
```
$ ngtd build --help
NAME:
   ngtd build - build ngtd index

USAGE:
   ngtd build [command options] [arguments...]

OPTIONS:
   --index value, -i value                 path to index (default: "/usr/share/ngtd/index")
   --dimension value, -d value             vector dimension size.(Must set if create new index) (default: -1)
   --database-type value, -t value         ngtd inner kvs type(redis, golevel, bolt or sqlite)
   --database-path value, -p value         ngtd inner kvs path(for golevel, bolt and sqlite) (default: "/usr/share/ngtd/db/kvs.db")
   --redis-host value                      redis running host (default: "localhost")
   --redis-port value                      redis running port (default: "6379")
   --redis-password value                  redis password
   --redis-database-index value, -I value  list up 2 redis database indexes (default: 0, 1)
   --text-delimiter value, -D value        delimiter for text input (default: "\t", " ")
   --pool value                            number of CPU using NGT indexing (default: 8)
   --parallel-parse value                  number of CPU using input parser (default: 8)
```

### Input format
Now, we support only text format.

It enables for id to use any character without delimiter1 and for vector elements to use decimal/hex format.

If you choice hex format, you must begin the value with "0x".
```
<id1><delimiter1><v11><delimiter2><v12><delimiter2>...<delimiter2><v1d>\n
<id2><delimiter1><v21><delimiter2><v22><delimiter2>...<delimiter2><v2d>\n
...
<id1><delimiter1><vn1><delimiter2><vn2><delimiter2>...<delimiter2><vnd>\n
```

License
-------

Copyright (C) 2018 Yahoo Japan Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this software except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Contributor License Agreement
-----------------------------

This project requires contributors to agree to a [Contributor License Agreement (CLA)](https://gist.github.com/ydnjp/3095832f100d5c3d2592).

Note that only for contributions to the ngtd repository on the GitHub (https://github.com/yahoojapan/ngtd), the contributors of them shall be deemed to have agreed to the CLA without individual written agreements.

Authors
-------

[Kosuke Morimoto](https://github.com/kou-m)  
[kpango](https://github.com/kpango)
