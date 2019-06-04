package lnd

import (
	"fmt"
	"os"
	"os/signal"
)

func Init() {
	fmt.Println("lnd init called")

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	foo(c)

}
