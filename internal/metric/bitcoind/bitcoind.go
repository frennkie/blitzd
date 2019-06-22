package bitcoind

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
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
	go blockCount()

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

func blockCount() {
	title := "block_count"
	logM.Info("started goroutine")

	// Connect to local bitcoind core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         config.C.Module.Bitcoind.RpcAddress,
		User:         config.C.Module.Bitcoind.RpcUser,
		Pass:         config.C.Module.Bitcoind.RpcPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = 60

		// connect
		client, err := rpcclient.New(connCfg, nil)
		if err != nil {
			logM.WithFields(log.Fields{"title": m.Title, "err": err}).Error("failed to create client")
		}

		// Get the current block count.
		blockCount, err := client.GetBlockCount()
		if err != nil {
			logM.WithFields(log.Fields{"title": m.Title, "err": err}).Error("client failed")
		}
		m.Value = fmt.Sprintf("%d", blockCount)
		m.Text = fmt.Sprintf("%d", blockCount)

		// update Metric in Cache
		data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.NoExpiration)
		logM.WithFields(log.Fields{"title": m.Title, "value": blockCount}).Trace("updated metric")

		// close client connection
		client.Shutdown()

		// sleep for Interval duration
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}

}
