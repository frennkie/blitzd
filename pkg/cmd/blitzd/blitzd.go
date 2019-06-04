package blitzd

import (
	"bufio"
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/metric/network"
	"github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/frennkie/blitzd/pkg/cmd/servers"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	// e.g. -ldflags '-X github.com/frennkie/blitzd/internal/blitzd.BuildTime=`date`'
	BuildVersion    = "unset"
	BuildTime       = "unset"
	BuildGitVersion = "unset"
)

func Init() {
	log.Printf("Starting version: %s, built at %s", BuildVersion, BuildTime)
	log.Printf("Git Version: %s", BuildGitVersion)

	if util.FileExists(viper.GetString("server.tlscert")) && util.FileExists(viper.GetString("server.tlskey")) {
		log.Printf("Using Key-Pair: %s;%s", viper.GetString("server.tlscert"), viper.GetString("server.tlskey"))
	} else {
		// ToDo add some checking (e.g. for existing files) here?!
		log.Printf("Need to generate Key-Pair")
		err := util.GenRootCaSignedClientServerCert(
			viper.GetString("alias"),
			viper.GetString("server.cacert"),
			viper.GetString("server.tlscert"),
			viper.GetString("server.tlskey"),
			viper.GetString("client.tlscert"),
			viper.GetString("client.tlskey"),
		)
		if err != nil {
			log.Fatalf("Failed to generate Key-Pair for: Server: %s", err)
		}
	}

	log.Printf("Server Root CA: %s", viper.GetString("server.cacert"))
	log.Printf("Server TLS Cert: %s", viper.GetString("server.tlscert"))
	log.Printf("Server TLS Key: %s", viper.GetString("server.tlskey"))

	log.Printf("Client Root CA: %s", viper.GetString("client.cacert"))
	log.Printf("Client TLS Cert: %s", viper.GetString("client.tlscert"))
	log.Printf("Client TLS Key: %s", viper.GetString("client.tlskey"))

	if viper.GetBool("server.http.enabled") {
		log.Printf("HTTP Server: enabled (http://localhost:%d)", viper.GetInt("server.http.port"))
		go servers.Welcome()
	} else {
		log.Printf("HTTP Server: disabled")
	}

	if viper.GetBool("server.https.enabled") {
		log.Printf("HTTPS Server: enabled (https://localhost:%d)", viper.GetInt("server.https.port"))
		go servers.Secure(&metric.Metrics)
	} else {
		log.Printf("HTTPS Server: disabled")
	}

	if viper.GetBool("server.rpc.enabled") {
		log.Printf("RPC Server: enabled (Port: %d)", viper.GetInt("server.rpc.port"))

		_ = servers.RunServer()

		//go func() {
		//	if err := cmd.RunServer(); err != nil {
		//		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		//		os.Exit(1)
		//	}
		//}()

	} else {
		log.Printf("RPC Server: disabled")
	}

	lnd.Init()
	network.Init()
	system.Init()

	if runtime.GOOS != "windows" {
		go UpdateLsbRelease()
		go UpdateFileBar()
	}

	// ToDo fix metrics
	//http.HandleFunc("/info/", serve.Secure(&metric.MetricsOld))

	select {}

}

func UpdateFileBar() {
	title := "file-bar"
	absFilePath := "/tmp/foo"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateFileBarFunc(title, absFilePath)
	go util.FileWatcher(title, absFilePath, UpdateFileBarFunc)
}

func UpdateFileBarFunc(title string, absFilePath string) {
	log.Printf("event-based update: %s (%s)", title, absFilePath)
	m := data.NewMetricEventBased(title)

	m.Value = fmt.Sprintf("%s", "foobar")

	metric.MetricsOldMux.Lock()
	metric.MetricsOld.FileBar = m
	metric.MetricsOldMux.Unlock()

}

// TODO replace with /etc/issue
func UpdateLsbRelease() {
	title := "lsb-release"
	absFilePath := "/etc/lsb-release"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateLsbReleaseFunc(title, absFilePath)
	go util.FileWatcher(title, absFilePath, UpdateLsbReleaseFunc)
}

func UpdateLsbReleaseFunc(title string, absFilePath string) {
	log.Printf("event-based update: %s (%s)", title, absFilePath)
	m := data.NewMetricEventBased(title)

	file, err := os.Open(absFilePath)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	_ = file.Close()

	tmp := txtlines[len(txtlines)-1]
	tmp2 := strings.Split(tmp, "=")[1]
	tmp3 := strings.Replace(tmp2, "\"", "", -1)
	m.Value = tmp3

	metric.MetricsOldMux.Lock()
	metric.MetricsOld.LsbRelease = m
	metric.MetricsOldMux.Unlock()

}
