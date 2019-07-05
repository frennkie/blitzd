package metric

import (
	"context"
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

type SetStatic interface {
	Set()
}

func Set(a SetStatic) { a.Set() }

type StartTimeBased interface {
	Start()
}

func Start(a StartTimeBased) { a.Start() }

//type StartTimeBasedNg interface {
//	DoStartLoad()
//	DoUpdateValueLoad() (string, error)
//	DoFormatTextLoad() string
//}

func DoStartLoad(m v1.Metric) {
	logM := log.WithFields(log.Fields{"module": m.Module})
	logM.WithFields(log.Fields{"title": m.Title}).Info("started goroutine")

	for {

		//value, err := m.Foo
		//if err != nil {
		//	log.WithError(err).Warnf("an error occurred")
		//}
		//m.Value = value
		m.Value = "foo"
		//m.Text = m.Foo

		// update Metric in Cache
		data.Cache.Set(fmt.Sprintf("%s.%s", m.Module, m.Title), m, cache.DefaultExpiration)
		logM.WithFields(log.Fields{"title": m.Title, "value": m.Value}).Trace("updated metric")

		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}

//
//func DoFormatTextLoad(s StartTimeBasedNg) string {
//	return "honk"
//	//return fmt.Sprintf("%s%s%s", s.Prefix, s.Value, s.Suffix)
//}

type UpdateEventBased interface {
	Update()
}

func Update(a UpdateEventBased) { a.Update() }

type Blitz struct {
	v1.Metric
}

func (b *Blitz) Sleep() {
	for {
		fmt.Println("sleeping...zzz")
		time.Sleep(2 * time.Second)
	}

}

type ValueGenerator interface {
	generateValue(ctx context.Context) (v1.Metric, error)
}

type GroupValueGenerator interface {
	generateGroupValue(ctx context.Context) ([]v1.Metric, error)
}
