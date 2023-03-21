FROM golang:1.18-alpine as builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add upx bash jq
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /workspace
COPY app/ /workspace/app
COPY cmd/ /workspace/cmd
COPY packages/ /workspace/packages
COPY doc/statik /workspace/doc/statik
COPY go.mod /workspace/
COPY go.sum /workspace/
COPY Makefile /workspace/

# See https://github.com/Oraichain/wasmvm/releases (tag v1.2.0)
RUN wget -O /lib/libwasmvm_muslc.a https://github.com/oraichain/wasmvm/releases/download/v1.2.0/libwasmvm_muslc.a

# RUN go mod tidy && go get ./...

# # force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN make build LEDGER_ENABLED=false BUILD_TAGS=muslc GOMOD_FLAGS= VERSION=0.41.3
RUN cp /workspace/build/oraid /bin/oraid
RUN upx --best --lzma /workspace/build/oraid