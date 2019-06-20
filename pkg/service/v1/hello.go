//go:generate ...

// Package v1 implements a server for Greeter service.
package v1

import (
	"context"
	"github.com/frennkie/blitzd/pkg/api/v1"
	"log"
)

// server is used to implement v1.GreeterServer.
type serverHello struct{}

// NewHelloServiceServer creates
func NewHelloServiceServer() v1.GreeterServer {
	return &serverHello{}
}

// SayHello implements v1.GreeterServer
func (s *serverHello) SayHello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloResponse, error) {
	log.Printf("Received: %v", req.Name)
	return &v1.HelloResponse{Message: "Hello " + req.Name}, nil
}
