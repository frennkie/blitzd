package lnd

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
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

	go foo()
	go foo2()
	go foo3.run(foo3{})
	go foo4()
	go foo5()

}

func foo() {

	title := "foo"
	log.Printf("starting goroutine: %s", title)

	for {
		foo := data.NewMetricTimeBased(title)
		foo.Value = "foo"
		foo.Interval = 4

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
		foo2 := data.NewMetricTimeBased(title)
		foo2.Value = "foo2"
		foo2.Interval = 12

		// update data in MetricCache
		//log.Printf("Updating: %s", foo2.Title)
		metric.LndMux.Lock()
		metric.Lnd.Foo = foo2
		metric.LndMux.Unlock()

		time.Sleep(time.Duration(foo2.Interval) * time.Second)
	}
}

type foo3 struct {
	data.Metric
}

func (foo3 foo3) run() {

	title := "foo3"
	foo3.Title = "foo3"
	log.Printf("starting goroutine: %s", foo3.Title)

	for {
		foo3 := data.NewMetricTimeBased(title)
		foo3.Value = "foo3"
		foo3.Interval = 6

		// update data in MetricCache
		//log.Printf("Updating: %s", foo3.Title)
		metric.LndMux.Lock()
		metric.Lnd.Foo = foo3
		metric.LndMux.Unlock()

		time.Sleep(time.Duration(foo3.Interval) * time.Second)
	}

}

// ToDo(frennkie) remove "foo4"
func foo4() {

	title := "foo4"
	log.Printf("starting goroutine: %s", title)

	for {

		//foo4 := data.NewMetricTimeBased(title)
		foo4 := v1.Metric{}
		foo4.Value = "foo4"
		foo4.Interval = 12

		foo4.Interval = time.Duration(5 * time.Second).Seconds()
		foo4.Timeout = time.Duration(10 * time.Second).Seconds()

		now := time.Now()
		foo4.UpdatedAt, _ = ptypes.TimestampProto(now)
		foo4.ExpiredAfter, _ = ptypes.TimestampProto(now.Add(data.DefaultExpireTime))

		// update data in MetricCache
		//log.Printf("Updating: %s", foo4.Title)

		//metric.MetricsAPIMux.Lock()
		//metric.MetricsAPI = append(metric.MetricsAPI, foo4)
		//metric.MetricsAPIMux.Unlock()
		metric.MetricsFoo4 = foo4

		time.Sleep(time.Duration(foo4.Interval) * time.Second)
	}
}

// ToDo(frennkie) remove "foo4"
func foo5() {

	title := "lnd.foo5"
	log.Printf("starting goroutine: %s", title)

	for {

		//foo5 := data.NewMetricTimeBased(title)
		foo5 := v1.Metric{Title: title}
		foo5.Value = "foo5"
		foo5.Text = "foo5"
		foo5.Interval = 12
		foo5.Kind = v1.Kind_TIME_BASED

		foo5.Interval = time.Duration(5 * time.Second).Seconds()
		foo5.Timeout = time.Duration(10 * time.Second).Seconds()

		now := time.Now()
		foo5.UpdatedAt, _ = ptypes.TimestampProto(now)
		foo5.ExpiredAfter, _ = ptypes.TimestampProto(now.Add(4 * time.Second))

		// update data in MetricCache
		//log.Printf("Updating: %s", foo5.Title)

		//metric.MetricsAPIMux.Lock()
		//metric.MetricsAPI = append(metric.MetricsAPI, foo5)
		//metric.MetricsAPIMux.Unlock()
		metric.Cache.Set(title, foo5, cache.NoExpiration)

		time.Sleep(time.Duration(foo5.Interval) * time.Second)
	}
}
