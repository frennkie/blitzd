package data

import "time"

type Cache struct {
	OperatingSystem Metric `json:"os"`
	Arch            Metric `json:"arch"`
	Foo             Metric `json:"foo"`
	Uptime          Metric `json:"uptime"`
}

type Metric struct {
	Interval     float64   `json:"interval"`
	Timeout      float64   `json:"timeout"`
	Title        string    `json:"title"`
	Value        string    `json:"value"`
	Text         string    `json:"text"`
	Prefix       string    `json:"prefix"`
	Suffix       string    `json:"suffix"`
	Style        string    `json:"style"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiredAfter time.Time `json:"expired_after"`
}
