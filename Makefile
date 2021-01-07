PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags --always) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# TODO: Update the ldflags with the app, client & server names
ldflags = -s -w \
	-X github.com/cosmos/cosmos-sdk/version.Name=NewApp \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=oraid \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=oraicli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
		@echo "--> Installing oraid & oraicli"
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/orai
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/oraicli
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/oraid
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/websocket


install-orai:
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/orai

install-oraicli:
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/oraicli

install-oraid:
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/oraid

install-websocket:
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/websocket						

watch-oraid:
		air -c oraid.toml

watch-oraicli:
		air -c oraicli.toml

watch-websocket:
		air -c websocket.toml		

# always rebuild all the binaries after saving 
build: 
		go build -o ./tmp/orai -mod=readonly $(BUILD_FLAGS) ./cmd/orai
		go build -o ./tmp/oraid -mod=readonly $(BUILD_FLAGS) ./cmd/oraid
		go build -o ./tmp/oraid -mod=readonly $(BUILD_FLAGS) ./cmd/orai
		go build -o ./tmp/oraicli -mod=readonly $(BUILD_FLAGS) ./cmd/oraicli
		go build -o ./tmp/websocket -mod=readonly $(BUILD_FLAGS) ./cmd/websocket

		go install -mod=readonly $(BUILD_FLAGS) ./cmd/orai
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/oraid
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/orai
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/oraicli
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/websocket

build-orai: 
		go build -o ./tmp/orai -mod=readonly $(BUILD_FLAGS) ./cmd/orai

build-oraid: 
		go build -o ./tmp/oraid -mod=readonly $(BUILD_FLAGS) ./cmd/oraid	

build-oraicli: 
		go build -o ./tmp/oraicli -mod=readonly $(BUILD_FLAGS) ./cmd/oraicli	

build-websocket: 
		go build -o ./tmp/websocket -mod=readonly $(BUILD_FLAGS) ./cmd/websocket					

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

# Uncomment when you have some tests
# test:
# 	@go test -mod=readonly $(PACKAGES)

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify
