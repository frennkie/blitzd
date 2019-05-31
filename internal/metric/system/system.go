package system

import (
	"fmt"
	"github.com/frennkie/blitzd/internal/metric"
)

func Init() {
	fmt.Println("system init called")

	// set static MetricsOld
	metric.MetricsMux.Lock()
	metric.Metrics.System.Arch = Arch()
	metric.Metrics.System.OperatingSystem = OperatingSystem()
	metric.MetricsMux.Unlock()

	// start goroutine for event-based

	// start goroutine for time-based
	go Uptime()
}
