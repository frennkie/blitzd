package main

import (
	"encoding/json"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzinfod/internal/utils"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/frennkie/blitzinfod/internal/data"
)

const (
	defaultRESTPort     = "7080"
	defaultRESTHostPort = "localhost:" + defaultRESTPort

	//defaultRPCPort     = "39735"
	//defaultRPCHostPort = "localhost:" + defaultRPCPort

	defaultExpireTime = 300 * time.Second // 5 minutes

	Normal = "normal"
	Red    = "red"
	Green  = "green"
	Yellow = "yellow"
	Purple = "purple"

	newline = "\r\n" // TODO windows newline
)

var (
	metrics    data.Cache
	metricsMux sync.Mutex

	// maxTime (Metric does not expire): "3000-01-01T00:00:00Z"
	maxTime = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)

	buildVersion = "0.8.15"
	buildTime    = "0"
)

func NewMetric(title string) data.Metric {
	metric := data.Metric{}

	metric.Title = title

	metric.Interval = time.Duration(5 * time.Second).Seconds()
	metric.Timeout = time.Duration(10 * time.Second).Seconds()

	metric.Value = "N/A"
	metric.Text = "N/A"
	metric.Prefix = ""
	metric.Suffix = ""
	metric.Style = Purple

	now := time.Now()
	metric.UpdatedAt = now
	metric.ExpiredAfter = now.Add(defaultExpireTime)

	return metric
}

// SetOperatingSystem sets the "os" from "runtime.GOOS" and returns it as a "Metric"
func SetOperatingSystem() data.Metric {
	metric := NewMetric("os")

	metric.Interval = 0
	metric.Timeout = 0
	metric.ExpiredAfter = maxTime

	metric.Value = runtime.GOOS

	return metric
}

func SetArch() data.Metric {
	metric := NewMetric("arch")

	metric.Interval = 0
	metric.Timeout = 0
	metric.ExpiredAfter = maxTime

	metric.Value = runtime.GOARCH

	return metric
}

func UpdateFoo() {
	// "warm-up" ToDo(frennkie) remove "Foo"
	time.Sleep(1 * time.Second)

	for {
		foo := NewMetric("foo")
		foo.Value = "foo"

		// update data in MetricCache
		log.Printf("Updating: %s", foo.Title)
		metricsMux.Lock()
		metrics.Foo = foo
		metricsMux.Unlock()

		time.Sleep(time.Duration(foo.Interval) * time.Second)
	}

}

func UpdateUptime() {
	for {
		m := NewMetric("uptime")
		m.Interval = time.Duration(1 * time.Second).Seconds()

		uptime, err := host.Uptime()
		if err != nil {
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {
			m.Value = fmt.Sprintf("%d", uptime)

			// update data in MetricCache
			log.Printf("Updating: %s", m.Title)
			metricsMux.Lock()
			metrics.Uptime = m
			metricsMux.Unlock()

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}

func UpdateNslookup() {
	for {
		m := NewMetric("nslookup")
		m.Interval = time.Duration(60 * time.Second).Seconds()

		cmdName := "nslookup"
		var cmdArgs []string

		if runtime.GOOS == "windows" {
			cmdArgs = []string{"google.com"}
		} else {
			cmdArgs = []string{"google.com"}
		}

		result, err := utils.TimeoutExec(cmdName, cmdArgs)
		if err != nil {
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {
			split := strings.Split(result, newline)
			last := strings.TrimSpace(split[len(split)-3])

			m.Value = last

			// update data in MetricCache
			log.Printf("Updating: %s", m.Title)
			metricsMux.Lock()
			metrics.Nslookup = m
			metricsMux.Unlock()

		}
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}

func UpdatePing() {
	for {
		m := NewMetric("ping")
		m.Interval = time.Duration(60 * time.Second).Seconds()

		cmdName := "ping"
		var cmdArgs []string

		if runtime.GOOS == "windows" {
			cmdArgs = []string{"-n", "2", "8.8.8.8"}
		} else {
			cmdArgs = []string{"-c", "2", "8.8.8.8"}
		}

		result, err := utils.TimeoutExec(cmdName, cmdArgs)
		if err != nil {
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {

			split := strings.Split(result, newline)
			last := strings.TrimSpace(split[len(split)-2])

			m.Value = last

			// update data in MetricCache
			log.Printf("Updating: %s", m.Title)
			metricsMux.Lock()
			metrics.Ping = m
			metricsMux.Unlock()

		}
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Version: buildVersion,
		Use:     "model",
		Short:   "RaspiBlitz Info Daemon",
		Long: `A service that retrieves and caches details about your RaspiBlitz.
                More info at: https://github.com/frennkie/model`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			blitzinfod()

		},
	}

	rootCmd.PersistentFlags().StringP("restHostPort", "R", defaultRESTHostPort, "REST: Listen on Host:Port")
	_ = viper.BindPFlag("restHostPort", rootCmd.PersistentFlags().Lookup("restHostPort"))
	viper.SetDefault("restHostPort", defaultRESTHostPort)

	//rootCmd.PersistentFlags().StringP("rpcHostPort", "H", defaultRPCHostPort, "RPC: Listen on Host:Port")
	//_ = viper.BindPFlag("rpcHostPort", rootCmd.PersistentFlags().Lookup("rpcHostPort"))
	//viper.SetDefault("rpcHostPort", defaultRPCHostPort)

	_ = rootCmd.Execute()
}

func blitzinfod() {
	log.Printf("Starting version: %s, built at %s", buildVersion, buildTime)

	// set static Metrics
	metrics.Arch = SetArch()
	metrics.OperatingSystem = SetOperatingSystem()

	// start Update of Metrics as goroutines
	go UpdateFoo()
	go UpdateUptime()
	go UpdateNslookup()
	go UpdatePing()

	box := rice.MustFindBox("../../web/")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))

	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/info/", serveInfo)
	http.HandleFunc("/api/", serveStaticApi)

	RESTHostPort := viper.GetString("RESTHostPort")
	log.Printf("REST: Listening on host: http://%s", RESTHostPort)

	//rpcHostPort := viper.GetString("rpcHostPort")
	//log.Printf("RPC: Listening on host: gRPC://%s", rpcHostPort)

	// now ListenAndServer
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", RESTHostPort), nil))

}

func serveStaticApi(w http.ResponseWriter, _ *http.Request) {

	js, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func serveRoot(w http.ResponseWriter, r *http.Request) {

	htmlRaw := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<title>BlitzInfo Daemon</title>
	</head>
	<body>
		<ul>	
			<li><a href="/api/">API</a></li>
			<li><a href="/info/">Info Page</a></li>
		</ul>
		<br>

		<hr>
		%s
		<br>

		<hr>
		Request: 
		<pre>%s</pre>
	</body>
	</html>
	`

	//values := []interface{}{r.RemoteAddr, r.RequestURI, r.URL.Path}
	values := []interface{}{r.RemoteAddr, r.URL.Path}

	html := fmt.Sprintf(htmlRaw, values...)

	_, _ = fmt.Fprintf(w, "%s", html)
}

func serveInfo(w http.ResponseWriter, _ *http.Request) {
	// find rice.Box
	templateBox, err := rice.FindBox("../../web")
	if err != nil {
		log.Fatal(err)
	}
	// get file contents as string
	templateString, err := templateBox.String("info.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	// parse and execute the template
	tmplMessage, err := template.New("info").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}

	if err := tmplMessage.Execute(w, metrics); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
