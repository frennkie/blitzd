package blitzd

import (
	"bufio"
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/metric/network"
	"github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/serve"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/spf13/viper"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	BuildVersion = "0.8.15" // ToDo (semi-)automatically update this?!
	BuildTime    = "0"      // ToDo set this a build time?!
)

func Init() {
	log.Printf("Starting version: %s, built at %s", BuildVersion, BuildTime)

	if util.FileExists(viper.GetString("servercrt")) && util.FileExists(viper.GetString("serverkey")) {
		log.Printf("Using Key-Pair: %s/%s", viper.GetString("servercrt"), viper.GetString("serverkey"))
	} else {
		log.Printf("Need to generate Key-Pair for: Server")
		err := util.GenCertPair(viper.GetString("servercrt"), viper.GetString("serverkey"))
		if err != nil {
			log.Fatalf("Failed to generate Key-Pair for: Server: %s", err)
		}
	}

	if util.FileExists(viper.GetString("clientcrt")) && util.FileExists(viper.GetString("clientkey")) {
		log.Printf("Using Key-Pair: %s/%s", viper.GetString("clientcrt"), viper.GetString("clientkey"))
	} else {
		log.Printf("Need to generate Key-Pair for: Client")
		err := util.GenCertPair(viper.GetString("clientcrt"), viper.GetString("clientkey"))
		if err != nil {
			log.Fatalf("Failed to generate Key-Pair for: Client: %s", err)
		}
	}

	lnd.Init()
	network.Init()
	system.Init()

	if runtime.GOOS != "windows" {
		go UpdateLsbRelease()
		go UpdateFileBar()
	}

	log.Printf("Client TLS Cert: %s", viper.GetString("client.tlscert"))
	log.Printf("Client TLS Key: %s", viper.GetString("client.tlskey"))

	log.Printf("Server TLS Cert: %s", viper.GetString("server.tlscert"))
	log.Printf("Server TLS Key: %s", viper.GetString("server.tlskey"))

	if viper.GetBool("server.http.enabled") {
		log.Printf("HTTP Server: enabled (Port: %d)", viper.GetInt("server.http.port"))
	} else {
		log.Printf("HTTP Server: disabled")
	}

	if viper.GetBool("server.https.enabled") {
		log.Printf("HTTPS Server: enabled (Port: %d)", viper.GetInt("server.https.port"))
	} else {
		log.Printf("HTTPS Server: disabled")
	}

	if viper.GetBool("server.rpc.enabled") {
		log.Printf("RPC Server: enabled (Port: %d)", viper.GetInt("server.rpc.port"))
	} else {
		log.Printf("RPC Server: disabled")
	}

	if viper.GetBool("server.http.enabled") {
		go serve.Welcome()
	}

	//if viper.GetBool("server.https.enabled") {
	//	go serve.Info(&metric.Metrics)
	//}

	select{}

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

	metric.MetricsMux.Lock()
	metric.Metrics.FileBar = m
	metric.MetricsMux.Unlock()

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

	metric.MetricsMux.Lock()
	metric.Metrics.LsbRelease = m
	metric.MetricsMux.Unlock()

}
