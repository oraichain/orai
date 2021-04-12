set -eo pipefail

BASEDIR=$(dirname $0)
PROJECTDIR=$BASEDIR/..
MODULE_SDK_DIR=$(realpath $PROJECTDIR/x)
COSMOS_SDK_DIR=${COSMOS_SDK_DIR:-$(go list -f "{{ .Dir }}" -m github.com/cosmos/cosmos-sdk)}

# Generate Go types from protobuf
for PROTO_DIR in $(ls $MODULE_SDK_DIR)
do   
  echo "processing: $MODULE_SDK_DIR/$PROTO_DIR"
  proto_files=$(find "x/$PROTO_DIR/types/" -maxdepth 4 -name '*.proto')  
  buf protoc \
    -I=. \
    -I=$COSMOS_SDK_DIR/third_party/proto \
    -I=$COSMOS_SDK_DIR/proto \
    --gocosmos_out=Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,plugins=interfacetype+grpc,paths=source_relative:. \
    --grpc-gateway_out=Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,paths=source_relative:. \
    $proto_files
done


