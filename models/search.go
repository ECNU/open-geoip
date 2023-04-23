package models

import (
	"fmt"
	"net"

	"github.com/toolkits/pkg/logger"

	"github.com/ECNU/open-geoip/g"
)

func (self IpGeo) ToString() string {
	//只返回必要的部分，给前端先使用
	res := fmt.Sprintf("%s %s %s %s %s %s", self.Continent, self.Country, self.Province, self.City, self.District, self.ISP)
	return res
}

func SearchIP(ipStr string) (result IpGeo) {
	//无论发生什么，IP 永远返回
	result.IP = ipStr

	////先检查campus ip，匹配则直接返回 District=华东师范大学， ISP=校园网，经纬度给空字符串
	//if IPCheck(ipStr, g.Config().Campus.IPs) {
	//	copier.Copy(&result, g.Config().Campus)
	//	return
	//}

	var err error
	result, err = GetIP(ipStr, g.Config().Source)

	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func CheckIPValid(ipStr string) bool {
	ipNet := net.ParseIP(ipStr)
	if ipNet == nil {
		return false
	}
	return true
}
