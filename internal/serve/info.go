package serve

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzd/internal/data"
	"html/template"
	"log"
	"net/http"
)

func Info(metrics *data.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
	}
}
