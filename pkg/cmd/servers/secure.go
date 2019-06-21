package servers

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/frennkie/blitzd/web"
	"github.com/goji/httpauth"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

var (
	authOpts = httpauth.AuthOptions{
		Realm:    "Restricted (Config)",
		AuthFunc: authFromConfig,
	}
)

func authFromConfig(username, password string, r *http.Request) bool {
	return username == viper.GetString("admin.username") &&
		util.CheckPasswordHash(password, viper.GetString("admin.password"))
}

func Secure() {

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
		<li><a href="/foobar/">Foobar</a> (password protected)</li>
		<li><a href="/info/">Info</a> (password protected)</li>
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

	infoMux.Handle("/foobar/",
		httpauth.BasicAuth(authOpts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon - /foobar/</title>
</head>
<body>
	<h2>Hello: %s</h2>
</body>
</html>`

			values := []interface{}{"foo2"}

			html := fmt.Sprintf(htmlRaw, values...)

			_, _ = fmt.Fprintf(w, "%s", html)
		})))

	infoMux.Handle("/info/",
		httpauth.BasicAuth(authOpts)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			infoTemplate, err := vfstemplate.ParseFiles(web.Assets, template.New("info.tmpl"), "info.tmpl")
			if err != nil {
				log.Fatal(err)
			}

			//// get copy of current Cache
			var mMap = make(map[string]*v1.Metric)
			//var mMap = make(map[string]string)
			var m = data.Cache.Items()

			for _, v := range m {
				metricObject := interface{}(v.Object).(v1.Metric)
				mMap[fmt.Sprintf("%s.%s", metricObject.Module, metricObject.Title)] = &metricObject
				//mMap[fmt.Sprintf("%s.%s", metricObject.Module, metricObject.Title)] = metricObject.Text
			}

			//var m v1.Metric
			//if x, found := data.Cache.Get("system.uptime"); found {
			//	m = x.(v1.Metric)
			//}

			if err := infoTemplate.Execute(w, mMap); err != nil {
				log.Println(err.Error())
				http.Error(w, http.StatusText(500), 500)
			}

		})))

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
