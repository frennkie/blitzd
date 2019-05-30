package network

import (
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/internal/util"
	"log"
	"runtime"
	"strings"
	"time"
)

func Ping() {
	title := "ping"
	log.Printf("starting goroutine: %s", title)
	for {
		m := data.NewMetric(title)
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

			// update data in MetricCache
			log.Printf("Updating: %s", m.Title)
			metric.NetworkMux.Lock()
			metric.Network.Ping = m
			metric.NetworkMux.Unlock()

		}
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}
