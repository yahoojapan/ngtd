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

package ngtdtest

import (
	"os"
	"encoding/csv"
	"fmt"
	"testing"

	"github.com/yahoojapan/gongt"
)

const (
	db = "../assets/test/index"
)

func CreateDB(t *testing.T) *Map {
	t.Helper()
	m := NewMap()
	gongt.Get().SetIndexPath(db).SetDimension(6).Open()

	f, err := os.Open("../assets/test/test.tsv")
	if err != nil {
		t.Errorf("Unexpected error: CreateDBWithDelete(%v)", err)
	}
	r := csv.NewReader(f)
	r.Comma = '\t'
	records, err := r.ReadAll()
	if err != nil {
		t.Errorf("Unexpected error: CreateDBWithDelete(%v)", err)
	}
	for _, record := range records {
		var a, b, c, d, e, f float64
		fmt.Sscanf(record[1], "%f %f %f %f %f %f", &a, &b, &c, &d, &e, &f)
		v := []float64{a, b, c, d, e, f}
		id, err := gongt.StrictInsert(v)
		if err != nil {
			t.Errorf("Unexpected error: CreateDBWithDelete(%v)", err)
		}
		m.Set([]byte(record[0]), id)
	}
	gongt.CreateAndSaveIndex(2)

	return m
}

func DeleteDB(t *testing.T) {
	t.Helper()
	gongt.Close()
	os.RemoveAll(db)
}
