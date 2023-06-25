package models

import (
	"fmt"
	"net"

	"github.com/ECNU/open-geoip/g"
	"github.com/toolkits/pkg/logger"
)

func (self IpGeo) ToString() string {
	//只返回必要的部分，给前端先使用
	res := fmt.Sprintf("%s %s %s %s %s %s %s %s %s", self.Continent, self.Country, self.CountryEnglish, self.CountryCode, self.AreaCode, self.Province, self.City, self.District, self.ISP)
	return res
}

func SearchIP(ipStr string, isApi bool, isAuth bool) (result IpGeo) {
	//无论发生什么，IP 永远返回
	result.IP = ipStr

	var err error
	result, err = GetIP(ipStr, g.Config().Source, isApi, isAuth)

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
