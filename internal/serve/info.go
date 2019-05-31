package serve

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

/*
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

*/

/*

func Info(metrics *data.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

}

*/

/*

func Info(metrics *data.Cache) {

	infoMux := http.NewServeMux()

	// favicon
	infoMux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../web/favicon.png")
	})

	// static
	infoMux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		box := rice.MustFindBox("../../web/")
		infoMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))
		//http.ServeFile(w, r, "../../web/status.css")
		//http.FileServer(box.HTTPBox())
	})

	// slash
	infoMux.HandleFunc("/info/", func(w http.ResponseWriter, r *http.Request) {
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

	infoMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)

		_, _ = io.WriteString(w, "Not found\n")
	})

	infoHostPort := fmt.Sprintf("localhost:%d", viper.GetInt("server.https.port"))
	log.Printf("Starting Info Server (http://%s)", infoHostPort)
	log.Fatal(http.ListenAndServe(infoHostPort, infoMux))
}



*/

func Info(metrics *data.Cache) {

	infoMux := http.NewServeMux()

	// favicon
	infoMux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../web/favicon.png")
	})

	// static
	infoMux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		box := rice.MustFindBox("../../web/")
		infoMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))
		//http.ServeFile(w, r, "../../web/status.css")
		//http.FileServer(box.HTTPBox())
	})

	infoMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// "/" matches everything - so only respond to exactly "/", "/info" and "/info/"
		if r.URL.Path != "/" && r.URL.Path != "/info" && r.URL.Path != "/info/" {
			http.NotFound(w, r)
			return
		}

		htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon</title>
</head>
<body>
	<ul>
		<li><a href="https://%s:%s/">Info Page 23</a></li>
	</ul>
	<br>

	<hr>
	%s
	<br>

	<hr>
	Request:
	<pre>%s</pre>
</body>
</html>`

		infoHost := "localhost"
		infoPort := fmt.Sprintf("%d", viper.GetInt("server.https.port"))

		values := []interface{}{infoHost, infoPort, r.RemoteAddr, r.URL.Path}

		html := fmt.Sprintf(htmlRaw, values...)

		_, _ = fmt.Fprintf(w, "%s", html)

	})

	infoHostPort := fmt.Sprintf("localhost:%d", viper.GetInt("server.https.port"))
	log.Printf("Starting Info Server (http://%s)", infoHostPort)
	log.Fatal(http.ListenAndServe(infoHostPort, infoMux))
}
