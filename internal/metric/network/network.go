package network

import "fmt"

func Init() {
	fmt.Println("network init called")

	go Nslookup()
	go Ping()
}
