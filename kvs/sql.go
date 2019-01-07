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

package kvs

import (
	"database/sql"
)

type SQL struct {
	db *sql.DB
}

func NewSQL(db *sql.DB) (*SQL, error) {
	_, err := db.Exec("create table if not exists kvs (key varchar(256) not null primary key, val integer not null unique)")
	if err != nil {
		return nil, err
	}
	return &SQL{
		db: db,
	}, nil
}

func (s *SQL) GetKey(val uint) ([]byte, error) {
	stmt, err := s.db.Prepare("select key from kvs where val = ? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var key string
	if err := stmt.QueryRow(val).Scan(&key); err != nil {
		return nil, err
	}

	return []byte(key), nil
}

func (s *SQL) GetKeys(vals []uint) ([][]byte, error) {
	ret := make([][]byte, len(vals))
	for i, val := range vals {
		k, err := s.GetKey(val)
		if err != nil {
			return nil, err
		}
		ret[i] = k
	}
	return ret, nil
}

func (s *SQL) GetVal(key []byte) (uint, error) {
	stmt, err := s.db.Prepare("select val from kvs where key = ? limit 1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val uint
	if err := stmt.QueryRow(string(key)).Scan(&val); err != nil {
		return 0, err
	}

	return val, nil
}

func (s *SQL) Set(key []byte, val uint) error {
	stmt, err := s.db.Prepare("insert into kvs (key, val) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(string(key), val)
	return err
}

func (s *SQL) Delete(key []byte) error {
	stmt, err := s.db.Prepare("delete from kvs where key = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(string(key))
	return err

}

func (s *SQL) Close() error {
	return s.db.Close()
}
