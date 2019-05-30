package network

import "fmt"

func init() {
	fmt.Println("network init called")

	go Nslookup()
	go Ping()
}
