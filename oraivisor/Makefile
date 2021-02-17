#!/usr/bin/make -f


all: oraivisor test

build:
	go build -mod=readonly ./cmd/oraivisor

test:
	go test -mod=readonly -race ./...

.PHONY: all build test
