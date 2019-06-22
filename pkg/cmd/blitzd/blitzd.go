package blitzd

import (
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/metric/bitcoind"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/metric/network"
	"github.com/frennkie/blitzd/internal/metric/raspiblitz"
	"github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/pkg/cmd/servers"
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

	if util.FileExists(config.C.Server.TlsCert) && util.FileExists(config.C.Server.TlsKey) {
		log.Printf("Using Key-Pair: %s;%s", config.C.Server.TlsCert, config.C.Server.TlsKey)
	} else {
		// ToDo add some checking (e.g. for existing files) here?!
		log.Printf("Need to generate Key-Pair")
		err := util.GenRootCaSignedClientServerCert(
			config.C.Alias,
			config.C.Server.CaCert,
			config.C.Server.TlsCert,
			config.C.Server.TlsKey,
			config.C.Client.TlsCert,
			config.C.Client.TlsKey,
		)
		if err != nil {
			log.Fatalf("Failed to generate Key-Pair for: Server: %s", err)
		}
	}

	log.Printf("Server Root CA: %s", config.C.Server.CaCert)
	log.Printf("Server TLS Cert: %s", config.C.Server.TlsCert)
	log.Printf("Server TLS Key: %s", config.C.Server.TlsKey)

	log.Printf("Client Root CA: %s", config.C.Client.CaCert)
	//log.Printf("Client TLS Cert: %s", config.C.Client.TlsCert)
	//log.Printf("Client TLS Key: %s", config.C.Client.TlsKey)

	if config.C.Server.Http.Enabled {
		log.Printf("HTTP Server: enabled (http://localhost:%d)", config.C.Server.Http.Port)
		go servers.Welcome()
	} else {
		log.Printf("HTTP Server: disabled")
	}

	if config.C.Server.Https.Enabled {
		log.Printf("HTTPS Server: enabled (https://localhost:%d)", config.C.Server.Https.Port)
		go servers.Secure()
	} else {
		log.Printf("HTTPS Server: disabled")
	}

	if config.C.Server.Rpc.Enabled {
		log.Printf("RPC Server: enabled (Port: %d)", config.C.Server.Rpc.Port)

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

	bitcoind.Init()
	lnd.Init()
	network.Init()
	raspiblitz.Init()
	system.Init()

	select {}

}
