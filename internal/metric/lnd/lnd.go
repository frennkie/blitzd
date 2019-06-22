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

func Init() {
	if config.C.Module.Lnd.Enabled {
		log.WithFields(log.Fields{"module": module}).Info("starting module")
	} else {
		log.WithFields(log.Fields{"module": module}).Info("skipping module - disabled by config")
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
	log.WithFields(log.Fields{"module": module, "title": title}).Debug("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = 12

		// gather and set data here
		m.Value = "foo5"
		m.Text = "foo5"

		// update Metric in Cache
		data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.NoExpiration)
		log.WithFields(log.Fields{"module": module, "title": title, "value": m.Value}).Trace("updated metric")

		// sleep for Interval duration
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}
