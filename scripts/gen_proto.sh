#!/usr/bin/env bash
path=$(cd `dirname $0`; pwd)/../driver/proto
path=$(cd ${path}; pwd)
# echo "protoc --proto_path=${path} --go_out=plugins=grpc:${path} ${path}/*.proto"
# protoc --proto_path=${path} --go_out=plugins=grpc:${path} ${path}/*.proto
# echo "protoc --go-grpc_out=. --go-grpc_opt paths=source_relative ${path}/*.proto"
# protoc --go-grpc_out=. --go-grpc_opt paths=source_relative ${path}/*.proto
protoc --proto_path=${path} --go_out=plugins=grpc:${path} ${path}/*.proto
