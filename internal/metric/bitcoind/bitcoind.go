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
	go bitcoindRpc()

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

func bitcoindRpc() {
	title := "bitcoin"
	logM.WithFields(log.Fields{"title": title}).Info("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = 60

		// gather and set data here
		// create new client instance
		client, err := rpcclient.New(&rpcclient.ConnConfig{
			HTTPPostMode: true,
			DisableTLS:   true,
			Host:         config.C.Module.Bitcoind.RpcAddress,
			User:         config.C.Module.Bitcoind.RpcUser,
			Pass:         config.C.Module.Bitcoind.RpcPassword,
		}, nil)
		if err != nil {
			log.Fatalf("error creating new btc client: %v", err)
		}

		// list accounts
		accounts, err := client.ListAccounts()
		if err != nil {
			log.Fatalf("error listing accounts: %v", err)
		}
		// iterate over accounts (map[string]btcutil.Amount) and write to stdout
		for label, amount := range accounts {
			log.Printf("%s: %s", label, amount)
		}

		// update Metric in Cache
		data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.NoExpiration)
		logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

		// sleep for Interval duration
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}

}
