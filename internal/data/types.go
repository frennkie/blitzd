package data

import (
	v1 "github.com/frennkie/blitzd/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"time"
)

type Cache struct {
	Lnd     Lnd     `json:"lnd"`
	Network Network `json:"network"`
	System  System  `json:"system"`
}

type CacheOld struct {
	OperatingSystem Metric `json:"os"`
	Arch            Metric `json:"arch"`
	Foo             Metric `json:"foo"`
	Uptime          Metric `json:"uptime"`
	Nslookup        Metric `json:"nslookup"`
	Ping            Metric `json:"ping"`
	LsbRelease      Metric `json:"lsb-release"`
	FileBar         Metric `json:"file-bar"`
}

type Lnd struct {
	Foo Metric `json:"foo"`
}

type Network struct {
	Nslookup Metric `json:"nslookup"`
	Ping     Metric `json:"ping"`
}

type System struct {
	Arch            Metric `json:"arch"`
	OperatingSystem Metric `json:"os"`
	Uptime          Metric `json:"uptime"`
}

type Metric struct {
	Interval     float64   `json:"interval"`
	Timeout      float64   `json:"timeout"`
	Kind         Kind      `json:"kind"`
	Title        string    `json:"title"`
	Value        string    `json:"value"`
	Text         string    `json:"text"`
	Prefix       string    `json:"prefix"`
	Suffix       string    `json:"suffix"`
	Style        string    `json:"style"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiredAfter time.Time `json:"expired_after"`
	Expired      bool      `json:"expired"`
}

type Kind string

const (
	Static     Kind = "static"
	TimeBased  Kind = "time-based"
	EventBased Kind = "event-based"
)

var Kinds = map[string]Kind{
	"static":      Static,
	"time-based":  TimeBased,
	"event-based": EventBased,
}

func NewMetric(title string) Metric {
	metric := Metric{}

	metric.Title = title

	metric.Interval = time.Duration(5 * time.Second).Seconds()
	metric.Timeout = time.Duration(10 * time.Second).Seconds()

	metric.Value = "N/A"
	metric.Text = "N/A"
	metric.Prefix = ""
	metric.Suffix = ""
	metric.Style = Purple

	now := time.Now()
	metric.UpdatedAt = now
	metric.ExpiredAfter = now.Add(DefaultExpireTime)

	return metric
}

func NewMetricStatic(title string) Metric {
	metric := NewMetric(title)
	metric.Kind = Static
	metric.Interval = 0
	metric.Timeout = 0
	metric.ExpiredAfter = MaxTime
	return metric
}

func NewMetricTimeBased(title string) Metric {
	metric := NewMetric(title)
	metric.Kind = TimeBased
	return metric
}

func NewMetricEventBased(title string) Metric {
	metric := NewMetric(title)
	metric.Kind = EventBased
	metric.Interval = 0
	metric.ExpiredAfter = MaxTime
	return metric
}

type MetricTimeBasedNg interface {
	run()
}

func NewMetricNg(module, title string) v1.Metric {
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

func NewMetricNgStatic(module, title string) v1.Metric {
	metric := NewMetricNg(module, title)
	metric.Kind = v1.Kind_STATIC
	metric.Interval = -1
	metric.Timeout = 0
	metric.ExpiredAfter = MaxTimeNg
	return metric
}

func NewMetricNgTimeBased(module, title string) v1.Metric {
	metric := NewMetricNg(module, title)
	metric.Kind = v1.Kind_TIME_BASED
	return metric
}

func NewMetricNgEventBased(module, title string) v1.Metric {
	metric := NewMetricNg(module, title)
	metric.Kind = v1.Kind_EVENT_BASED
	metric.Interval = -1
	metric.ExpiredAfter = MaxTimeNg
	return metric
}
