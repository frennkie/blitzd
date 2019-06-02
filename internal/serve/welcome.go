package serve

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/web"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Welcome() {
	welcomeMux := http.NewServeMux()

	welcomeMux.Handle("/favicon.ico", http.FileServer(web.Assets))
	welcomeMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(web.Assets)))

	// "/" matches everything
	welcomeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var baseUrls []string

		securePort := fmt.Sprintf("%d", viper.GetInt("server.https.port"))

		remoteAddrPort := r.RemoteAddr
		remoteAddr := strings.Split(remoteAddrPort, ":")
		clientIsRemote := remoteAddr[0] != "127.0.0.1"

		if clientIsRemote {
			localAddr := strings.Split(r.Host, ":")
			baseUrls = append(baseUrls, "https://"+localAddr[0]+":"+securePort+"/")
		} else {

			log.Printf(r.Host)
			if viper.GetBool("server.https.localhost_only") {
				baseUrls = util.GetBaseUrls("https", securePort, true, false)
			} else {
				baseUrls = util.GetBaseUrls("https", securePort, true, true)
			}

		}

		welcomeTemplate, err := vfstemplate.ParseFiles(web.Assets, template.New("welcome.tmpl"), "welcome.tmpl")
		if err != nil {
			log.Fatal(err)
		}

		if err := welcomeTemplate.Execute(w, baseUrls); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}

	})

	port := fmt.Sprintf("%d", viper.GetInt("server.http.port"))

	if viper.GetBool("server.http.localhost_only") {
		log.Printf("Starting Welcome Server: http://localhost:%s) / http://127.0.0.1:%s / http://[::1]:%s", port, port, port)
		go func() {

			log.Fatal(http.ListenAndServe("127.0.0.1:"+port, welcomeMux))
		}()

		go func() {

			log.Fatal(http.ListenAndServe("[::1]:"+port, welcomeMux))
		}()

	} else {
		go func() {
			// ToDo: Get proper any here
			log.Printf("Starting Welcome Server (http://localhost:%s - and all other IPs)", port)
			log.Fatal(http.ListenAndServe(":"+port, welcomeMux))
		}()
	}

}
