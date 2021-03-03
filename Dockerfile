FROM golang:1.15-alpine3.12 AS go-builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git jq bash ncurses upx
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /workspace
COPY . /workspace/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.13.0/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 39dc389cc6b556280cbeaebeda2b62cf884993137b83f90d1398ac47d09d3900

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN make build LEDGER_ENABLED=false BUILD_TAGS=muslc GOMOD_FLAGS=
RUN go get github.com/pwaller/goupx
RUN go get github.com/cosmtrek/air

# then remove
RUN rm -rf /workspace/*
