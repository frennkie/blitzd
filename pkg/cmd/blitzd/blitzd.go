package blitzd

import (
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/metric/network"
	"github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/pkg/cmd/servers"
	"github.com/spf13/viper"
	"log"
)

var (
	// e.g. -ldflags '-X github.com/frennkie/blitzd/internal/blitzd.BuildTime=`date`'
	BuildVersion    = "unset"
	BuildTime       = "unset"
	BuildGitVersion = "unset"
)

func Init() {
	log.Printf("Starting version: %s, built at %s", BuildVersion, BuildTime)
	log.Printf("Git Version: %s", BuildGitVersion)

	if util.FileExists(viper.GetString("server.tlscert")) && util.FileExists(viper.GetString("server.tlskey")) {
		log.Printf("Using Key-Pair: %s;%s", viper.GetString("server.tlscert"), viper.GetString("server.tlskey"))
	} else {
		// ToDo add some checking (e.g. for existing files) here?!
		log.Printf("Need to generate Key-Pair")
		err := util.GenRootCaSignedClientServerCert(
			viper.GetString("alias"),
			viper.GetString("server.cacert"),
			viper.GetString("server.tlscert"),
			viper.GetString("server.tlskey"),
			viper.GetString("client.tlscert"),
			viper.GetString("client.tlskey"),
		)
		if err != nil {
			log.Fatalf("Failed to generate Key-Pair for: Server: %s", err)
		}
	}

	log.Printf("Server Root CA: %s", viper.GetString("server.cacert"))
	log.Printf("Server TLS Cert: %s", viper.GetString("server.tlscert"))
	log.Printf("Server TLS Key: %s", viper.GetString("server.tlskey"))

	log.Printf("Client Root CA: %s", viper.GetString("client.cacert"))
	//log.Printf("Client TLS Cert: %s", viper.GetString("client.tlscert"))
	//log.Printf("Client TLS Key: %s", viper.GetString("client.tlskey"))

	if viper.GetBool("server.http.enabled") {
		log.Printf("HTTP Server: enabled (http://localhost:%d)", viper.GetInt("server.http.port"))
		go servers.Welcome()
	} else {
		log.Printf("HTTP Server: disabled")
	}

	if viper.GetBool("server.https.enabled") {
		log.Printf("HTTPS Server: enabled (https://localhost:%d)", viper.GetInt("server.https.port"))
		go servers.Secure(&metric.Metrics)
	} else {
		log.Printf("HTTPS Server: disabled")
	}

	if viper.GetBool("server.rpc.enabled") {
		log.Printf("RPC Server: enabled (Port: %d)", viper.GetInt("server.rpc.port"))

		_ = servers.RunServer()

		//go func() {
		//	if err := cmd.RunServer(); err != nil {
		//		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		//		os.Exit(1)
		//	}
		//}()

	} else {
		log.Printf("RPC Server: disabled")
	}

	lnd.Init()
	network.Init()
	system.Init()

	select {}

}
