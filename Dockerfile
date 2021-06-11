FROM cosmwasm/prototools-docker AS go-builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk --no-cache add build-base jq bash ncurses;
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /workspace
# COPY . /workspace/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.14.0/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 220b85158d1ae72008f099a7ddafe27f6374518816dd5873fd8be272c5418026

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
# RUN make build GOMOD_FLAGS=
# RUN go get github.com/cosmtrek/air

# # then remove
# RUN rm -rf /workspace/*
