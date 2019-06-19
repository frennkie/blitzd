package system

import (
	"bufio"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"github.com/frennkie/blitzd/internal/util"
	"log"
	"os"
	"strings"
)

// TODO replace with /etc/issue
func UpdateLsbRelease() {
	title := "lsb-release"
	absFilePath := "/etc/lsb-release"

	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		log.Printf("file does not exist - skipping: %s", absFilePath)
		return
	}

	log.Printf("starting goroutine: %s (%s)", title, absFilePath)

	UpdateLsbReleaseFunc(title, absFilePath)
	go util.FileWatcher(title, absFilePath, UpdateLsbReleaseFunc)
}

func UpdateLsbReleaseFunc(title string, absFilePath string) {
	log.Printf("event-based update: %s (%s)", title, absFilePath)
	m := data.NewMetricEventBased(title)

	file, err := os.Open(absFilePath)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	_ = file.Close()

	tmp := txtlines[len(txtlines)-1]
	tmp2 := strings.Split(tmp, "=")[1]
	tmp3 := strings.Replace(tmp2, "\"", "", -1)
	m.Value = tmp3

	metric.MetricsOldMux.Lock()
	metric.MetricsOld.LsbRelease = m
	metric.MetricsOldMux.Unlock()

}
