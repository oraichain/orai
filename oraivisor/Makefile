#!/usr/bin/make -f


all: oraivisor test

build:
	go get ./... && go build -mod=readonly -o build/oraivisor ./cmd/oraivisor

test:
	go test -mod=readonly -race ./...

.PHONY: all build test
