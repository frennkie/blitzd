package lnd

import (
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"log"
	"os"
	"time"
)

// ToDo(frennkie) remove "Foo"
func foo(c chan os.Signal) {

	title := "foo"
	log.Printf("starting goroutine: %s", title)

	for {

		time.Sleep(1 * time.Second)

		select {
		case sig := <-c:
			log.Printf("Got %s signal. Aborting...\n", sig)
			return
		default:
			log.Printf("test1")

			//Do something useful with message here
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
}
