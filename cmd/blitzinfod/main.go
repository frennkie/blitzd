package main

import (
	"bufio"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/frennkie/blitzinfod/internal/serve"
	"github.com/frennkie/blitzinfod/internal/utils"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
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

	Normal = "normal"
	Red    = "red"
	Green  = "green"
	Yellow = "yellow"
	Purple = "purple"
)

var (
	metrics    data.Cache
	metricsMux sync.Mutex

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
	metric.ExpiredAfter = now.Add(data.DefaultExpireTime)

	return metric
}

func NewMetricStatic(title string) data.Metric {
	metric := NewMetric(title)
	metric.Kind = data.Static
	metric.Interval = 0
	metric.Timeout = 0
	metric.ExpiredAfter = data.MaxTime
	return metric
}

func NewMetricTimeBased(title string) data.Metric {
	metric := NewMetric(title)
	metric.Kind = data.TimeBased
	return metric
}

func NewMetricEventBased(title string) data.Metric {
	metric := NewMetric(title)
	metric.Kind = data.EventBased
	metric.Interval = 0
	metric.ExpiredAfter = data.MaxTime
	return metric
}

// SetOperatingSystem sets the "os" from "runtime.GOOS" and returns it as a "Metric"
func SetOperatingSystem() data.Metric {
	title := "os"
	log.Printf("setting: %s", title)

	metric := NewMetricStatic(title)
	metric.Value = runtime.GOOS

	return metric
}

func SetArch() data.Metric {
	title := "arch"
	log.Printf("setting: %s", title)

	metric := NewMetricStatic(title)
	metric.Value = runtime.GOARCH

	return metric
}

// ToDo(frennkie) remove "Foo"
func UpdateFoo() {
	title := "foo"
	log.Printf("starting goroutine: %s", title)

	for {
		foo := NewMetric(title)
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
	title := "uptime"
	log.Printf("starting goroutine: %s", title)
	for {
		m := NewMetricTimeBased(title)
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
	title := "nslookup"
	log.Printf("starting goroutine: %s", title)
	for {
		m := NewMetric(title)
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
			split := strings.Split(result, utils.GetNewLine())
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
	title := "ping"
	log.Printf("starting goroutine: %s", title)
	for {
		m := NewMetric(title)
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

			split := strings.Split(result, utils.GetNewLine())
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

func UpdateFileBar() {
	title := "file-bar"
	absFilePath := "/tmp/foo"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateFileBarFunc(title, absFilePath)
	go utils.FileWatcher(title, absFilePath, UpdateFileBarFunc)
}

func UpdateFileBarFunc(title string, absFilePath string) {
	log.Printf("event-based update: %s (%s)", title, absFilePath)
	m := NewMetricEventBased(title)

	m.Value = fmt.Sprintf("%s", "foobar")

	metricsMux.Lock()
	metrics.FileBar = m
	metricsMux.Unlock()

}

func UpdateLsbRelease() {
	title := "lsb-release"
	absFilePath := "/etc/lsb-release"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateLsbReleaseFunc(title, absFilePath)
	go utils.FileWatcher(title, absFilePath, UpdateLsbReleaseFunc)
}

func UpdateLsbReleaseFunc(title string, absFilePath string) {
	log.Printf("event-based update: %s (%s)", title, absFilePath)
	m := NewMetricEventBased(title)

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

	metricsMux.Lock()
	metrics.LsbRelease = m
	metricsMux.Unlock()

}

func main() {
	var rootCmd = &cobra.Command{
		Version: buildVersion,
		Use:     "blitzinfod",
		Short:   "RaspiBlitz Info Daemon",
		Long: `A service that retrieves and caches details about your RaspiBlitz.
                More info at: https://github.com/frennkie/blitzinfod`,
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

	if runtime.GOOS != "windows" {
		go UpdateLsbRelease()
		go UpdateFileBar()
	}

	box := rice.MustFindBox("../../web/")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))

	http.HandleFunc("/", serve.Root)
	http.HandleFunc("/info/", serve.Info(&metrics))
	http.HandleFunc(data.APIv1, serve.StaticApi(&metrics))

	RESTHostPort := viper.GetString("RESTHostPort")
	log.Printf("REST: Listening on host: http://%s", RESTHostPort)

	//rpcHostPort := viper.GetString("rpcHostPort")
	//log.Printf("RPC: Listening on host: gRPC://%s", rpcHostPort)

	// now ListenAndServer
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", RESTHostPort), nil))

}
