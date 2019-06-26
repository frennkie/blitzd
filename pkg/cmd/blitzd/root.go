package blitzd

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/web/assets"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var RootCmd = &cobra.Command{
	Version: BuildVersion + " (built: " + BuildTime + ")",
	Use:     "blitzd",
	Short:   "RaspiBlitz Daemon",
	Long: `A service that retrieves and caches details about your RaspiBlitz.
More info at: https://github.com/frennkie/blitzd`,
	Run: func(cmd *cobra.Command, args []string) {
		Init()
	},
}

var DemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo Code",
	Run: func(cmd *cobra.Command, args []string) {
		http.Handle("/favicon.ico", http.FileServer(assets.Assets))
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.Assets)))

		log.Printf("HTTP Server: http://localhost:18080/")
		log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:18080"), nil))
	},
}

var GenCertCmd = &cobra.Command{
	Use:   "gencert",
	Short: "Generate Certificates",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("gencert called")
		//util.GenRootCa("ca.crt", "ca.key", "foobar", false)
		err := util.GenRootCaSignedClientServerCert(
			config.C.Alias,
			config.C.Server.CaCert,
			config.C.Server.TlsCert,
			config.C.Server.TlsKey,
			config.C.Client.TlsCert,
			config.C.Client.TlsKey,
		)
		if err != nil {
			log.Printf("an error occured")
		}
		log.Printf("success!")

	},
}

var GraceCmd = &cobra.Command{
	Use:   "grace",
	Short: "Grace",
	Run: func(cmd *cobra.Command, args []string) {
		lnd.Init()
	},
}
