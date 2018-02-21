NAME := ngtd
VERSION := v0.0.1
GO_VERSION:=$(shell go version)
REVISION := $(shell git rev-parse --short HEAD)

.PHONY: build deps

deps:
	glide init
	glide update
	glide install

proto/ngtd.pb.go: proto/ngtd.proto
	protoc --gofast_out=plugins=grpc:. proto/ngtd.proto

build: proto/ngtd.pb.go
	go build -ldflags="-w -s"

test: proto/ngtd.pb.go
	go test -v ./...
