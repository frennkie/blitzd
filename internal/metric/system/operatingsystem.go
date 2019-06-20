package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/patrickmn/go-cache"
	"log"
	"runtime"
)

// SetOperatingSystem sets the "os" from "runtime.GOOS" and returns it as a "Metric"
func OperatingSystem() {
	module := "system"
	title := "os"
	log.Printf("setting: %s.%s", module, title)

	m := data.NewMetricStatic(module, title)
	m.Value = runtime.GOOS
	m.Text = runtime.GOOS
	data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.NoExpiration)
}
