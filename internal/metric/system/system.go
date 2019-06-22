package system

import (
	"fmt"
	"runtime"
)

func Init() {
	fmt.Println("system init called")

	// set static
	Arch()
	OperatingSystem()

	// start goroutine for event-based

	// start goroutine for time-based
	go Uptime()

	// ToDo(frennkie) sort out
	if runtime.GOOS != "windows" {
		go UpdateLsbRelease()
	}
}
