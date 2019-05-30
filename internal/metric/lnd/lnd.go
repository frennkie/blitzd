package lnd

import "fmt"

func Init() {
	fmt.Println("lnd init called")

	go Foo()
}
