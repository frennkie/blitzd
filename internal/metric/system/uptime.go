package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/patrickmn/go-cache"
	"github.com/shirou/gopsutil/host"
	"log"
	"time"
)

func Uptime() {
	module := "system"
	title := "uptime"
	log.Printf("starting goroutine: %s.%s", module, title)
	for {
		m := data.NewMetricNgTimeBased(module, title)
		m.Interval = time.Duration(1 * time.Second).Seconds()

		uptime, err := host.Uptime()
		if err != nil {
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {
			m.Value = fmt.Sprintf("%d", uptime)

			// update data in MetricCache
			//log.Printf("Updating: %s", m.Title)
			metric.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.DefaultExpiration)

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}

func UptimeOld() {
	title := "uptime_old"
	log.Printf("starting goroutine: %s", title)
	for {
		m := data.NewMetricTimeBased(title)
		m.Interval = time.Duration(1 * time.Second).Seconds()

		uptime, err := host.Uptime()
		if err != nil {
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {
			m.Value = fmt.Sprintf("%d", uptime)

			// update data in MetricCache
			//log.Printf("Updating: %s", m.Title)

			metric.MetricsMux.Lock()
			metric.Metrics.System.Uptime = m
			metric.MetricsMux.Unlock()

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}
