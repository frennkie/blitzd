package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
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
		m := data.NewMetricTimeBased(module, title)
		m.Interval = time.Duration(1 * time.Second).Seconds()

		uptime, err := host.Uptime()
		if err != nil {
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {
			m.Value = fmt.Sprintf("%d", uptime)
			m.Text = fmt.Sprintf("%ds", uptime)

			// update data in metric.Cache
			//log.Printf("Updating: %s.%s", module, title)
			data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.DefaultExpiration)

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}
