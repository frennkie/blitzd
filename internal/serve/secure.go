package serve

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/web"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

// this currently uses basic auth.
// Mutual TLS could also be an idea:
// https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go

func secret(user, realm string) string {
	if user == viper.GetString("admin.username") {
		// password is "changeme"

		return viper.GetString("admin.password")
	}
	return ""
}

func Secure(metrics *data.Cache) {

	authenticator := auth.NewBasicAuthenticator("example.com", secret)

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

		secureSchema := "https"
		secureHost := "localhost"
		securePort := fmt.Sprintf("%d", viper.GetInt("server.https.port"))
		secureBase := fmt.Sprintf("%s://%s:%s", secureSchema, secureHost, securePort)

		values := []interface{}{secureBase, secureBase, secureBase, r.RemoteAddr, r.URL.Path}

		html := fmt.Sprintf(htmlRaw, values...)

		_, _ = fmt.Fprintf(w, "%s", html)

	})

	infoMux.HandleFunc("/api/", authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

		htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon - REST API</title>
</head>
<body>
	<h1>REST API</h1>
	<h2>Hello: %s</h2>
</body>
</html>`

		values := []interface{}{r.Username}

		html := fmt.Sprintf(htmlRaw, values...)

		_, _ = fmt.Fprintf(w, "%s", html)

	}))

	//infoMux.HandleFunc("/info/", func(w http.ResponseWriter, r *http.Request) {
	infoMux.HandleFunc("/info/", authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

		infoTemplate, err := vfstemplate.ParseFiles(web.Assets, template.New("info.tmpl"), "info.tmpl")
		if err != nil {
			log.Fatal(err)
		}

		if err := infoTemplate.Execute(w, metrics); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	}))

	infoHostPort := config.GetServerHttpsHostPort()
	log.Printf("Starting Secure Info/REST Server (https://%s)", infoHostPort)
	log.Fatal(http.ListenAndServeTLS(infoHostPort,
		viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
}
