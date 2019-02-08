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

// Package ngtd daemonize NGT
package ngtd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kpango/glg"
	"github.com/yahoojapan/gongt"
	"github.com/yahoojapan/ngtd/handler"
	"github.com/yahoojapan/ngtd/kvs"
	pb "github.com/yahoojapan/ngtd/proto"
	"github.com/yahoojapan/ngtd/router"
	"github.com/yahoojapan/ngtd/service"

	"google.golang.org/grpc"
)

// NGTD is base struct
type NGTD struct {
	sigCh   chan os.Signal
	l       net.Listener
	port    string
	running bool
}

type ServerType int

const (
	HTTP ServerType = 1
	GRPC ServerType = 2
)

var (
	ErrServerAlreadyRunning = errors.New("NGTD is already running")
)

// NewNGTD create NGTD struct
func NewNGTD(index string, db kvs.KVS, port int) (*NGTD, error) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	gongt.SetIndexPath(index).Open()
	service.SetDB(db)

	if errs := gongt.GetErrors(); len(errs) > 0 {
		glg.Fatalln(errs)
		return nil, fmt.Errorf("%v", errs)
	}

	p := strconv.Itoa(port)

	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		return nil, err
	}

	return &NGTD{
		sigCh: sigCh,
		l:     l,
		port:  p,
	}, nil
}

func (n *NGTD) ListenAndServe(t ServerType) error {
	switch t {
	case HTTP:
		return n.listenAndServeHTTP()
	case GRPC:
		return n.listenAndServeGRPC()
	}
	return nil
}

func (n *NGTD) listenAndServeHTTP() error {
	if n.running {
		return ErrServerAlreadyRunning
	}
	defer gongt.Close()
	srv := &http.Server{
		Addr:    ":" + n.port,
		Handler: router.NewRouter(),
	}

	go func() {
		n.running = true
		glg.Info("NGTD HTTP API Server Starting ...")
		if err := srv.Serve(n.l); err != nil {
			glg.Error(err)
			n.sigCh <- syscall.SIGINT
		}
	}()
	// wait terminate signal
	<-n.sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	glg.Info("NGTD HTTP API Server Shutdown ...")

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (n *NGTD) listenAndServeGRPC() error {
	if n.running {
		return ErrServerAlreadyRunning
	}

	defer gongt.Close()
	srv := grpc.NewServer()
	pb.RegisterNGTDServer(srv, &handler.GRPC{})

	go func() {
		n.running = true
		glg.Info("NGTD GRPC Server Starting ...")
		if err := srv.Serve(n.l); err != nil {
			glg.Error(err)
			n.sigCh <- syscall.SIGINT
		}
	}()

	// wait terminate signal
	<-n.sigCh
	glg.Info("NGTD GRPC Server Shutdown ...")
	srv.Stop()

	return nil
}

func (n *NGTD) ListenAndServeProfile(port int) error {
	return http.ListenAndServe(":"+strconv.Itoa(port), router.NewPprofRouter())
}

func (n *NGTD) Stop() {
	n.sigCh <- syscall.SIGINT
}
