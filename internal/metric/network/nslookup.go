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

func Nslookup() {
	title := "nslookup"
	log.Printf("starting goroutine: %s", title)
	for {
		m := data.NewMetric(title)
		m.Interval = time.Duration(60 * time.Second).Seconds()

		cmdName := "nslookup"
		var cmdArgs []string

		if runtime.GOOS == "windows" {
			cmdArgs = []string{"google.com"}
		} else {
			cmdArgs = []string{"google.com"}
		}

		result, err := util.TimeoutExec(cmdName, cmdArgs)
		if err != nil {
			log.Printf("Error Updating: %s - %s", m.Title, err)
		} else {
			split := strings.Split(result, util.GetNewLine())
			last := strings.TrimSpace(split[len(split)-3])

			m.Value = last

			// update data in MetricCache
			log.Printf("Updating: %s", m.Title)
			metric.NetworkMux.Lock()
			metric.Network.Nslookup = m
			metric.NetworkMux.Unlock()

		}
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}
