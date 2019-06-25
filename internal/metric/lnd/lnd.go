package lnd

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"os/user"
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
	go version()

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

func version() {
	title := "version"
	var interval float64 = 8
	logM.WithFields(log.Fields{"title": title}).Info("started goroutine")

	// ToDo(frennkie) check and transfer this
	errored := false
	for {
		if errored {
			// Last execution ran into an error and did not sleep at end of
			// for loop. Therefore sleep now for interval.
			time.Sleep(time.Duration(interval) * time.Second)
			errored = false
		}

		m := data.NewMetricTimeBased(module, title)
		m.Interval = interval

		usr, err := user.Current()
		if err != nil {
			logM.WithFields(log.Fields{"title": title, "err": err}).Error("Cannot get current user:")
			errored = true
			continue
		}
		logM.WithFields(log.Fields{"title": title, "home": usr.HomeDir}).Info("The user home directory: " + usr.HomeDir)
		tlsCertPath := config.C.Module.Lnd.TlsCert
		macaroonPath := config.C.Module.Lnd.Macaroon

		tlsCreds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
		if err != nil {
			logM.WithFields(log.Fields{"title": title, "err": err}).Error("Cannot get node tls credentials")
			errored = true
			continue
		}

		macaroonBytes, err := ioutil.ReadFile(macaroonPath)
		if err != nil {
			logM.WithFields(log.Fields{"title": title, "err": err}).Error("Cannot read macaroon file")
			errored = true
			continue
		}

		mac := &macaroon.Macaroon{}
		if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
			logM.WithFields(log.Fields{"title": title, "err": err}).Error("Cannot unmarshal macaroon")
			errored = true
			continue
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(tlsCreds),
			grpc.WithBlock(),
			grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(mac)),
		}

		conn, err := grpc.Dial("localhost:10009", opts...)
		if err != nil {
			logM.WithFields(log.Fields{"title": title, "err": err}).Error("cannot dial to lnd")
			errored = true
			continue
		}
		client := lnrpc.NewLightningClient(conn)

		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		getInfoResp, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
		if err != nil {
			logM.WithFields(log.Fields{"title": title, "err": err}).Error("Cannot get info from node:")
			errored = true
			continue
		}

		//// ToDo(frennkie) was defer cancel() .. but IDE complains about defer in for loops
		//defer cancel()
		//spew.Dump(getInfoResp)

		// gather and set data here
		m.Value = getInfoResp.Version
		m.Text = getInfoResp.Version

		// update Metric in Cache
		data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.NoExpiration)
		logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

		// sleep for Interval duration
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}
