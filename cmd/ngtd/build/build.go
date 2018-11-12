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

package build

import (
	"fmt"
	"io"
	"sync"

	"github.com/kpango/glg"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd/kvs"
	"github.com/yahoojapan/ngtd/service"
)

type data struct {
	id     []byte
	vector []float64
}

type builder struct {
	db                kvs.KVS
	r                 Reader
	p                 Parser
	parallelParseSize int
	rCh               chan []byte
	wCh               chan data
	wg                *sync.WaitGroup
}

func NewBuilder(db kvs.KVS, r Reader, p Parser, parallelParseSize int) *builder {
	service.SetDB(db)
	return &builder{
		r:                 r,
		p:                 p,
		parallelParseSize: parallelParseSize,
		rCh:               make(chan []byte),
		wCh:               make(chan data, 1),
		wg:                new(sync.WaitGroup),
	}
}

func (b *builder) Run(index string, dimension, poolSize int) error {
	gongt.SetIndexPath(index)
	if dimension > 0 {
		gongt.SetDimension(dimension)
	}
	gongt.Open()

	if errs := gongt.GetErrors(); len(errs) > 0 {
		return fmt.Errorf("Get gongt errors: %v", errs)
	}

	b.build()

	gongt.CreateAndSaveIndex(poolSize)

	return nil
}

func (b *builder) build() {
	b.wg.Add(1)
	go b.read()

	b.wg.Add(1)
	go b.parallelParse()

	b.wg.Add(1)
	go b.write()

	b.wg.Wait()
}

func (b *builder) read() {
	defer close(b.rCh)
	defer b.wg.Done()
	for {
		row, err := b.r.Next()
		if err == io.EOF {
			return
		} else if err != nil {
			glg.Warn(err)
		} else {
			b.rCh <- row
		}
	}
}

func (b *builder) parse(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case buf, ok := <-b.rCh:
			if !ok {
				return
			}
			id, vector, err := b.p.Parse(buf)
			if err != nil {
				glg.Warn(err)
			}
			b.wCh <- data{id, vector}
		default:
		}
	}
}

func (b *builder) parallelParse() {
	defer close(b.wCh)
	defer b.wg.Done()
	w := &sync.WaitGroup{}
	for i := 0; i < b.parallelParseSize; i++ {
		w.Add(1)
		go b.parse(w)
	}
	w.Wait()
}

func (b *builder) write() {
	defer b.wg.Done()
	for {
		if d, ok := <-b.wCh; ok {
			service.Insert(d.vector, d.id)
		} else {
			return
		}
	}
}
