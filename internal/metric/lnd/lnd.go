package lnd

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
	module = "lnd"
)

var (
	logM = log.WithFields(log.Fields{"module": module})
)

func Init() {
	if config.C.Module.Lnd.Enabled {
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

	go foo5()

}

// ToDo(frennkie) remove "foo5"
func foo5() {
	title := "foo5"
	logM.WithFields(log.Fields{"title": title}).Info("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = 12

		// gather and set data here
		m.Value = "foo5"
		m.Text = "foo5"

		// update Metric in Cache
		data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.NoExpiration)
		logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

		// sleep for Interval duration
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}
