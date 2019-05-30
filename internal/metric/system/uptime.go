package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/shirou/gopsutil/host"
	"log"
	"time"
)

func Uptime() {
	title := "uptime"
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
			metric.SystemMux.Lock()
			metric.System.Uptime = m
			metric.SystemMux.Unlock()

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}
