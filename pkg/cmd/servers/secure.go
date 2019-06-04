package servers

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/web"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/spf13/viper"
	"github.com/zenazn/goji/graceful"
	"html/template"
	"log"
	"net/http"
)

// this currently uses basic auth.
// Mutual TLS could also be an idea:
// https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go
// http://www.levigross.com/2015/11/21/mutual-tls-authentication-in-go/

func secret(user, realm string) string {
	if user == viper.GetString("admin.username") {
		// password is "changeme"

		return viper.GetString("admin.password")
	}
	return ""
}

func Secure(metrics *data.Cache) {

	authenticator := auth.NewBasicAuthenticator("RaspiBlitz Secure Server", secret)

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
		<li><a href="/foobar/">Foobar</a></li>
		<li><a href="/info/">Info</a></li>
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

		values := []interface{}{secureBase, r.RemoteAddr, r.URL.Path}

		html := fmt.Sprintf(htmlRaw, values...)

		_, _ = fmt.Fprintf(w, "%s", html)

	})

	infoMux.HandleFunc("/foobar/", authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

		htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon - REST API</title>
</head>
<body>
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

	port := fmt.Sprintf("%d", viper.GetInt("server.https.port"))

	if viper.GetBool("server.https.localhost_only") {
		log.Printf("Starting Secure Info Server: https://localhost:%s) / https://127.0.0.1:%s / https://[::1]:%s", port, port, port)
		go func() {

			log.Fatal(graceful.ListenAndServeTLS("127.0.0.1:"+port,
				viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
		}()

		go func() {

			log.Fatal(graceful.ListenAndServeTLS("[::1]:"+port,
				viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
		}()

	} else {
		go func() {
			// ToDo: Get proper any here
			log.Printf("Starting Secure Info Server (https://ANY:%s)", port)
			log.Fatal(graceful.ListenAndServeTLS(":"+port,
				viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
		}()
	}

}
