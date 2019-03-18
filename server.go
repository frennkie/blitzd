package main

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/host"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
)

const (
	defaultHost = "localhost"
	defaultPort = 7080

	defaultExpireTime = 300 * time.Second // 5 minutes

	Normal = "normal"
	Red    = "red"
	Green  = "green"
	Yellow = "yellow"
	Purple = "purple"
)

var (
	metrics    Cache
	metricsMux sync.Mutex

	// maxTime (Metric does not expire): "3000-01-01T00:00:00Z"
	maxTime = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
)

type Cache struct {
	OperatingSystem Metric `json:"os"`
	Arch            Metric `json:"arch"`
	Foo             Metric `json:"foo"`
	Uptime          Metric `json:"uptime"`
}

type Metric struct {
	Interval     float64   `json:"interval"`
	Timeout      float64   `json:"timeout"`
	Title        string    `json:"title"`
	Value        string    `json:"value"`
	Text         string    `json:"text"`
	Prefix       string    `json:"prefix"`
	Suffix       string    `json:"suffix"`
	Style        string    `json:"style"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiredAfter time.Time `json:"expired_after"`
}

func NewMetric(title string) Metric {
	metric := Metric{}

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
func SetOperatingSystem() Metric {
	metric := NewMetric("os")

	metric.Interval = 0
	metric.Timeout = 0
	metric.ExpiredAfter = maxTime

	metric.Value = runtime.GOOS

	return metric
}

func SetArch() Metric {
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

// config defines the configuration options for lnd.
//
// See loadConfig for further details regarding the configuration
// loading+parsing process.
type config struct {
	ShowVersion bool `short:"V" long:"version" description:"Display version information and exit"`

	Host string `short:"H" long:"host" description:"The host to listen on"`
	Port uint16 `short:"P" long:"port" description:"The port to listen on"`
}

func loadConfig() (*config, error) {
	defaultCfg := config{
		Host: defaultHost,
		Port: defaultPort,
	}

	// Pre-parse the command line options to pick up an alternative config
	// file.
	preCfg := defaultCfg
	if _, err := flags.Parse(&preCfg); err != nil {
		return nil, err
	}

	return &preCfg, nil

}

func main() {

	// Load the configuration, and parse any command line options. This
	// function will also set up logging properly.
	loadedConfig, err := loadConfig()
	if err != nil {
		return
	}
	cfg := loadedConfig

	// set static Metrics
	metrics.Arch = SetArch()
	metrics.OperatingSystem = SetOperatingSystem()

	// start Update of Metrics as goroutines
	go UpdateFoo()
	go UpdateUptime()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/info/", serveInfo)
	http.HandleFunc("/api/", serveStaticApi)

	log.Printf("Listening on host: http://%s:%d", cfg.Host, cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))

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

func serveInfo(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "view.html")

	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "view.html", metrics); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
