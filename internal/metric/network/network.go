package network

import (
	"context"
	"github.com/frennkie/blitzd/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

const (
	module = "network"
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

	go Nslookup()
	go Ping()
}
