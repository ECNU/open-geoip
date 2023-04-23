package g

import (
	"encoding/json"

	"log"
	"sync"

	"github.com/toolkits/file"
)

/*
GlobalConfig 全局配置
*/
type GlobalConfig struct {
	Logger       LoggerSection      `json:"logger"`
	DB           DBConfig           `json:"db"`
	Campus       CampusConfig       `json:"campus"`
	InternalDB   InternalDBConfig   `json:"internaldb"`
	Source       SourceConfig       `json:"source"`
	AutoDownload AutoDownloadConfig `json:"autoDownload"`
	Http         HttpConfig         `json:"http"`
}

/*
AutoDownloadConfig 自动下载的配置
*/
type AutoDownloadConfig struct {
	Enabled           bool   `json:"enabled"`
	MaxmindLicenseKey string `json:"maxmindLicenseKey"`
	TargetFilePath    string `json:"targetFilePath"`
	Timeout           int    `json:"timeout"`
	Interval          int    `json:"interval"`
}

/*
DBConfig DB 配置
*/
type DBConfig struct {
	Maxmind  string `json:"maxmind"`
	Qqzengip string `json:"qqzengip"`
	Ipdb     string `json:"ipdb"`
}

type InternalDBConfig struct {
	Source string `json:"source"`
	Enable bool   `json:"enable"`
	DB     string `json:"db"`
}

/*
CampusConfig 园区内网配置
*/
type CampusConfig struct {
	Continent      string   `json:"continent"`      //州
	Country        string   `json:"country"`        //国家
	Province       string   `json:"province"`       //省
	City           string   `json:"city"`           //城市
	District       string   `json:"district"`       //区县(行政区）
	ISP            string   `json:"isp"`            //运营商
	AreaCode       string   `json:"areaCode"`       //行政区划代码（国内）
	CountryEnglish string   `json:"countryEnglish"` //国家英文名
	CountryCode    string   `json:"countryCode"`    //国家英文代码
	Longitude      string   `json:"longitude"`      //经度
	Latitude       string   `json:"latitude"`       //纬度
	IPs            []string `json:"ips"`
}

/*
SourceConfig 数据源
*/
type SourceConfig struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

/*
HttpConfig Http 配置
*/
type HttpConfig struct {
	Listen     string   `json:"listen"`
	TrustProxy []string `json:"trustProxy"`
	XAPIKey    string   `json:"x-api-key"`
	CORS       []string `json:"cors"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

/*
Config 安全的读取和修改配置
*/
func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

/*
ParseConfig 加载配置
*/
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c
}
