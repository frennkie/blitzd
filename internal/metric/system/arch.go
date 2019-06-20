package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/patrickmn/go-cache"
	"log"
	"runtime"
)

func Arch() {
	module := "system"
	title := "arch"
	log.Printf("setting: %s.%s", module, title)

	m := data.NewMetricStatic(module, title)
	m.Value = runtime.GOARCH
	m.Text = runtime.GOARCH
	data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.NoExpiration)
}
