package data

import "time"

type Cache struct {
	OperatingSystem Metric `json:"os"`
	Arch            Metric `json:"arch"`
	Foo             Metric `json:"foo"`
	Uptime          Metric `json:"uptime"`
	Nslookup        Metric `json:"nslookup"`
	Ping            Metric `json:"ping"`
	LsbRelease      Metric `json:"lsb-release"`
	FileBar         Metric `json:"file-bar,omitempty"`
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
