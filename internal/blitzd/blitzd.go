package blitzd

import (
	"bufio"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzd/internal/api"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/internal/metric/lnd"
	"github.com/frennkie/blitzd/internal/metric/network"
	"github.com/frennkie/blitzd/internal/metric/system"
	"github.com/frennkie/blitzd/internal/serve"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/spf13/viper"
	"log"
	"net/http"
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

	//if _, err := os.Stat(customConfig); os.IsNotExist(err) {
	//	log.Printf("custom config file does not exist - skipping: %s", customConfig)
	//	return
	//}

	lnd.Init()
	network.Init()
	system.Init()

	if runtime.GOOS != "windows" {
		go UpdateLsbRelease()
		go UpdateFileBar()
	}

	box := rice.MustFindBox("../../web/")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))

	http.HandleFunc("/", serve.Root)
	// ToDo fix metrics
	http.HandleFunc("/info/", serve.Info(&metric.Metrics))

	// ToDo fix.. every sub url matches
	http.HandleFunc(data.APIv1, api.All())
	http.HandleFunc(data.APIv1+"config/", api.Config())

	RESTHostPort := viper.GetString("RESTHostPort")
	log.Printf("REST: Listening on host: http://%s", RESTHostPort)

	//rpcHostPort := viper.GetString("rpcHostPort")
	//log.Printf("RPC: Listening on host: gRPC://%s", rpcHostPort)

	// now ListenAndServer
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", RESTHostPort), nil))

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
