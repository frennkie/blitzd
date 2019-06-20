package network

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/patrickmn/go-cache"
	"log"
	"runtime"
	"strings"
	"time"
)

func Ping() {
	module := "network"
	title := "ping"
	log.Printf("starting goroutine: %s.%s", module, title)
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
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {

			split := strings.Split(result, util.GetNewLine())
			last := strings.TrimSpace(split[len(split)-2])

			m.Value = last
			m.Text = last

			// update data in MetricCache
			log.Printf("Updating: %s.%s", module, title)
			data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.DefaultExpiration)

			time.Sleep(time.Duration(m.Interval) * time.Second)
		}
	}
}
