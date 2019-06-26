package cli

import (
	"context"
	"fmt"
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var cmdFoo5 = &cobra.Command{
	Use:   "foo",
	Short: "gRPC: Get Metric Foo",
	Run: func(cmd *cobra.Command, args []string) {
		if jsonFlag && formattedFlag {
			fmt.Println("do not use both --json or --formatted simultaneously")
			os.Exit(1)
		}
		foo()
	},
}

func foo() {

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

	c := v1.NewMetricServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMetricFoo(ctx, &v1.GetMetricFooRequest{Api: "v1"})
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

		fmt.Println("---")
		fmt.Println(r.Metric.Expired)
		fmt.Println(r.Metric.Expired == v1.Tribool_TRIBOOL_TRUE)
		fmt.Println("---")

		os.Exit(0)
	}
}
