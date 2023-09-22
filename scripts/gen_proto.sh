#!/usr/bin/env bash
path=$(cd `dirname $0`; pwd)/../driver/proto
echo "protoc --proto_path=${path} --go_out=plugins=grpc:${path} ${path}/*.proto"
protoc --proto_path=${path} --go_out=plugins=grpc:${path} ${path}/*.proto