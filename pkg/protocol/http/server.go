package http

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/web/assets"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Welcome() {
	welcomeMux := http.NewServeMux()

	welcomeMux.Handle("/favicon.ico", http.FileServer(assets.AssetsFs))
	welcomeMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.AssetsFs)))

	// "/" matches everything
	welcomeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var baseUrls []string

		securePort := fmt.Sprintf("%d", config.C.Server.Https.Port)

		remoteAddrPort := r.RemoteAddr
		remoteAddrPortSplit := strings.Split(remoteAddrPort, ":")
		remoteAddrSplit := remoteAddrPortSplit[:len(remoteAddrPortSplit)-1]
		remoteAddr := strings.Join(remoteAddrSplit, ":")
		log.Printf("remoteAddr: %s", remoteAddr)

		clientIsRemote := remoteAddr != "127.0.0.1" && remoteAddr != "[::1]"

		var localAddr string

		if clientIsRemote {
			log.Printf("Client is Remote!")
			switch colonCount := strings.Count(r.Host, ":"); colonCount {
			case 0:
				log.Printf("ColonCount: %d", colonCount)
				panic("Somthing went wrong")
			case 1:
				//log.Printf("IPv4")
				rHostSplit := strings.Split(r.Host, ":")
				localAddr = rHostSplit[0]
			default:
				//log.Printf("IPv6")
				rHostSplit := strings.Split(r.Host, ":")
				rHostSplit = rHostSplit[:len(rHostSplit)-1]
				localAddr = strings.Join(rHostSplit, ":")
			}

			baseUrls = append(baseUrls, "https://"+localAddr+":"+securePort+"/")
		} else {
			log.Printf("Client is Local!")
			log.Printf(r.Host)
			if config.C.Server.Https.LocalhostOnly {
				baseUrls = util.GetBaseUrls("https", securePort, true, false)
			} else {
				baseUrls = util.GetBaseUrls("https", securePort, true, true)
			}

		}

		welcomeTemplate, err := vfstemplate.ParseFiles(assets.AssetsFs, template.New("welcome.tmpl"), "welcome.tmpl")
		if err != nil {
			log.Fatal(err)
		}

		if err := welcomeTemplate.Execute(w, baseUrls); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}

	})

	port := fmt.Sprintf("%d", config.C.Server.Http.Port)

	if config.C.Server.Http.LocalhostOnly {
		log.Printf("Starting Welcome Server: http://localhost:%s) / http://127.0.0.1:%s / http://[::1]:%s", port, port, port)
		go func() {

			//log.Fatal(graceful.ListenAndServe("127.0.0.1:"+port, welcomeMux))
			log.Fatal(http.ListenAndServe("127.0.0.1:"+port, welcomeMux))
		}()

		go func() {

			//log.Fatal(graceful.ListenAndServe("[::1]:"+port, welcomeMux))
			log.Fatal(http.ListenAndServe("[::1]:"+port, welcomeMux))
		}()

	} else {
		go func() {
			// ToDo: Get proper any here
			log.Printf("Starting Welcome Server (http://localhost:%s - and all other IPs)", port)
			//log.Fatal(graceful.ListenAndServe(":"+port, welcomeMux))
			log.Fatal(http.ListenAndServe(":"+port, welcomeMux))
		}()
	}

}
