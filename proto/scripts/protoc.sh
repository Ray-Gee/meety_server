#!/bin/sh

set -xe

SERVER_OUTPUT_DIR=server/grpc
CLIENT_OUTPUT_DIR=client/src/messenger
GOOGLEAPIS_DIR=/

protoc --version
protoc --proto_path=proto messenger.proto \
  --go_out=plugins="grpc:${SERVER_OUTPUT_DIR}" \
  --go_opt=module=github.com/Ryuichi-g/meety_server/tree/main/proto \
  --js_out=import_style=commonjs:${CLIENT_OUTPUT_DIR} \
  --grpc-web_out=import_style=typescript,mode=grpcwebtext:${CLIENT_OUTPUT_DIR}
# protoc --proto_path=proto crud.proto \
#   -I${GOOGLEAPIS_DIR} -I. --include_imports --include_source_info \
#   --descriptor_set_out=server/grpc/crud/crud_descriptor.pb server/grpc/crud/crud_descriptor.proto
#   --go_out=plugins="grpc:${SERVER_OUTPUT_DIR}" \
#   --go_opt=module=github.com/Ryuichi-g/meety_server/tree/main/proto \
#   --js_out=import_style=commonjs:${CLIENT_OUTPUT_DIR} \
#   --grpc-web_out=import_style=typescript,mode=grpcwebtext:${CLIENT_OUTPUT_DIR}