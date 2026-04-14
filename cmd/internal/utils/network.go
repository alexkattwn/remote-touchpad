package utils

import (
	"net"
	"strings"
)

func GetLocalIP() string {
	addrs, _ := net.InterfaceAddrs()

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {

			ip := ipnet.IP.String()

			if strings.HasSuffix(ip, ".1") {
				continue
			}

			if strings.HasPrefix(ip, "192.168.") {
				return ip
			}
		}
	}
	return "localhost"
}