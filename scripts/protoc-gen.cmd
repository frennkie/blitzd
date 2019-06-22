REM ### run this from main dir ###
protoc -I. -Iapi\proto\v1 --go_out=plugins=grpc:pkg\api\v1               api\proto\v1\helloworld.proto
protoc -I. -Iapi\proto\v1 --go_out=plugins=grpc:pkg\api\v1               api\proto\v1\hello.proto
protoc -I. -Iapi\proto\v1 --go_out=plugins=grpc:pkg\api\v1               api\proto\v1\blitzd.proto
protoc -I. -Iapi\proto\v1 --grpc-gateway_out=logtostderr=true:pkg\api\v1 api\proto\v1\blitzd.proto

REM ### protoc -I. -Iapi\proto\v1 --swagger_out=logtostderr=true:api\swagger\v1  api\proto\v1\blitzd.proto
REM ### MOVE api\swagger\v1\api\proto\v1\blitzd.swagger.json api\swagger\v1
REM ### RMDIR api\swagger\v1\api /S /Q

