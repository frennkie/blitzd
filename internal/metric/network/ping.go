package network

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"time"
)

func ping() {
	title := "ping"
	logM.WithFields(log.Fields{"title": title}).Info("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = time.Duration(60 * time.Second).Seconds()

		cmdName := "ping"
		var cmdArgs []string

		if runtime.GOOS == "windows" {
			cmdArgs = []string{"-n", "2", "8.8.8.8"}
		} else {
			cmdArgs = []string{"-c", "2", "8.8.8.8"}
		}

		result, err := util.TimeoutExec(cmdName, cmdArgs)
		if err != nil {
			logM.WithFields(log.Fields{"title": m.Title, "err": err}).Warn("error updating metric")
		} else {

			split := strings.Split(result, util.GetNewLine())
			last := strings.TrimSpace(split[len(split)-2])

			m.Value = last
			m.Text = last

			// update data in MetricCache
			data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.DefaultExpiration)
			logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}
