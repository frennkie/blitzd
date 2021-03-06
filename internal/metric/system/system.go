package system

import (
	"bufio"
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/internal/util"
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/patrickmn/go-cache"
	"github.com/shirou/gopsutil/host"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	module = "system"
)

var (
	logM = log.WithFields(log.Fields{"module": module})
)

func Init() {
	if config.C.Module.System.Enabled {
		logM.Info("starting")
	} else {
		logM.Warn("disabled by config - skipping")
		return
	}

	// set static
	operatingSystem()

	arch := Arch{Metric: data.NewMetricStatic(module, "arch")}
	metric.Set(arch)

	// start goroutines for event-based
	go etcIssue()
	go lsbRelease()

	// start goroutines for time-based

	mUptime := Uptime{Metric: data.NewMetricTimeBased(module, "uptimeNG")}
	mUptime.Metric.Interval = time.Duration(1 * time.Second).Seconds()
	go metric.Start(mUptime)

	ctx := context.Background()

	go func(ctx context.Context) {

		mLoad := Load{Metric: data.NewMetricTimeBased(module, "loadNG")}
		mLoad.Interval = time.Duration(5 * time.Second).Seconds()
		mLoad.Prefix = "Hallo, "
		mLoad.Suffix = " Konstantin!"

		for {
			m, _ := mLoad.generateValue(ctx)
			// mLoad.Text = fmt.Sprintf("%ds", uptime)

			// update Metric in Cache
			data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.DefaultExpiration)
			logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}

	}(ctx)

	go uptime()

	// config.C.Module.System.Metrics[0]

}

// SetOperatingSystem sets the "os" from "runtime.GOOS" and returns it as a "Metric"
func operatingSystem() {
	title := "os"

	m := data.NewMetricStatic(module, title)
	m.Value = runtime.GOOS
	m.Text = runtime.GOOS

	// update Metric in Cache
	data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.NoExpiration)
	logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Info("set")
}

func uptime() {
	title := "uptime"
	logM.WithFields(log.Fields{"title": title}).Info("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = time.Duration(1 * time.Second).Seconds()

		uptime, err := host.Uptime()
		if err != nil {
			logM.WithFields(log.Fields{"title": m.Title, "err": err}).Warn("error updating metric")
		} else {
			m.Value = fmt.Sprintf("%d", uptime)
			m.Text = fmt.Sprintf("%ds", uptime)

			// update Metric in Cache
			data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.DefaultExpiration)
			logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}

// TODO replace with /etc/issue
func lsbRelease() {
	absFilePath := "/etc/lsb-release"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		logM.WithFields(log.Fields{"file": absFilePath}).Warn("does not exist - skipping")
		return
	}

	logM.WithFields(log.Fields{"file": absFilePath}).Info("initial update")
	lsbReleaseFunc(absFilePath)

	logM.WithFields(log.Fields{"file": absFilePath}).Info("starting goroutine")
	go util.FileWatcher(absFilePath, lsbReleaseFunc)
}

func lsbReleaseFunc(absFilePath string) {
	logM.WithFields(log.Fields{"file": absFilePath, "kind": v1.Kind_KIND_EVENT_BASED}).Debug("update")

	title := "lsb_release"
	m := data.NewMetricEventBased(module, title)

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
	m.Text = tmp3

	data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.DefaultExpiration)
}

func etcIssue() {
	absFilePath := "/etc/issue"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		logM.WithFields(log.Fields{"file": absFilePath}).Warn("does not exist - skipping")
		return
	}

	logM.WithFields(log.Fields{"file": absFilePath}).Info("initial update")
	etcIssueFunc(absFilePath)

	logM.WithFields(log.Fields{"file": absFilePath}).Info("starting goroutine")
	go util.FileWatcher(absFilePath, etcIssueFunc)
}

func etcIssueFunc(absFilePath string) {
	logM.WithFields(log.Fields{"file": absFilePath, "kind": v1.Kind_KIND_EVENT_BASED}).Debug("update")

	title := "etc_issue"
	m := data.NewMetricEventBased(module, title)

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

	firstLine := txtlines[0]

	issue := strings.Replace(firstLine, "\\b", "", -1)
	issue = strings.Replace(issue, "\\d", "", -1)
	issue = strings.Replace(issue, "\\s", "", -1)
	issue = strings.Replace(issue, "\\l", "", -1)
	issue = strings.Replace(issue, "\\m", "", -1)
	issue = strings.Replace(issue, "\\n", "", -1)
	issue = strings.Replace(issue, "\\o", "", -1)
	issue = strings.Replace(issue, "\\r", "", -1)
	issue = strings.Replace(issue, "\\t", "", -1)
	issue = strings.Replace(issue, "\\u", "", -1)
	issue = strings.Replace(issue, "\\U", "", -1)
	issue = strings.Replace(issue, "\\v", "", -1)

	issue = strings.Trim(issue, " ")

	m.Value = issue
	m.Text = issue

	data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.DefaultExpiration)
}

type Arch struct{ v1.Metric }

func (m Arch) Set() {
	m.Value = runtime.GOARCH
	m.Text = runtime.GOARCH

	// update Metric in Cache
	data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m.Metric, cache.NoExpiration)
	logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Info("set")
}

type Uptime struct{ v1.Metric }

func (m Uptime) Start() {
	fmt.Println("from Uptime: start: ", m.Title)

	logM.WithFields(log.Fields{"title": m.Title}).Info("started goroutine")

	for {

		uptime, err := host.Uptime()
		if err != nil {
			logM.WithFields(log.Fields{"title": m.Title, "err": err}).Warn("error updating metric")
		} else {
			m.Value = fmt.Sprintf("%d", uptime)
			m.Text = fmt.Sprintf("%ds", uptime)

			// update Metric in Cache
			data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m.Metric, cache.DefaultExpiration)
			logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}

}

//Load
type Load struct {
	v1.Metric
}

func (m Load) generateValue(ctx context.Context) (v1.Metric, error) {
	mMetric := data.NewMetricTimeBased("system", "loadNNG")
	mMetric.Value = "honk"
	return mMetric, nil
}
