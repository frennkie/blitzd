package main

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzd/internal/blitzd"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)


var RootCmd = &cobra.Command{
	Version: blitzd.BuildVersion,
	Use:     "blitzd",
	Short:   "RaspiBlitz Daemon",
	Long: `A service that retrieves and caches details about your RaspiBlitz.
More info at: https://github.com/frennkie/blitzd`,
	Run: func(cmd *cobra.Command, args []string) {
		go serveInfo(&data.Cache{})
		blitzd.Init()

	},
}

func main() {
	cobra.OnInitialize(config.InitConfig)

	RootCmd.PersistentFlags().StringVar(&config.BlitzdDir, "dir",
		config.DefaultBlitzdDir, "blitzd home directory (default is $HOME/.blitzd")

	_ = RootCmd.Execute()

}



func serveInfo(metrics *data.Cache) {
	//http := http.NewServeMux()

	// slash
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving / request for: %s", r.RequestURI)

		// find rice.Box
		templateBox, err := rice.FindBox("../../web")
		if err != nil {
			log.Fatal(err)
		}
		// get file contents as string
		templateString, err := templateBox.String("info.tmpl")
		if err != nil {
			log.Fatal(err)
		}
		// parse and execute the template
		tmplMessage, err := template.New("info").Parse(templateString)
		if err != nil {
			log.Fatal(err)
		}

		if err := tmplMessage.Execute(w, metrics); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	})

	// favicon
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("serving favicon request for: %s", r.RequestURI)
		http.ServeFile(w, r, "../../web/favicon.png")
	})

	// static
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("serving /static/ request for: %s", r.RequestURI)

		//w.Header().Add("Content-Type","text/plain")

		http.StripPrefix("/static/", http.FileServer(rice.MustFindBox("../../web/").HTTPBox()))
	})

	InfoHostPort := fmt.Sprintf("localhost:%d", viper.GetInt("server.https.port"))

	log.Printf("Starting Info Server (http://%s)", InfoHostPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", InfoHostPort), nil))

	//log.Printf("Starting Info Server (https://%s)", InfoHostPort)
	//log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s", InfoHostPort),
	//	viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), nil))
	//

}