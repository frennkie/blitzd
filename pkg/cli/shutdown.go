package cli

import (
	"context"
	"fmt"
	pb "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var cmdShutdown = &cobra.Command{
	Use:   "shutdown",
	Short: "gRPC: Call Shutdown Service",
	Run: func(cmd *cobra.Command, args []string) {
		shutdown()
	},
}

func shutdown() {

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

	c := pb.NewShutdownClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.DoShutdown(ctx, &pb.ShutdownRequest{})
	if err != nil {
		log.Fatalf("an error occurred: %v", err)
	}

	if r != nil {
		fmt.Println(r)
	}

}
