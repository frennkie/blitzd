#!/bin/sh

### run this from main dir ###
protoc -I. -Ithird_party --go_out=plugins=grpc:pkg                      api/proto/v1/helloworld.proto
protoc -I. -Ithird_party --go_out=plugins=grpc:pkg                      api/proto/v1/hello.proto
protoc -I. -Ithird_party --go_out=plugins=grpc:pkg                      api/proto/v1/blitzd.proto

protoc -I. -Ithird_party --grpc-gateway_out=logtostderr=true:pkg        api/proto/v1/blitzd.proto

mv pkg/api/proto/v1/* pkg/api/v1/
rm -rf pkg/api/proto

protoc -I. -Ithird_party --swagger_out=logtostderr=true:api/swagger/v1  api/proto/v1/blitzd.proto
mv api/swagger/v1/api/proto/v1/blitzd.swagger.json api/swagger/v1
rm -rf api/swagger/v1/api

go generate web/assets.go
go generate web/swagger.go
go generate web/swagger_json.go