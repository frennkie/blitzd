//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package v1

import (
	"context"
	"github.com/frennkie/blitzd/pkg/api/v1"
	"log"
)

// server is used to implement v1.GreeterServer.
type server struct{}

// NewServiceServer creates
func NewServiceServer() v1.GreeterServer {
	return &server{}
}

// SayHello implements v1.GreeterServer
func (s *server) SayHello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloReply, error) {
	log.Printf("Received: %v", req.Name)
	return &v1.HelloReply{Message: "Hello " + req.Name}, nil
}
