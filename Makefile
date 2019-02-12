NAME := ngtd
VERSION := v0.0.1
GO_VERSION:=$(shell go version)
REVISION := $(shell git rev-parse --short HEAD)

.PHONY: build

docker:
	docker build -t yahoojapan/ngtd:latest .

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

test: proto/ngtd.pb.go
	GO111MODULE=on go test -v ./...
