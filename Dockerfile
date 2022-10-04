FROM golang:1.18-alpine as builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add upx
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

# See https://github.com/CosmWasm/wasmvm/releases
ADD ./libwasmvm_muslc_112.a /lib/libwasmvm_muslc.a
# # RUN sha256sum /lib/libwasmvm_muslc.a | grep 39dc389cc6b556280cbeaebeda2b62cf884993137b83f90d1398ac47d09d3900

# RUN go mod tidy && go get ./...

# # force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN make build LEDGER_ENABLED=false BUILD_TAGS=muslc GOMOD_FLAGS= VERSION=0.41.0
RUN cp /workspace/oraid /bin/oraid
# RUN upx --best --lzma /workspace/oraid