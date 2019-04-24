NAME := ngtd
VERSION := v0.0.1
GO_VERSION:=$(shell go version)
REVISION := $(shell git rev-parse --short HEAD)

.PHONY: build clean

clean:
	go clean ./...
	go clean -modcache
	rm -rf ./*.log
	rm -rf ./*.svg
	rm -rf ./go.mod
	rm -rf ./go.sum
	rm -rf bench
	rm -rf pprof
	rm -rf vendor

init:
	GO111MODULE=on go mod init
	GO111MODULE=on go mod vendor
	sleep 3

deps:
	curl -LO https://github.com/yahoojapan/NGT/archive/v$(NGT_VERSION).tar.gz
	tar zxf v$(NGT_VERSION).tar.gz -C /tmp
	cd /tmp/NGT-$(NGT_VERSION); cmake .
	make -j -C /tmp/NGT-$(NGT_VERSION)
	make install -C /tmp/NGT-$(NGT_VERSION)

proto/ngtd.pb.go: proto/ngtd.proto
	protoc --gofast_out=plugins=grpc:. proto/ngtd.proto

build: proto/ngtd.pb.go
	GO111MODULE=on go build -ldflags="-w -s"

test: clean init proto/ngtd.pb.go
	GO111MODULE=on go test --race -v ./...

docker-build:
	docker build --pull=true --file=Dockerfile -t yahoojapan/ngtd:latest .

docker-push:
	docker push yahoojapan/ngtd:latest

coverage:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out
