package utils

import (
	"net"
	"time"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/19 20:24
  @describe :
*/

// GetMasterIP 获取主IP地址
func GetMasterIP() string {
	conn, err := net.DialTimeout("udp", "8.8.8.8:8", time.Second*5)
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GetLocalIPv4 获取当前IPv4地址
func GetLocalIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && !ipNet.IP.IsLinkLocalUnicast() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
