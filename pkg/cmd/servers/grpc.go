package servers

import (
	"context"
	"github.com/frennkie/blitzd/pkg/protocol/grpc"
	"github.com/frennkie/blitzd/pkg/service/v1"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	hello := v1.NewHelloServiceServer()
	helloWorld := v1.NewHelloWorldServiceServer()
	metric := v1.NewMetricServer()
	shutdown := v1.NewShutdownServer()

	return grpc.RunServer(ctx, hello, helloWorld, metric, shutdown)

}
