package serve

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

func Info(metrics *data.Cache) {

	infoMux := http.NewServeMux()

	// favicon
	infoMux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../web/favicon.png")
	})

	// static
	infoMux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		//box := rice.MustFindBox("../../web/")
		//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))
		http.ServeFile(w, r, "../../web/status.css")
		//http.FileServer(box.HTTPBox())
	})

	// slash
	infoMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	InfoHostPort := fmt.Sprintf("localhost:%d", viper.GetInt("server.https.port"))
	log.Printf("Starting Info Server (http://%s)", InfoHostPort)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s", InfoHostPort),
		viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
}
