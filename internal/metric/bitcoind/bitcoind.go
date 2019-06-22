package bitcoind

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

const (
	module = "bitcoind"
)

var (
	logM = log.WithFields(log.Fields{"module": module})
)

func Init() {
	if config.C.Module.Bitcoind.Enabled {
		logM.Info("starting")
	} else {
		logM.Warn("disabled by config - skipping")
		return
	}

	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	go foo6()

}

// ToDo(frennkie) remove "foo6"
func foo6() {
	title := "foo6"
	logM.WithFields(log.Fields{"title": title}).Info("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = 6

		// gather and set data here
		m.Value = "foo6"
		m.Text = "foo6"

		// update Metric in Cache
		data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.NoExpiration)
		logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

		// sleep for Interval duration
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}
