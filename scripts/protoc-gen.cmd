REM ### run this from main dir ###
protoc -I. -Iapi\proto\v1 --go_out=plugins=grpc:pkg\api\v1               api\proto\v1\helloworld.proto
protoc -I. -Iapi\proto\v1 --go_out=plugins=grpc:pkg\api\v1               api\proto\v1\hello.proto
protoc -I. -Iapi\proto\v1 --go_out=plugins=grpc:pkg\api\v1               api\proto\v1\blitzd.proto
