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
	logM.WithFields(log.Fields{"title": title}).Info("started goroutine")

	for {
		m := data.NewMetricTimeBased(module, title)
		m.Interval = 12

		usr, err := user.Current()
		if err != nil {
			fmt.Println("Cannot get current user:", err)
			return
		}
		fmt.Println("The user home directory: " + usr.HomeDir)
		tlsCertPath := config.C.Module.Lnd.TlsCert
		macaroonPath := config.C.Module.Lnd.Macaroon

		tlsCreds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
		if err != nil {
			fmt.Println("Cannot get node tls credentials", err)
			return
		}

		macaroonBytes, err := ioutil.ReadFile(macaroonPath)
		if err != nil {
			fmt.Println("Cannot read macaroon file", err)
			return
		}

		mac := &macaroon.Macaroon{}
		if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
			fmt.Println("Cannot unmarshal macaroon", err)
			return
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(tlsCreds),
			grpc.WithBlock(),
			grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(mac)),
		}

		conn, err := grpc.Dial("localhost:10009", opts...)
		if err != nil {
			fmt.Println("cannot dial to lnd", err)
			return
		}
		client := lnrpc.NewLightningClient(conn)

		ctx := context.Background()
		getInfoResp, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
		if err != nil {
			fmt.Println("Cannot get info from node:", err)
			return
		}
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
