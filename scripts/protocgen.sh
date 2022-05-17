#!/usr/bin/env bash

# set -eo pipefail

BASEDIR=$(dirname $0)
PROJECTDIR=$BASEDIR/..
MODULE_SDK_DIR=$(realpath $PROJECTDIR/x)
COSMOS_SDK_DIR="cosmos-sdk" # ${COSMOS_SDK_DIR:-$(go list -f "{{ .Dir }}" -m github.com/cosmos/cosmos-sdk)}

# Generate Go types from protobuf
for PROTO_DIR in $(ls $MODULE_SDK_DIR)
do
  echo "processing: $MODULE_SDK_DIR/$PROTO_DIR"
  proto_files=$(find "x/$PROTO_DIR/types/" -maxdepth 4 -name '*.proto')
  buf protoc \
    -I=$MODULE_SDK_DIR \
    -I=$COSMOS_SDK_DIR/third_party/proto \
    -I=$COSMOS_SDK_DIR/proto \
    --gocosmos_out=Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,plugins=interfacetype+grpc,paths=source_relative:. \
    --grpc-gateway_out=Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,paths=source_relative:. \
    $proto_files

  # move proto files to the right places
  cp -r $PROTO_DIR/* $MODULE_SDK_DIR/$PROTO_DIR
  rm -rf $PROTO_DIR
done

# if got this error: include directory /work/cosmos-sdk/proto is within include directory /work which is not allowed

# then remove the replace cosmos sdk line to run proto-gen first

# this script can be run in proto container (in docker-compose.yml file)