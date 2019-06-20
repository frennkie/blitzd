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

var cmdGetFoo5 = &cobra.Command{
	Use:   "foo5",
	Short: "gRPC: Get Metric Foo5 (Subcommand)",
	Run: func(cmd *cobra.Command, args []string) {
		getfoo5()
	},
}

func getfoo5() {

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
	r, err := c.GetMetricFoo5(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("an error occured: %v", err)
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
