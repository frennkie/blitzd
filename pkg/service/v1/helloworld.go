//go:generate ...

// Package v1 implements a server for Greeter service.
package v1

import (
	"context"
	"github.com/frennkie/blitzd/pkg/api/v1"
	"log"
)

// server is used to implement v1.GreeterServer.
type serverHelloWorldServer struct{}

// NewHelloServiceServer creates
func NewHelloWorldServiceServer() v1.HelloWorldGreeterServer {
	return &serverHelloWorldServer{}
}

// SayHello implements v1.GreeterServer
func (s *serverHelloWorldServer) SayHelloWorld(ctx context.Context, req *v1.HelloWorldRequest) (*v1.HelloWorldResponse, error) {
	log.Printf("Received: HelloWorld Request")
	return &v1.HelloWorldResponse{Message: "Hello World"}, nil
}
