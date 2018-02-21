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
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type TextReader struct {
	f *os.File
	r *bufio.Reader
	v bool
}

func NewTextReader(p string) (*TextReader, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	return &TextReader{
		f: f,
		r: r,
	}, nil
}

func (x TextReader) Next() ([]byte, error) {
	line, _, err := x.r.ReadLine()
	return line, err
}

func (x TextReader) Close() error {
	return x.f.Close()
}

type TextParser struct {
	kvDelimiter string
	vDelimiter  string
}

func NewTextParser(kvd, vd string) (*TextParser, error) {
	return &TextParser{kvDelimiter: kvd, vDelimiter: vd}, nil
}

const (
	hexPrefix = "0x"
	base      = 16
	bitSize   = 64
)

func (p TextParser) Parse(in []byte) ([]byte, []float64, error) {
	kv := strings.SplitN(string(in), p.kvDelimiter, 2)
	if len(kv) < 2 {
		return nil, nil, fmt.Errorf("cannot split by %v", p.kvDelimiter)
	}
	id := *(*[]byte)(unsafe.Pointer(&kv[0]))
	v := strings.Split(kv[1], p.vDelimiter)
	vector := make([]float64, len(v))
	for i, e := range v {
		if strings.HasPrefix(e, hexPrefix) {
			bits, err := strconv.ParseUint(strings.TrimPrefix(e, hexPrefix), base, bitSize)
			if err != nil {
				return nil, nil, err
			}
			vector[i] = math.Float64frombits(bits)
		} else {
			v, err := strconv.ParseFloat(e, bitSize)
			if err != nil {
				return nil, nil, err
			}
			vector[i] = v
		}
	}
	return id, vector, nil
}
