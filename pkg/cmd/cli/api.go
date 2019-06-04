package cli

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	pb "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultName = "world"
)

var CmdApi = &cobra.Command{
	Use:   "api [string to echo]",
	Short: "AOU Echo anything to the screen",
	Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
		run()
	},
}

func run() {

	// load peer cert/key, cacert
	clientCert, err := tls.LoadX509KeyPair(viper.GetString("client.tlscert"), viper.GetString("client.tlskey"))
	if err != nil {
		log.Printf("load client cert/key error:%v", err)
		return
	}

	serverRootCaCert, err := ioutil.ReadFile(viper.GetString("server.cacert"))
	if err != nil {
		log.Printf("read ca cert file error:%v", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(serverRootCaCert)

	//_ = credentials.NewTLS(&tls.Config{
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	})

	rpcAddress := fmt.Sprintf("127.0.0.1:%d", viper.GetInt("server.rpc.port"))
	log.Printf("rpcAddress: %s", rpcAddress)
	//conn, err := grpc.Dial(rpcAddress)
	conn, err := grpc.Dial(rpcAddress, grpc.WithTransportCredentials(ta))
	if err != nil {
		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
