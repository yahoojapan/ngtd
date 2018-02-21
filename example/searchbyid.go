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

package main

import (
	"context"
	"io"
	"sync"

	"github.com/kpango/glg"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd/cmd/ngtd/build"
	pb "github.com/yahoojapan/ngtd/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8200", grpc.WithInsecure())
	if err != nil {
		glg.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewNGTDClient(conn)
	st, err := c.StreamSearchByID(context.Background())
	if err != nil {
		glg.Fatalln(err)
	}

	r, err := build.NewTextReader("assets/random/input.tsv")
	if err != nil {
		glg.Fatalln(err)
	}
	defer r.Close()
	p, err := build.NewTextParser("\t", " ")
	if err != nil {
		glg.Fatalln(err)
	}

	glg.Info("Search")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			line, err := r.Next()
			if err == io.EOF {
				st.CloseSend()
				break
			} else if err != nil {
				glg.Warn(err)
				continue
			}
			id, _, err := p.Parse(line)
			if err != nil {
				glg.Warn(err)
			}

			if err := st.Send(&pb.SearchRequest{
				Id:      id,
				Size_:   10,
				Epsilon: gongt.DefaultEpsilon,
			}); err != nil {
				glg.Warn(err)
			}
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			res, err := st.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				glg.Warn(err)
			} else {
				if res.Error != "" {
					glg.Warn(res.Error)
				} else {
					glg.Info(res.Result)
				}
			}
		}
	}(&wg)

	wg.Wait()
}
