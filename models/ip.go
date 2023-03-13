package models

import (
	"fmt"
	"net"
	"strings"

	"github.com/ECNU/go-geoip/g"
	"github.com/oschwald/geoip2-golang"
)

type IpGeo struct {
	IP             string `json:"ip"`             //IP地址
	Continent      string `json:"continent"`      //州
	Country        string `json:"country"`        //国家
	Province       string `json:"province"`       //省
	City           string `json:"city"`           //城市
	District       string `json:"district"`       //区县(行政区）
	ISP            string `json:"isp"`            //运营商
	AreaCode       string `json:"areaCode"`       //行政区划代码（国内）
	CountryEnglish string `json:"countryEnglish"` //国家英文名
	CountryCode    string `json:"countryCode"`    //国家英文代码
	Longitude      string `json:"longitude"`      //经度
	Latitude       string `json:"latitude"`       //纬度
}

type IpReader struct {
	MaxMindReader *geoip2.Reader
	QqzengReader  *IpSearch
}

var ipReader *IpReader

func InitReader(config g.DBConfig) (err error) {
	ipReader = &IpReader{}
	if config.Qqzengip != "" {
		instance, err := LoadDat(config.Qqzengip)
		if err != nil {
			return err
		}
		ipReader.QqzengReader = instance
	}
	if config.Maxmind != "" {
		db, err := geoip2.Open(config.Maxmind)
		if err != nil {
			return err
		}
		ipReader.MaxMindReader = db
	}

	return
}

func copyGeoIP(src *geoip2.City, dst *IpGeo) {
	if _, ok := src.Continent.Names["zh-CN"]; ok {
		dst.Continent = src.Continent.Names["zh-CN"]
	}
	if _, ok := src.Country.Names["zh-CN"]; ok {
		dst.Country = src.Country.Names["zh-CN"]
	}
	if _, ok := src.Country.Names["en"]; ok {
		dst.CountryEnglish = src.Country.Names["en"]
	}
	if _, ok := src.City.Names["zh-CN"]; ok {
		dst.City = src.City.Names["zh-CN"]
	}
	dst.CountryCode = src.Country.IsoCode
	if len(src.Subdivisions) > 0 {
		dst.Province = src.Subdivisions[0].Names["zh-CN"]
	}
	dst.Latitude = fmt.Sprintf("%f", src.Location.Latitude)
	dst.Longitude = fmt.Sprintf("%f", src.Location.Longitude)
	return
}

func readSource(ipNet net.IP, source string) (ipGeo IpGeo, err error) {
	ipGeo.IP = ipNet.String()
	switch source {
	case "maxmind":
		var record *geoip2.City
		record, err = ipReader.MaxMindReader.City(ipNet)
		if err != nil {
			return
		}
		copyGeoIP(record, &ipGeo)
		return
	case "qqzengip":
		ipGeo = ipReader.QqzengReader.Get(ipNet.String())
		return
	default:
		return
	}
	return
}

func GetIP(ipStr string, config g.SourceConfig) (ipGeo IpGeo, err error) {
	ipGeo.IP = ipStr
	ipNet := net.ParseIP(ipStr)
	if ipNet == nil {
		err = fmt.Errorf("invalid IP address format: %s", ipStr)
		return
	}
	//ipv6
	if strings.Contains(ipStr, ":") {
		ipGeo, err = readSource(ipNet, config.IPv6)
		if err != nil {
			return
		}
		return
	}
	//ipv4
	ipGeo, err = readSource(ipNet, config.IPv4)
	if err != nil {
		return
	}
	return
}
