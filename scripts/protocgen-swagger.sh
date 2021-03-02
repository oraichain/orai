#!/usr/bin/env bash

set -eo pipefail

# go get ./...
# apk add nodejs-current 
# npm install -g swagger-combine

BASEDIR=$(dirname $0)
PROJECTDIR=$BASEDIR/..
# default is tmp folder
SOURCEDIR=$(realpath ${1:-$PROJECTDIR/tmp})
MODULE_SDK_DIR=$(realpath $PROJECTDIR/x)
DOC_DIR=$(realpath $PROJECTDIR/doc)

COSMOS_SDK_DIR=${COSMOS_SDK_DIR:-$(go list -f "{{ .Dir }}" -m github.com/cosmos/cosmos-sdk)}
COSMOS_WASM_DIR=${COSMOS_WASM_DIR:-$(go list -f "{{ .Dir }}" -m github.com/CosmWasm/wasmd)}

# scan all folders that contain proto file
proto_dirs=$(find $MODULE_SDK_DIR $COSMOS_SDK_DIR/proto $COSMOS_WASM_DIR -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

GEN_DIR=$SOURCEDIR/swagger-gen
# clean swagger files
rm -rf $GEN_DIR
mkdir -p $GEN_DIR

for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then   
    buf protoc  \
    -I=. \
    -I="$COSMOS_WASM_DIR" \
    -I="$COSMOS_SDK_DIR/third_party/proto" \
    -I="$COSMOS_SDK_DIR/proto" --gocosmos_out=Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,plugins=interfacetype+grpc,paths=source_relative:$COSMOS_SDK_DIR \
    --swagger_out=$GEN_DIR \
    --swagger_opt=logtostderr=true,fqn_for_swagger_name=true,simple_operation_ids=true \
    $query_file
        
  fi
done

node -e "var fs = require('fs'),file='$COSMOS_SDK_DIR/client/docs/config.json',result = fs.readFileSync(file).toString().replace('./client','$COSMOS_SDK_DIR/client').replace(/.\/tmp-swagger-gen/g, '$GEN_DIR');
var baseModuleDir = '$GEN_DIR/x', obj = JSON.parse(result);
fs.readdirSync(baseModuleDir).forEach(dir=>{
   obj.apis.push({
        url: baseModuleDir + '/' + dir + (dir === 'wasm' ? '/internal' : '' ) + '/types/query.swagger.json',
        operationIds: {
            rename: {
                Params: dir[0].toUpperCase() + dir.slice(1) + 'Params'
            }    
        }
   });
});

fs.writeFileSync('$GEN_DIR/config.json', JSON.stringify(obj, null, 2));
"


# combine swagger files
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine $GEN_DIR/config.json -o $DOC_DIR/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true
