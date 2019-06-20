package cli

import (
	"context"
	"fmt"
	pb "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var cmdGetAll = &cobra.Command{
	Use:   "all",
	Short: "gRPC: Get All Metrics",
	Run: func(cmd *cobra.Command, args []string) {
		getAll()
	},
}

func getAll() {

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

	c := pb.NewMetricServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMetricAll(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("an error occured: %v", err)
	}

	if r != nil {
		result, _ := jsonMarshaler.MarshalToString(r)
		fmt.Println(result)
		os.Exit(0)
	}

	os.Exit(1)
}
