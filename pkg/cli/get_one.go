package cli

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	pb "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "gRPC: Get a Metric by it's path (e.g. \"system.uptime\")",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if jsonFlag && formattedFlag {
			fmt.Println("do not use both --json or --formatted simultaneously")
			os.Exit(1)
		}
		if config.Verbose {
			fmt.Println("Args: " + strings.Join(args, " "))
			fmt.Println("jsonFlag: ", jsonFlag)
			fmt.Println("formattedFlag: ", formattedFlag)
		}
		get(args)
	},
}

func get(args []string) {

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
	r, err := c.GetMetricByPath(ctx, &pb.GetMetricByPathRequest{Path: args[0]})
	if err != nil {
		log.Fatalf("an error occurred: %v", err)
	}

	if r != nil {
		if jsonFlag {
			result, _ := jsonMarshaler.MarshalToString(r)
			fmt.Println(result)
		} else {
			if formattedFlag {
				fmt.Println(r.Metric.Text)
			} else {
				fmt.Println(r.Metric.Value)
			}
		}
		os.Exit(0)
	}
}
