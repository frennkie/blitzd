package raspiblitz

import (
	"bufio"
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
	"strings"
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

	go raspiBlitzConfig()

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

	file, err := os.Open(absFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	config := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "=")
		if len(s) == 2 {
			config[s[0]] = s[1]
		}
	}

	log.Info(config)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	versionTitle := "version"
	version := data.NewMetricEventBased(module, versionTitle)
	version.Value = config["raspiBlitzVersion"]
	version.Text = config["raspiBlitzVersion"]
	data.Cache.Set(fmt.Sprintf("%s.%s", module, versionTitle), version, cache.NoExpiration)

	networkTitle := "network"
	network := data.NewMetricEventBased(module, networkTitle)
	network.Value = config["network"]
	network.Text = config["network"]
	data.Cache.Set(fmt.Sprintf("%s.%s", module, networkTitle), network, cache.NoExpiration)

	chainTitle := "chain"
	chain := data.NewMetricEventBased(module, chainTitle)
	chain.Value = config["chain"]
	chain.Text = config["chain"]
	data.Cache.Set(fmt.Sprintf("%s.%s", module, chainTitle), chain, cache.NoExpiration)

}
