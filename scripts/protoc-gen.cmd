REM ### run this from main dir ###
protoc -I. -Ithird_party --go_out=plugins=grpc:pkg                      api\proto\v1\helloworld.proto
protoc -I. -Ithird_party --go_out=plugins=grpc:pkg                      api\proto\v1\hello.proto
protoc -I. -Ithird_party --go_out=plugins=grpc:pkg                      api\proto\v1\blitzd.proto

protoc -I. -Ithird_party --grpc-gateway_out=logtostderr=true:pkg        api\proto\v1\blitzd.proto

MOVE pkg\api\proto\v1\* pkg\api\v1\
RMDIR pkg\api\proto /S /Q

protoc -I. -Ithird_party --swagger_out=logtostderr=true:api\swagger\v1  api\proto\v1\blitzd.proto
MOVE api\swagger\v1\api\proto\v1\blitzd.swagger.json api\swagger\v1
RMDIR api\swagger\v1\api /S /Q

go generate web\assets.go
go generate web\swagger.go
go generate web\swagger_json.go