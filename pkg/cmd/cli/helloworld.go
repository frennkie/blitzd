package cli

import (
	"context"
	pb "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var CmdHelloWorld = &cobra.Command{
	Use:   "helloworld",
	Short: "gRPC: Print Hello World",
	Long:  `Print Hello World`,
	Run: func(cmd *cobra.Command, args []string) {
		helloWorld()
	},
}

func helloWorld() {

	conn, err := setupConnection()
	if err != nil {
		log.Fatalf("could not setup connection: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := pb.NewHelloWorldGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHelloWorld(ctx, &pb.HelloWorldRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", r.Message)
}
