package models

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
	"strings"

	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/util"
	"github.com/ipipdotnet/ipdb-go"
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
	MaxMindReader         *geoip2.Reader
	InternalMaxMindReader *InternalReader
	QqzengReader          *IpSearch
	IpdbReader            *ipdb.City
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

	if g.Config().InternalDB.Enabled {
		db, err := Open(g.Config().InternalDB.DB)
		if err != nil {
			return err
		}
		ipReader.InternalMaxMindReader = db
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

// Todo 国际化拓展
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

func copyInternalGeoIP(src *InternalGeoIP, dst *IpGeo, language string) {
	if _, ok := src.City.Continent.Names[language]; ok {
		dst.Continent = src.City.Continent.Names[language]
	}
	if _, ok := src.City.Country.Names[language]; ok {
		dst.Country = src.City.Country.Names[language]
	}
	if _, ok := src.City.Country.Names["en"]; ok {
		dst.CountryEnglish = src.City.Country.Names["en"]
	}
	if _, ok := src.City.City.Names[language]; ok {
		dst.City = src.City.City.Names[language]
	}
	dst.CountryCode = src.City.Country.IsoCode
	dst.AreaCode = src.Internal.AreaCode
	if len(src.City.Subdivisions) > 0 {
		dst.Province = src.City.Subdivisions[0].Names[language]
	}
	if _, ok := src.Internal.ISP[language]; ok {
		dst.ISP = src.Internal.ISP[language]
	}
	if _, ok := src.Internal.District[language]; ok {
		dst.District = src.Internal.District[language]
	}
	dst.Latitude = fmt.Sprintf("%f", src.City.Location.Latitude)
	dst.Longitude = fmt.Sprintf("%f", src.City.Location.Longitude)

	return
}

func readInternalSource(ipNet net.IP, source, language string) (ipGeo IpGeo, err error) {
	ipGeo.IP = ipNet.String()
	switch source {
	case "maxmind":

		var record *InternalGeoIP

		record, err = ipReader.InternalMaxMindReader.City(ipNet)

		if err != nil {
			return
		}

		copyInternalGeoIP(record, &ipGeo, language)
		return
	}

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
	// 本地优先
	if g.Config().InternalDB.Enabled {
		ipGeo, err = readInternalSource(ipNet, g.Config().InternalDB.Source, language)

		if err != nil {
			return
		}

		// ToString 9个字段 为空是8
		if len(ipGeo.ToString()) != 8 {
			return
		}

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
