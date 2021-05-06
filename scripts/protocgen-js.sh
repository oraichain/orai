#!/usr/bin/env bash

set -eo pipefail

# go get ./...

if [ ! `command -v pbjs` ]; then  
  apk add nodejs-current npm
  npm install -g protobufjs
fi 

BASEDIR=$(dirname $0)
PROJECTDIR=$(realpath $BASEDIR/..)
# default is tmp folder
SOURCEDIR=$PROJECTDIR/tmp
MODULE_SDK_DIR=$PROJECTDIR/x

COSMOS_SDK_DIR=${COSMOS_SDK_DIR:-$(go list -f "{{ .Dir }}" -m github.com/cosmos/cosmos-sdk)}
COSMOS_WASM_DIR=${COSMOS_WASM_DIR:-$(go list -f "{{ .Dir }}" -m github.com/CosmWasm/wasmd)}

# scan all folders that contain proto file
proto_dirs=$(find $MODULE_SDK_DIR $COSMOS_SDK_DIR/proto $COSMOS_SDK_DIR/third_party/proto $COSMOS_WASM_DIR -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
proto_files=()

for dir in $proto_dirs; do
  proto_files=("${proto_files[@]} $(find "${dir}" -maxdepth 1 -name '*.proto')")
done

# create dir & file if it does not exist
mkdir -p $SOURCEDIR/generated  

# echo ${proto_files[@]}

# gen files
pbjs \
  -o $SOURCEDIR/generated/proto.js \
  -t static-module \
  -w es6 \
  --es6 \
  --force-long \
  --keep-case \
  --no-create "$@" \
  ${proto_files[@]}

pbts \
  -o $SOURCEDIR/generated/proto.d.ts \
  $SOURCEDIR/generated/proto.js

# show results
du -hd1 $SOURCEDIR/generated/*