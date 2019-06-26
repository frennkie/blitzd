package blitzd

import (
	"context"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/metric/bitcoind"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/metric/network"
	"github.com/frennkie/blitzd/internal/metric/raspiblitz"
	"github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/pkg/protocol/grpc"
	"github.com/frennkie/blitzd/pkg/protocol/http"
	"github.com/frennkie/blitzd/pkg/protocol/https"
	v1 "github.com/frennkie/blitzd/pkg/service/v1"
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

	if util.FileExists(config.C.Server.Tls.Cert) && util.FileExists(config.C.Server.Tls.Key) {
		log.Printf("Using Key-Pair: %s;%s", config.C.Server.Tls.Cert, config.C.Server.Tls.Key)
	} else {
		// ToDo add some checking (e.g. for existing files) here?!
		log.Printf("Need to generate Key-Pair")
		err := util.GenRootCaSignedClientServerCert(
			config.C.Alias,
			config.C.Server.Tls.Ca,
			config.C.Server.Tls.Cert,
			config.C.Server.Tls.Key,
			config.C.Client.Tls.Cert,
			config.C.Client.Tls.Key,
		)
		if err != nil {
			log.Fatalf("Failed to generate Key-Pair for: Server: %s", err)
		}
	}

	log.WithFields(log.Fields{"tls": "server", "path": config.C.Server.Tls.Ca}).Debug("Root CA")
	log.WithFields(log.Fields{"tls": "server", "path": config.C.Server.Tls.Cert}).Debug("TLS Cert")
	log.WithFields(log.Fields{"tls": "server", "path": config.C.Server.Tls.Key}).Debug("TLS Key")

	log.WithFields(log.Fields{"tls": "client", "path": config.C.Client.Tls.Ca}).Debug("Root CA")
	log.WithFields(log.Fields{"tls": "client", "path": config.C.Client.Tls.Cert}).Debug("TLS Cert")
	log.WithFields(log.Fields{"tls": "client", "path": config.C.Client.Tls.Ca}).Debug("TLS Key")

	log.WithFields(log.Fields{"server": "HTTP",
		"enabled":        config.C.Server.Http.Enabled,
		"localhost_only": config.C.Server.Http.LocalhostOnly,
		"port":           config.C.Server.Http.Port}).Debug("starting server if enabled")
	if config.C.Server.Http.Enabled {
		go http.Welcome()
	}

	log.WithFields(log.Fields{"server": "HTTPS",
		"enabled":        config.C.Server.Https.Enabled,
		"localhost_only": config.C.Server.Https.LocalhostOnly,
		"port":           config.C.Server.Https.Port}).Debug("starting server if enabled")
	if config.C.Server.Https.Enabled {
		go https.Secure()

		log.WithFields(log.Fields{"server": "HTTPS", "feature": "rest-server",
			"enabled": config.C.Server.Https.Rest.Enabled}).Debug("activating feature if enabled")

		log.WithFields(log.Fields{"server": "HTTPS", "feature": "rest-decs",
			"enabled": config.C.Server.Https.Rest.Enabled}).Debug("activating feature if enabled")

	}

	log.WithFields(log.Fields{"server": "RPC",
		"enabled":        config.C.Server.Grpc.Enabled,
		"localhost_only": config.C.Server.Grpc.LocalhostOnly,
		"port":           config.C.Server.Grpc.Port}).Debug("starting server if enabled")
	if config.C.Server.Grpc.Enabled {
		ctx := context.Background()

		hello := v1.NewHelloServiceServer()
		helloWorld := v1.NewHelloWorldServiceServer()
		metric := v1.NewMetricServer()
		shutdown := v1.NewShutdownServer()

		_ = grpc.RunServer(ctx, hello, helloWorld, metric, shutdown)

		//go func() {
		//	if err := cmd.RunServer(); err != nil {
		//		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		//		os.Exit(1)
		//	}
		//}()
	}

	//if config.C.Server.Rest.Enabled {
	//
	//	ctx := context.Background()
	//
	//	// run HTTP gateway
	//	go func() {
	//		_ = https.RunServer(ctx,
	//			fmt.Sprintf("%d", config.C.Server.Grpc.Port),
	//			fmt.Sprintf("%d", config.C.Server.Rest.Port))
	//	}()
	//}

	bitcoind.Init()
	lnd.Init()
	network.Init()
	raspiblitz.Init()
	system.Init()

	select {}

}
