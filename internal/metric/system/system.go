package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/metric"
)

func init() {
	fmt.Println("system init called")

	// set static Metrics
	metric.System.Arch = Arch()
	metric.System.OperatingSystem = OperatingSystem()

	// start goroutine for event-based

	// start goroutine for time-based
	go Uptime()
}
