#!/usr/bin/env bash

# apk add nodejs-current 
# npm install -g protobufjs

BASEDIR=$(dirname $0)
PROJECTDIR=$BASEDIR/..
# default is tmp folder
SOURCEDIR=$(realpath ${1:-$PROJECTDIR/tmp})
MODULE_SDK_DIR=$(realpath $PROJECTDIR/x)

COSMOS_SDK_DIR=${COSMOS_SDK_DIR:-$(go list -f "{{ .Dir }}" -m github.com/cosmos/cosmos-sdk)}
# scan all folders that contain proto file
proto_dirs=$(find $MODULE_SDK_DIR $COSMOS_SDK_DIR/proto $COSMOS_SDK_DIR/third_party/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
proto_files=()

for dir in $proto_dirs; do
  proto_files=("${proto_files[@]} $(find "${dir}" -maxdepth 1 -name '*.proto')")
done

# create dir & file if it does not exist
if [[ ! -d $SOURCEDIR/generated ]]; then
  mkdir $SOURCEDIR/generated
  touch $SOURCEDIR/generated/proto.js
fi

# echo ${proto_files[@]}

# gen files
pbjs \
  -o $SOURCEDIR/generated/proto.js \
  -t static-module \
  -w es6 \
  --es6 \
  --force-long \
  --keep-case \
  --no-create \
  ${proto_files[@]}

pbts \
  -o $SOURCEDIR/generated/proto.d.ts \
  $SOURCEDIR/generated/proto.js

# fix for node 14
# sed -i 's/import \* as/import/' $SOURCEDIR/generated/proto.js
node -e "var fs = require('fs'),file='$SOURCEDIR/generated/proto.js',result = fs.readFileSync(file).toString().replace('import * as', 'import');fs.writeFileSync(file, result)"

# show results
du -hd1 $SOURCEDIR/generated/*