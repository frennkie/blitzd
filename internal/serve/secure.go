package serve

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/web"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

/*
func Secure(metrics *data.CacheOld) http.HandlerFunc {
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

func Secure(metrics *data.CacheOld) http.HandlerFunc {
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



		InfoHostPort := fmt.Sprintf("localhost:%d", viper.GetInt("server.https.port"))
		log.Printf("Starting Secure Server (http://%s)", InfoHostPort)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s", InfoHostPort),
			viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
	}

}

*/

/*

func Secure(metrics *data.CacheOld) {

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
	log.Printf("Starting Secure Server (http://%s)", infoHostPort)
	log.Fatal(http.ListenAndServe(infoHostPort, infoMux))
}



*/

func Secure(metrics *data.Cache) {

	infoMux := http.NewServeMux()

	// favicon && /static
	infoMux.Handle("/favicon.ico", http.FileServer(web.Assets))
	infoMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(web.Assets)))

	infoMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// "/" matches everything - so only respond to exactly "/", "/info" and "/info/"
		if r.URL.Path != "/" && r.URL.Path != "/about" && r.URL.Path != "/about/" {
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
	<h1>About</h1>
	<ul>
		<li>Me: %s</li>
		<li><a href="%s/api/">REST API</a></li>
		<li><a href="%s/info/">Info</a></li>
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

		secureSchema := "http"
		secureHost := "localhost"
		securePort := fmt.Sprintf("%d", viper.GetInt("server.https.port"))
		secureBase := fmt.Sprintf("%s://%s:%s", secureSchema, secureHost, securePort)

		values := []interface{}{secureBase, secureBase, secureBase, r.RemoteAddr, r.URL.Path}

		html := fmt.Sprintf(htmlRaw, values...)

		_, _ = fmt.Fprintf(w, "%s", html)

	})

	infoMux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {

		htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon - REST API</title>
</head>
<body>
	<h1>REST API</h1>
</body>
</html>`

		_, _ = fmt.Fprintf(w, "%s", fmt.Sprintf(htmlRaw))

	})

	infoMux.HandleFunc("/info/", func(w http.ResponseWriter, r *http.Request) {

		infoTemplate, err := vfstemplate.ParseFiles(web.Assets, template.New("info.tmpl"), "info.tmpl")
		if err != nil {
			log.Fatal(err)
		}

		if err := infoTemplate.Execute(w, metrics); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	})

	infoHostPort := fmt.Sprintf("localhost:%d", viper.GetInt("server.https.port"))
	log.Printf("Starting Secure Info/REST Server (https://%s)", infoHostPort)
	//log.Fatal(http.ListenAndServe(infoHostPort, infoMux))
	log.Fatal(http.ListenAndServeTLS(infoHostPort,
		viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
}
