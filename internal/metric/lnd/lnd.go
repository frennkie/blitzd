package lnd

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
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

	go foo()
	go foo2()

}

func foo() {

	title := "foo"
	log.Printf("starting goroutine: %s", title)

	for {
		foo := data.NewMetricTimeBased(title)
		foo.Value = "foo"
		foo.Interval = 2

		// update data in MetricCache
		log.Printf("Updating: %s", foo.Title)
		metric.LndMux.Lock()
		metric.Lnd.Foo = foo
		metric.LndMux.Unlock()

		time.Sleep(time.Duration(foo.Interval) * time.Second)
	}
}

// ToDo(frennkie) remove "foo2"
func foo2() {

	title := "foo2"
	log.Printf("starting goroutine: %s", title)

	for {
		foo := data.NewMetricTimeBased(title)
		foo.Value = "foo"
		foo.Interval = 12

		// update data in MetricCache
		log.Printf("Updating: %s", foo.Title)
		metric.LndMux.Lock()
		metric.Lnd.Foo = foo
		metric.LndMux.Unlock()

		time.Sleep(time.Duration(foo.Interval) * time.Second)
	}
}
