package system

import (
	"github.com/frennkie/blitzd/internal/data"
	"log"
	"runtime"
)

// SetOperatingSystem sets the "os" from "runtime.GOOS" and returns it as a "Metric"
func OperatingSystem() data.Metric {
	title := "os"
	log.Printf("setting: %s", title)

	metric := data.NewMetricStatic(title)
	metric.Value = runtime.GOOS

	return metric
}
