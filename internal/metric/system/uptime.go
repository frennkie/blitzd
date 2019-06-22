package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/patrickmn/go-cache"
	"github.com/shirou/gopsutil/host"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	module = "system"
)

func Uptime() {
	title := "uptime"

	logCtx := log.WithFields(log.Fields{"module": module, "title": title})
	logCtx.Debug("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = time.Duration(1 * time.Second).Seconds()

		uptime, err := host.Uptime()
		if err != nil {
			logCtx.WithFields(log.Fields{"err": err}).Warn("error updating metric")
		} else {
			m.Value = fmt.Sprintf("%d", uptime)
			m.Text = fmt.Sprintf("%ds", uptime)

			// update Metric in Cache
			data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.DefaultExpiration)
			logCtx.WithFields(log.Fields{"value": m.Value}).Trace("updated metric")

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}
