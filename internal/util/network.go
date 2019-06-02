package util

import (
	"bytes"
	"net"
	"os"
	"sort"
)

// getIPAddresses returns a slice of net.IP that contains all IPs of this system.
func getIPAddresses() []net.IP {

	ipAddresses := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}

	// addIP appends an IP address only if it isn't already in the slice.
	addIP := func(ipAddr net.IP) {
		for _, ip := range ipAddresses {
			if bytes.Equal(ip, ipAddr) {
				return
			}
		}
		ipAddresses = append(ipAddresses, ipAddr)
	}

	// Add all the interface IPs that aren't already in the slice.
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, a := range addrs {
		ipAddr, _, err := net.ParseCIDR(a.String())
		if err == nil {
			addIP(ipAddr)
		}
	}

	return ipAddresses

}

func GetLocalAddresses() []string {

	var addresses []string

	for _, ip := range getIPAddresses() {
		// skip if this is not a loopback address
		if ip.IsLoopback() {
			res := ip.To4()
			if res == nil {
				addresses = append(addresses, "["+ip.String()+"]")
			} else {
				addresses = append(addresses, ip.String())
			}
		}
	}

	addresses = append(addresses, "localhost")

	return addresses

}

func GetNonLocalAddresses() []string {

	var addresses []string

	for _, ip := range getIPAddresses() {
		// skip if this is a loopback address and exclude is set
		if ip.IsLoopback() {
			continue
		} else {
			res := ip.To4()
			if res == nil {
				addresses = append(addresses, "["+ip.String()+"]")
			} else {
				addresses = append(addresses, ip.String())
			}
		}
	}

	// Collect the hostname's names into a slice.
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// do not add "localhost"
	if hostname != "localhost" {
		addresses = append(addresses, hostname)
	}

	return addresses

}

func GetBaseUrls(scheme, portStr string, includeLocal, includeNonLcal bool) []string {
	var baseUrls []string

	if includeLocal {
		for _, addr := range GetLocalAddresses() {
			baseUrls = append(baseUrls, scheme+"://"+addr+":"+portStr+"/")
		}
	}

	if includeNonLcal {
		for _, addr := range GetNonLocalAddresses() {
			baseUrls = append(baseUrls, scheme+"://"+addr+":"+portStr+"/")
		}
	}

	sort.Strings(baseUrls)

	return baseUrls
}
