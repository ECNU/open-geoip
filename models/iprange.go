package models

import (
	"encoding/binary"
	"net"

	"strings"
)

// IPCheck 检查 ip 是否在特定的 ip 地址范围内
func IPCheck(thisip string, ips []string) bool {
	for _, ip := range ips {
		ip = strings.TrimRight(ip, "/")
		if strings.Contains(ip, "/") {
			if ipCheckMask(thisip, ip) {
				return true
			}
		} else if strings.Contains(ip, "-") {
			ipRange := strings.SplitN(ip, "-", 2)
			if ipCheckRange(thisip, ipRange[0], ipRange[1]) {
				return true
			}
		} else {
			if thisip == ip {
				return true
			}
		}
	}
	return false
}

func ipCheckRange(ip, ipStart, ipEnd string) bool {
	thisIP := net.ParseIP(ip)
	firstIP := net.ParseIP(ipStart)
	endIP := net.ParseIP(ipEnd)
	if thisIP.To4() == nil || firstIP.To4() == nil || endIP.To4() == nil {
		return false
	}
	firstIPNum := ipToInt(firstIP.To4())
	endIPNum := ipToInt(endIP.To4())
	thisIPNum := ipToInt(thisIP.To4())
	if thisIPNum >= firstIPNum && thisIPNum <= endIPNum {
		return true
	}
	return false
}

func ipCheckMask(ip, ipMask string) bool {
	_, subnet, _ := net.ParseCIDR(ipMask)

	thisIP := net.ParseIP(ip)
	return subnet.Contains(thisIP)
}

func ipToInt(ip net.IP) int32 {
	return int32(binary.BigEndian.Uint32(ip.To4()))
}
