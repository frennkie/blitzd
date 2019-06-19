package servers

import (
	"context"
	"github.com/frennkie/blitzd/pkg/protocol/grpc"
	"github.com/frennkie/blitzd/pkg/service/v1"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	v1API := v1.NewServiceServer()

	return grpc.RunServer(ctx, v1API)

}
