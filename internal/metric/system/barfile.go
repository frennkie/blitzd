package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/util"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
)

func UpdateFileBar() {
	module := "system"
	title := "file-bar"
	absFilePath := "/tmp/foo"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateFileBarFunc(module, title, absFilePath)
	go util.FileWatcher(module, title, absFilePath, UpdateFileBarFunc)
}

func UpdateFileBarFunc(module, title string, absFilePath string) {
	log.Printf("event-based update: %s.%s (%s)", module, title, absFilePath)
	m := data.NewMetricEventBased(module, title)

	m.Value = fmt.Sprintf("%s", "foobar")
	m.Text = fmt.Sprintf("%s", "foobar")

	data.Cache.Set(fmt.Sprintf("%s.%s", module, title), m, cache.DefaultExpiration)

}
