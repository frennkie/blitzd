package data

import (
	"github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	DefaultExpireTime = 300 * time.Second // 5 minutes

)

var (
	// Global Cache
	Cache = cache.New(5*time.Minute, 10*time.Minute)

	// maxTime (Metric does not expire): "3000-01-01T00:00:00Z"
	MaxTime, _ = ptypes.TimestampProto(time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC))
)

func NewMetric(module, title string) v1.Metric {
	now := time.Now()
	updatedAt, _ := ptypes.TimestampProto(now)
	expiredAfter, _ := ptypes.TimestampProto(now.Add(DefaultExpireTime))

	metric := v1.Metric{
		Kind:         0,
		Module:       module,
		Title:        title,
		Interval:     time.Duration(5 * time.Second).Seconds(),
		Timeout:      time.Duration(10 * time.Second).Seconds(),
		UpdatedAt:    updatedAt,
		ExpiredAfter: expiredAfter,
		Expired:      false,
		Value:        "",
		Prefix:       "",
		Suffix:       "",
		Style:        0,
		Text:         "",
	}

	return metric
}

func NewMetricStatic(module, title string) v1.Metric {
	metric := NewMetric(module, title)
	metric.Kind = v1.Kind_STATIC
	metric.Interval = -1
	metric.Timeout = 0
	metric.ExpiredAfter = MaxTime
	return metric
}

func NewMetricTimeBased(module, title string) v1.Metric {
	metric := NewMetric(module, title)
	metric.Kind = v1.Kind_TIME_BASED
	return metric
}

func NewMetricEventBased(module, title string) v1.Metric {
	metric := NewMetric(module, title)
	metric.Kind = v1.Kind_EVENT_BASED
	metric.Interval = -1
	metric.ExpiredAfter = MaxTime
	return metric
}
