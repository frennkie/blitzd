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
		Kind:         v1.Kind_KIND_UNSPECIFIED,
		Module:       module,
		Title:        title,
		Interval:     time.Duration(5 * time.Second).Seconds(),
		Timeout:      time.Duration(10 * time.Second).Seconds(),
		UpdatedAt:    updatedAt,
		ExpiredAfter: expiredAfter,
		Expired:      v1.Tribool_TRIBOOL_UNSPECIFIED,
		Value:        "",
		Prefix:       "",
		Suffix:       "",
		Style:        v1.Style_STYLE_UNSPECIFIED,
		Text:         "",
	}

	return metric
}

func NewMetricStatic(module, title string) v1.Metric {
	metric := NewMetric(module, title)
	metric.Kind = v1.Kind_KIND_STATIC
	metric.Interval = -1
	metric.ExpiredAfter = MaxTime
	metric.Style = v1.Style_STYLE_NORMAL
	return metric
}

func NewMetricTimeBased(module, title string) v1.Metric {
	metric := NewMetric(module, title)
	metric.Kind = v1.Kind_KIND_TIME_BASED
	metric.Style = v1.Style_STYLE_NORMAL
	return metric
}

func NewMetricEventBased(module, title string) v1.Metric {
	metric := NewMetric(module, title)
	metric.Kind = v1.Kind_KIND_EVENT_BASED
	metric.Interval = -1
	metric.ExpiredAfter = MaxTime
	metric.Style = v1.Style_STYLE_NORMAL
	return metric
}
