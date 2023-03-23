package models

import (
	"fmt"
	"net"
	"strings"

	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/util"
	"github.com/ipipdotnet/ipdb-go"
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
	IpdbReader    *ipdb.City
}

var ipReader *IpReader

func InitReader() error {
	ipReader = &IpReader{}
	if g.Config().DB.Qqzengip != "" {
		instance, err := LoadDat(g.Config().DB.Qqzengip)
		if err != nil {
			return err
		}
		ipReader.QqzengReader = instance
	}
	if g.Config().DB.Ipdb != "" {
		city, err := ipdb.NewCity(g.Config().DB.Ipdb)
		if err != nil {
			return err
		}
		ipReader.IpdbReader = city

	}

	if g.Config().Source.IPv4 == "maxmind" || g.Config().Source.IPv6 == "maxmind" {
		if g.Config().AutoDownload.Enabled {
			dbPath, err := util.AutoDownloadMaxmindDatabase(g.Config().AutoDownload)
			if err != nil {
				return err
			}
			db, err := geoip2.Open(dbPath)
			if err != nil {
				return err
			}
			ipReader.MaxMindReader = db
		} else {
			db, err := geoip2.Open(g.Config().DB.Maxmind)
			if err != nil {
				return err
			}
			ipReader.MaxMindReader = db
		}
	}

	return nil
}

func copyIPdb(src *ipdb.CityInfo, dst *IpGeo, language string) {
	dst.AreaCode = src.ChinaAdminCode
	dst.City = src.CityName
	dst.Continent = src.ContinentCode
	dst.Country = src.CountryName
	dst.CountryCode = src.CountryCode
	dst.District = src.DistrictName
	dst.ISP = src.IspDomain
	dst.Latitude = src.Latitude
	dst.Longitude = src.Longitude
	dst.Province = src.RegionName
	return
}

//Todo 国际化拓展
func switchIpdbLanguage(lan string) string {
	switch lan {
	case "zh-CN":
		return "CN"
	default:
		return "CN"
	}
}

func readSource(ipNet net.IP, source, language string) (ipGeo IpGeo, err error) {
	ipGeo.IP = ipNet.String()
	switch source {
	case "maxmind":
		var record *geoip2.City
		record, err = ipReader.MaxMindReader.City(ipNet)
		if err != nil {
			return
		}
		copyGeoIP(record, &ipGeo, language)
		return
	case "qqzengip":
		ipGeo = ipReader.QqzengReader.Get(ipNet.String())
		return
	case "ipdb":
		lan := switchIpdbLanguage(language)
		var cityInfo *ipdb.CityInfo
		cityInfo, err = ipReader.IpdbReader.FindInfo(ipNet.String(), lan)
		if err != nil {
			return
		}
		copyIPdb(cityInfo, &ipGeo, language)
		return
	default:
		return
	}
	return
}

func copyGeoIP(src *geoip2.City, dst *IpGeo, language string) {
	if _, ok := src.Continent.Names[language]; ok {
		dst.Continent = src.Continent.Names[language]
	}
	if _, ok := src.Country.Names[language]; ok {
		dst.Country = src.Country.Names[language]
	}
	if _, ok := src.Country.Names["en"]; ok {
		dst.CountryEnglish = src.Country.Names["en"]
	}
	if _, ok := src.City.Names[language]; ok {
		dst.City = src.City.Names[language]
	}
	dst.CountryCode = src.Country.IsoCode
	if len(src.Subdivisions) > 0 {
		dst.Province = src.Subdivisions[0].Names[language]
	}
	dst.Latitude = fmt.Sprintf("%f", src.Location.Latitude)
	dst.Longitude = fmt.Sprintf("%f", src.Location.Longitude)
	return
}

func GetIP(ipStr string, config g.SourceConfig) (ipGeo IpGeo, err error) {
	//ToDO 支持国际化
	language := "zh-CN"
	ipGeo.IP = ipStr
	ipNet := net.ParseIP(ipStr)
	if ipNet == nil {
		err = fmt.Errorf("invalid IP address format: %s", ipStr)
		return
	}
	//ipv6
	if strings.Contains(ipStr, ":") {
		ipGeo, err = readSource(ipNet, config.IPv6, language)
		if err != nil {
			return
		}
		return
	}
	//ipv4
	ipGeo, err = readSource(ipNet, config.IPv4, language)
	if err != nil {
		return
	}
	return
}
