FROM golang:1.19-alpine as builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add upx bash jq
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

ARG WASMVM_VERSION
ARG VERSION

WORKDIR /workspace
COPY orai/app/ /workspace/app
COPY orai/cmd/ /workspace/cmd
COPY orai/doc/statik /workspace/doc/statik
COPY orai/go.mod /workspace/
COPY orai/go.sum /workspace/
COPY orai/Makefile /workspace/

# See https://github.com/Oraichain/wasmvm/releases
RUN set -eux; \    
    export ARCH=$(uname -m); \
    wget -O /lib/libwasmvm_muslc.a https://github.com/oraichain/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.${ARCH}.a;

RUN go mod download

# # force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN make build LEDGER_ENABLED=false BUILD_TAGS=muslc GOMOD_FLAGS= VERSION=${VERSION}
RUN cp /workspace/build/oraid /bin/oraid
RUN upx --best --lzma /workspace/build/oraid
