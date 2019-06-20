package cli

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

var jsonMarshaler = jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: true,
	Indent:       "  ",
	OrigName:     false,
	AnyResolver:  nil,
}

func setupConnection() (*grpc.ClientConn, error) {

	// load peer cert/key, cacert
	clientCert, err := tls.LoadX509KeyPair(viper.GetString("client.tlscert"), viper.GetString("client.tlskey"))
	if err != nil {
		log.Printf("load client cert/key error:%v", err)
		return nil, err
	}

	serverRootCaCert, err := ioutil.ReadFile(viper.GetString("server.cacert"))
	if err != nil {
		log.Printf("read ca cert file error:%v", err)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(serverRootCaCert)

	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	})

	rpcAddress := viper.GetString("rpcHostPort")

	if config.Verbose {
		log.Printf("rpcAddress: %s", rpcAddress)
	}

	conn, err := grpc.Dial(rpcAddress, grpc.WithTransportCredentials(ta))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
