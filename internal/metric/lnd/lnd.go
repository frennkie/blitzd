package lnd

import "fmt"

func init() {
	fmt.Println("lnd init called")

	go Foo()
}
