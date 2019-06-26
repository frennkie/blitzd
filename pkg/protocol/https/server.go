package https

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/frennkie/blitzd/web/assets"
	"github.com/frennkie/blitzd/web/swagger"
	"github.com/goji/httpauth"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/philips/go-bindata-assetfs"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	authOpts = httpauth.AuthOptions{
		Realm:    "Restricted (Config)",
		AuthFunc: authFromConfig,
	}
)

func authFromConfig(username, password string, r *http.Request) bool {
	return username == config.C.Admin.Username &&
		util.CheckPasswordHash(password, config.C.Admin.Password)
}

func serveSwagger(mux *http.ServeMux) {
	_ = mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func Secure() {

	infoMux := http.NewServeMux()

	ctx := context.Background()

	demoAddr := fmt.Sprintf("localhost:%d", config.C.Server.Rpc.Port)

	// START

	// load peer cert/key, cacert
	clientCert, err := tls.LoadX509KeyPair(config.C.Client.TlsCert, config.C.Client.TlsKey)
	if err != nil {
		log.Fatalf("load client cert/key error:%v", err)
	}

	serverRootCaCert, err := ioutil.ReadFile(config.C.Server.CaCert)
	if err != nil {
		log.Fatalf("read ca cert file error:%v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(serverRootCaCert)

	ta := credentials.NewTLS(&tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	})

	if config.Verbose {
		log.Printf("rpcAddress: %s", config.RpcHostPort)
	}

	dopts := []grpc.DialOption{grpc.WithTransportCredentials(ta)}

	gwmux := runtime.NewServeMux()
	err = v1.RegisterMetricServiceHandlerFromEndpoint(ctx, gwmux, demoAddr, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}

	// END

	infoMux.Handle("/api/", gwmux)
	serveSwagger(infoMux)

	infoMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.Copy(w, strings.NewReader(v1.Swagger))
	})

	// favicon && /static
	infoMux.Handle("/favicon.ico", http.FileServer(assets.Assets))
	infoMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.Assets)))

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
			securePort := fmt.Sprintf("%d", config.C.Server.Https.Port)
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

			infoTemplate, err := vfstemplate.ParseFiles(assets.Assets, template.New("info.tmpl"), "info.tmpl")
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

	port := fmt.Sprintf("%d", config.C.Server.Https.Port)

	if config.C.Server.Https.LocalhostOnly {
		log.Printf("Starting Secure Info Server: https://localhost:%s) / https://127.0.0.1:%s / https://[::1]:%s", port, port, port)
		go func() {

			//log.Fatal(graceful.ListenAndServeTLS("127.0.0.1:"+port,
			//	config.C.Server.TlsCert, config.C.Server.TlsKey, infoMux))
			log.Fatal(http.ListenAndServeTLS("127.0.0.1:"+port,
				config.C.Server.TlsCert, config.C.Server.TlsKey, infoMux))
		}()

		go func() {

			//log.Fatal(graceful.ListenAndServeTLS("[::1]:"+port,
			//	config.C.Server.TlsCert, config.C.Server.TlsKey, infoMux))
			log.Fatal(http.ListenAndServeTLS("[::1]:"+port,
				config.C.Server.TlsCert, config.C.Server.TlsKey, infoMux))
		}()

	} else {
		go func() {
			// ToDo: Get proper ANY here?!
			log.Printf("Starting Secure Info Server (https://ANY:%s)", port)
			//log.Fatal(graceful.ListenAndServeTLS(":"+port,
			//	config.C.Server.TlsCert, config.C.Server.TlsKey, infoMux))
			log.Fatal(http.ListenAndServeTLS(":"+port,
				config.C.Server.TlsCert, config.C.Server.TlsKey, infoMux))
		}()
	}

}

// RunServer runs HTTP/REST gateway
func RunServerBasic(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterMetricServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	log.Println("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
