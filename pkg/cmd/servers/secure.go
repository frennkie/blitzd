package servers

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/web"
	"github.com/goji/httpauth"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"runtime"
)

var (
	authOptsPam = httpauth.AuthOptions{
		Realm:    "Restricted (PAM)",
		AuthFunc: authFromPam,
	}

	authOptsConfig = httpauth.AuthOptions{
		Realm:    "Restricted (Config)",
		AuthFunc: authFromConfig,
	}

	authOpts httpauth.AuthOptions
)

func authFromPam(username, password string, r *http.Request) bool {
	// ToDo(frennkie) do something
	_ = util.CheckUser(username)
	return util.CheckUserPassword(username, password)
}

func authFromConfig(username, password string, r *http.Request) bool {
	return username == viper.GetString("admin.username") &&
		util.CheckPasswordHash(password, viper.GetString("admin.password"))
}

// ToDo remove this and all usages!
func authSecretFromConfig(user, _ string) string {
	if user == viper.GetString("admin.username") {
		// default password is "changeme"
		return viper.GetString("admin.password")
	}
	return ""
}

func Secure(metrics *data.Cache) {

	switch runtime.GOOS {
	case "windows":
		authOpts = authOptsConfig
	case "linux":
		authOpts = authOptsPam
	default:
		authOpts = authOptsConfig
	}

	authenticator := auth.NewBasicAuthenticator("RaspiBlitz Secure Server", authSecretFromConfig)

	infoMux := http.NewServeMux()

	// favicon && /static
	infoMux.Handle("/favicon.ico", http.FileServer(web.Assets))
	infoMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(web.Assets)))

	infoMux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			// "/" matches everything - so only respond to exactly "/", "/about" and "/about/"
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

	infoMux.Handle("/foobar2/",
		httpauth.BasicAuth(authOpts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon - REST API</title>
</head>
<body>
	<h2>Hello (Config|PAM): %s</h2>
</body>
</html>`

			values := []interface{}{"foo2"}

			html := fmt.Sprintf(htmlRaw, values...)

			_, _ = fmt.Fprintf(w, "%s", html)
		})))

	infoMux.Handle("/foobar3/",
		httpauth.BasicAuth(authOptsConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon - REST API</title>
</head>
<body>
	<h2>Hello (Config): %s</h2>
</body>
</html>`

			values := []interface{}{"foo3"}

			html := fmt.Sprintf(htmlRaw, values...)

			_, _ = fmt.Fprintf(w, "%s", html)
		})))

	infoMux.Handle("/foobar4/",
		httpauth.BasicAuth(authOptsPam)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon - REST API</title>
</head>
<body>
	<h2>Hello (PAM): %s</h2>
</body>
</html>`

			values := []interface{}{"foo4"}

			html := fmt.Sprintf(htmlRaw, values...)

			_, _ = fmt.Fprintf(w, "%s", html)
		})))

	infoMux.HandleFunc("/info/",
		authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

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

			//log.Fatal(graceful.ListenAndServeTLS("127.0.0.1:"+port,
			//	viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
			log.Fatal(http.ListenAndServeTLS("127.0.0.1:"+port,
				viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
		}()

		go func() {

			//log.Fatal(graceful.ListenAndServeTLS("[::1]:"+port,
			//	viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
			log.Fatal(http.ListenAndServeTLS("[::1]:"+port,
				viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
		}()

	} else {
		go func() {
			// ToDo: Get proper ANY here?!
			log.Printf("Starting Secure Info Server (https://ANY:%s)", port)
			//log.Fatal(graceful.ListenAndServeTLS(":"+port,
			//	viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
			log.Fatal(http.ListenAndServeTLS(":"+port,
				viper.GetString("server.tlscert"), viper.GetString("server.tlskey"), infoMux))
		}()
	}

}
