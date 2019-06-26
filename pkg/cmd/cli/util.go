package cli

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/golang/protobuf/jsonpb"
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
	clientCert, err := tls.LoadX509KeyPair(config.C.Client.Tls.Cert, config.C.Client.Tls.Key)
	if err != nil {
		log.Printf("load client cert/key error:%v", err)
		return nil, err
	}

	serverRootCaCert, err := ioutil.ReadFile(config.C.Server.Tls.Ca)
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

	if config.Verbose {
		log.Printf("rpcAddress: %s", config.RpcHostPort)
	}

	conn, err := grpc.Dial(config.RpcHostPort, grpc.WithTransportCredentials(ta))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
