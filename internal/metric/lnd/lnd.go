package lnd

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"os/signal"
	"time"
)

func Init() {
	fmt.Println("lnd init called")

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
	module := "lnd"
	title := "foo5"
	log.Printf("starting goroutine: %s.%s", module, title)

	for {
		foo5 := data.NewMetricTimeBased(module, title)
		foo5.Value = "foo5"
		foo5.Text = "foo5"
		foo5.Interval = 12

		// update data in metric.Cache
		//log.Printf("Updating: %s.%s", module, title)
		data.Cache.Set(title, foo5, cache.NoExpiration)

		time.Sleep(time.Duration(foo5.Interval) * time.Second)
	}
}
