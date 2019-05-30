package system

import (
	"github.com/frennkie/blitzd/internal/data"
	"log"
	"runtime"
)

func Arch() data.Metric {
	title := "arch"
	log.Printf("setting: %s", title)

	metric := data.NewMetricStatic(title)
	metric.Value = runtime.GOARCH

	return metric
}
