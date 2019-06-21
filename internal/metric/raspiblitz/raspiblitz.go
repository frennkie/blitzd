package raspiblitz

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

const (
	module = "raspiblitz"
)

func Init() {
	if viper.GetBool("module.raspiblitz.enabled") {
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

	go version()
	go raspiBlitzConfig()

}

// ToDo(frennkie) remove "foo5"
func version() {
	title := "version"
	log.WithFields(log.Fields{"module": module, "title": title}).Debug("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = 12

		// gather and set data here
		m.Value = "0.1.2"
		m.Text = "v0.1.2"

		// update Metric in Cache
		data.Cache.Set(title, m, cache.NoExpiration)
		log.WithFields(log.Fields{"module": module, "title": title, "value": m.Value}).Trace("updated metric")

		// sleep for Interval duration
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}

func raspiBlitzConfig() {
	absFilePath := viper.GetString("module.raspiblitz.path")

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", module, absFilePath)

	raspiBlitzConfigFunc(absFilePath)
	go util.FileWatcher(absFilePath, raspiBlitzConfigFunc)
}

func raspiBlitzConfigFunc(absFilePath string) {
	log.WithFields(log.Fields{"absFilePath": absFilePath, "kind": v1.Kind_EVENT_BASED, "module": module}).Debug("update")

	versionTitle := "version"
	version := data.NewMetricEventBased(module, versionTitle)
	version.Value = fmt.Sprintf("%s", "foobar")
	version.Text = fmt.Sprintf("%s", "foobar")
	data.Cache.Set(fmt.Sprintf("%s.%s", module, versionTitle), version, cache.NoExpiration)

}
