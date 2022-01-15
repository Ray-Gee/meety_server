#!/bin/sh

protoc proto/hello.proto \
    --js_out=import_style=commonjs:client/src/hello \
    --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client/src/hello \
    --go-grpc_out=server/hello \
    --go_out=server/hello
