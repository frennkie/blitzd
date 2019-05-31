package serve

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/web"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func Welcome() {
	welcomeMux := http.NewServeMux()

	welcomeMux.Handle("/favicon.ico", http.FileServer(web.Assets))
	welcomeMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(web.Assets)))

	// "/" matches everything
	welcomeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		htmlRaw := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>BlitzInfo Daemon</title>
</head>
<body>
	<ul>
		<li><a href="https://%s:%s/">Secure Page</a></li>
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

	welcomeHostPort := config.GetServerHttpHostPort()
	log.Printf("Starting Welcome Server (http://%s)", welcomeHostPort)
	log.Fatal(http.ListenAndServe(welcomeHostPort, welcomeMux))
}
