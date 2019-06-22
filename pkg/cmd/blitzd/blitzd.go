package blitzd

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/metric/bitcoind"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/metric/network"
	"github.com/frennkie/blitzd/internal/metric/raspiblitz"
	"github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/pkg/cmd/servers"
	"github.com/frennkie/blitzd/pkg/protocol/rest"
	log "github.com/sirupsen/logrus"
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

	log.WithFields(log.Fields{"tls": "server", "path": config.C.Server.CaCert}).Debug("Root CA")
	log.WithFields(log.Fields{"tls": "server", "path": config.C.Server.TlsCert}).Debug("TLS Cert")
	log.WithFields(log.Fields{"tls": "server", "path": config.C.Server.TlsKey}).Debug("TLS Key")

	log.WithFields(log.Fields{"tls": "client", "path": config.C.Client.CaCert}).Debug("Root CA")
	log.WithFields(log.Fields{"tls": "client", "path": config.C.Client.TlsCert}).Debug("TLS Cert")
	log.WithFields(log.Fields{"tls": "client", "path": config.C.Client.TlsKey}).Debug("TLS Key")

	log.WithFields(log.Fields{"server": "HTTP",
		"enabled":        config.C.Server.Http.Enabled,
		"localhost_only": config.C.Server.Http.LocalhostOnly,
		"port":           config.C.Server.Http.Port}).Debug("starting if enabled")
	if config.C.Server.Http.Enabled {
		go servers.Welcome()
	}

	log.WithFields(log.Fields{"server": "HTTPs",
		"enabled":        config.C.Server.Https.Enabled,
		"localhost_only": config.C.Server.Https.LocalhostOnly,
		"port":           config.C.Server.Https.Port}).Debug("starting if enabled")
	if config.C.Server.Https.Enabled {
		go servers.Secure()
	}

	log.WithFields(log.Fields{"server": "RPC",
		"enabled":        config.C.Server.Rpc.Enabled,
		"localhost_only": config.C.Server.Rpc.LocalhostOnly,
		"port":           config.C.Server.Rpc.Port}).Debug("starting if enabled")
	if config.C.Server.Rpc.Enabled {
		_ = servers.RunServer()

		//go func() {
		//	if err := cmd.RunServer(); err != nil {
		//		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		//		os.Exit(1)
		//	}
		//}()
	}

	log.WithFields(log.Fields{"server": "REST",
		"enabled":        config.C.Server.Rest.Enabled,
		"localhost_only": config.C.Server.Rest.LocalhostOnly,
		"port":           config.C.Server.Rest.Port}).Debug("starting if enabled")
	if config.C.Server.Rest.Enabled {

		ctx := context.Background()

		// run HTTP gateway
		go func() {
			_ = rest.RunServer(ctx,
				fmt.Sprintf("%d", config.C.Server.Rpc.Port),
				fmt.Sprintf("%d", config.C.Server.Rest.Port))
		}()
	}

	bitcoind.Init()
	lnd.Init()
	network.Init()
	raspiblitz.Init()
	system.Init()

	select {}

}
