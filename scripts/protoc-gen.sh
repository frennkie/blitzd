#!/bin/sh

protoc --proto_path=api/proto/v1 --go_out=plugins=grpc:pkg/api/v1 api/proto/v1/helloworld.proto
protoc --proto_path=api/proto/v1 --go_out=plugins=grpc:pkg/api/v1 api/proto/v1/hello.proto
protoc --proto_path=api/proto/v1 --go_out=plugins=grpc:pkg/api/v1 api/proto/v1/blitzd.proto
