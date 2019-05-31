package main

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/blitzd"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVar(&config.BlitzdDir, "dir",
		config.DefaultBlitzdDir, "blitzd home directory (default is $HOME/.blitzd")

	rootCmd.AddCommand(demoCmd)
	rootCmd.AddCommand(genCertCmd)

	_ = rootCmd.Execute()

}

var rootCmd = &cobra.Command{
	Version: blitzd.BuildVersion + " (built: " + blitzd.BuildTime + ")",
	Use:     "blitzd",
	Short:   "RaspiBlitz Daemon",
	Long: `A service that retrieves and caches details about your RaspiBlitz.
More info at: https://github.com/frennkie/blitzd`,
	Run: func(cmd *cobra.Command, args []string) {
		blitzd.Init()
	},
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling Request for: %s", r.URL)
	http.StripPrefix("/static/", http.FileServer(web.Assets))
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo Code",
	Run: func(cmd *cobra.Command, args []string) {
		http.Handle("/favicon.ico", http.FileServer(web.Assets))
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(web.Assets)))

		log.Printf("HTTP Server: http://localhost:18080/")
		log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:18080"), nil))
	},
}

var genCertCmd = &cobra.Command{
	Use:   "gencert",
	Short: "Generate Certificates",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("gencert called")
		//util.GenRootCa("ca.crt", "ca.key", "fobar", false)
		err := util.GenRootCaSignedClientServerCert(
			viper.GetString("server.cacert"),
			viper.GetString("server.tlscert"),
			viper.GetString("server.tlskey"),
			viper.GetString("client.tlscert"),
			viper.GetString("client.tlskey"),
		)
		if err != nil {
			log.Printf("an error occured")
		}
		log.Printf("success!")

	},
}
