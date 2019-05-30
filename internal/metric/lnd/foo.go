package lnd

import (
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"log"
	"time"
)

// ToDo(frennkie) remove "Foo"
func Foo() {
	title := "foo"
	log.Printf("starting goroutine: %s", title)

	for {
		foo := data.NewMetricTimeBased(title)
		foo.Value = "foo"

		// update data in MetricCache
		log.Printf("Updating: %s", foo.Title)
		metric.LndMux.Lock()
		metric.Lnd.Foo = foo
		metric.LndMux.Unlock()

		time.Sleep(time.Duration(foo.Interval) * time.Second)
	}
}
