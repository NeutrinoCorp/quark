package quark

import (
	"net"
)

// 	deprecated
// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	// handle err...
	if err != nil {
		return ""
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
