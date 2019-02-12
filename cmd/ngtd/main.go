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
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/kpango/glg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd"
	"github.com/yahoojapan/ngtd/cmd/ngtd/build"
	"github.com/yahoojapan/ngtd/kvs"
	"golang.org/x/sync/errgroup"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	// Version is ngtd version
	Version = "1.1.0"
	// Revision is release revision
	Revision = "profilable"

	index     string
	dbType    string
	dimension int
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(runtime.Error); ok {
				glg.Fatalln(err)
			}
			glg.Info(err.(error))
		}
	}()
	glg.Get().SetMode(glg.STD)

	app := cli.NewApp()
	app.Name = "ngtd"
	app.Usage = "NGT Daemonize"
	app.Version = Version + "-" + Revision

	flags := func(f []cli.Flag) []cli.Flag {
		commonFlags := []cli.Flag{
			cli.StringFlag{
				Name:        "index, i",
				Value:       "/usr/share/ngtd/index",
				Usage:       "path to index",
				Destination: &index,
			},
			cli.IntFlag{
				Name:        "dimension, d",
				Value:       -1,
				Usage:       "vector dimension size.(Must set if create new index)",
				Destination: &dimension,
			},

			cli.StringFlag{
				Name:        "database-type, t",
				Value:       "",
				Usage:       "ngtd inner kvs type(redis, golevel, bolt, sqlite or inmem)",
				Destination: &dbType,
			},
			cli.StringFlag{
				Name:  "database-path, p",
				Value: "/usr/share/ngtd/db/kvs.db",
				Usage: "ngtd inner kvs path(for golevel, bolt and sqlite)",
			},
			cli.StringFlag{
				Name:  "redis-host",
				Value: "localhost",
				Usage: "redis running host",
			},
			cli.StringFlag{
				Name:  "redis-port",
				Value: "6379",
				Usage: "redis running port",
			},
			cli.StringFlag{
				Name:  "redis-password",
				Value: "",
				Usage: "redis password",
			},
			cli.IntSliceFlag{
				Name:  "redis-database-index, I",
				Usage: "list up 2 redis database indexes",
			},
			cli.IntFlag{
				Name:  "redis-ping-timeout",
				Value: 600,
				Usage: "wait until redis finish loading dump.rdb or timeout. (unit=second)",
			},
			cli.IntFlag{
				Name:  "redis-ping-retry-freq",
				Value: 10,
				Usage: "try ping this frequency. (unit=second)",
			},
		}

		return append(commonFlags, f...)
	}

	database := func(c *cli.Context) (kvs.KVS, error) {
		p := c.String("database-path")
		switch dbType {
		case "redis":
			indexes := c.IntSlice("redis-database-index")
			if len(indexes) == 0 {
				indexes = cli.IntSlice{0, 1}
			}
			pingTimeout := c.Int("redis-ping-timeout")
			if pingTimeout <= 0 {
				return nil, fmt.Errorf("invalid value: redis-ping-timeout should be greater than 0: %d", pingTimeout)
			}
			pingRetryFreq := c.Int("redis-ping-retry-freq")
			if pingRetryFreq <= 0 {
				return nil, fmt.Errorf("invalid value: redis-ping-retry-freq should be greater than 0: %d", pingRetryFreq)
			}
			return kvs.NewRedis(c.String("redis-host"), c.String("redis-port"), c.String("redis-password"), indexes[0], indexes[1], time.Duration(pingTimeout)*time.Second, time.Duration(pingRetryFreq)*time.Second)
		case "bolt":
			return kvs.NewBoltDB(p)
		case "golevel":
			return kvs.NewGoLevel(p)
		case "sqlite":
			s, err := sql.Open("sqlite3", p)
			if err != nil {
				return nil, err
			}
			return kvs.NewSQL(s)
		case "inmem":
			return kvs.NewMemory(), nil
		default:
			return nil, fmt.Errorf("unsupported database type: %v", dbType)
		}
	}

	serve := func(name string, alias []string, t ngtd.ServerType) cli.Command {
		return cli.Command{
			Name:    name,
			Aliases: alias,
			Usage:   "serve ngtd index by " + name,
			Flags: flags([]cli.Flag{
				cli.IntFlag{
					Name:  "port, P",
					Value: 8200,
					Usage: "listening port",
				},
				cli.IntFlag{
					Name:  "pprof-port, PP",
					Value: 6060,
					Usage: "listening pprof port",
				},
				cli.BoolFlag{
					Name:  "pprof pp",
					Usage: "enable pprof server",
				},
			}),
			Action: func(c *cli.Context) error {
				if dimension > 0 {
					gongt.SetDimension(dimension)
				}
				db, err := database(c)
				if err != nil {
					return err
				}
				n, err := ngtd.NewNGTD(index, db, c.Int("port"))
				if err != nil {
					return err
				}
				eg := new(errgroup.Group)
				if c.Bool("pprof") {
					eg.Go(func() error {
						return n.ListenAndServeProfile(c.Int("pprof-port"))
					})
				}
				eg.Go(func() error {
					return n.ListenAndServe(t)
				})

				return eg.Wait()
			},
		}
	}

	app.Commands = []cli.Command{
		serve("http", []string{"H"}, ngtd.HTTP),
		serve("grpc", []string{"g"}, ngtd.GRPC),
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "build ngtd index",
			Flags: flags([]cli.Flag{
				cli.StringSliceFlag{
					Name:  "text-delimiter, D",
					Value: &cli.StringSlice{"\t", " "},
					Usage: "delimiter for text input",
				},
				cli.IntFlag{
					Name:  "pool",
					Value: runtime.NumCPU(),
					Usage: "number of CPU using NGT indexing",
				},
				cli.IntFlag{
					Name:  "parallel-parse",
					Value: runtime.NumCPU(),
					Usage: "number of CPU using input parser",
				},
			}),
			Action: func(c *cli.Context) error {
				db, err := database(c)
				if err != nil {
					return err
				}
				in := c.Args().Get(0)
				r, err := build.NewTextReader(in)
				if err != nil {
					return err
				}
				d := c.StringSlice("text-delimiter")
				p, err := build.NewTextParser(d[0], d[1])
				if err != nil {
					return err
				}
				return build.NewBuilder(db, r, p, c.Int("parallel-parse")).Run(index, dimension, c.Int("pool"))
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		glg.Fatal(err)
	}
}
