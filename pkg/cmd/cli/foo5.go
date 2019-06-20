package cli

import (
	"context"
	"fmt"
	pb "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var CmdFoo5 = &cobra.Command{
	Use:   "foo5",
	Short: "gRPC: Get Metric Foo",
	Run: func(cmd *cobra.Command, args []string) {
		foo5()
	},
}

func foo5() {

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
	r, err := c.GetMetricFoo5(ctx, &pb.GetMetricRequest{})
	if err != nil {
		log.Fatalf("an error occured: %v", err)
	}

	if r != nil {
		fmt.Println(proto.MarshalTextString(r.Metric))

		result, _ := jsonMarshaler.MarshalToString(r)
		fmt.Println(result)
	}

}
