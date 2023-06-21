# Open-GeoIP
Open-GeoIP: 简单且高性能的 IP 地址地理信息查询服务

![](https://github.com/ECNU/open-geoip/blob/main/demo.jpg?raw=true)

- [Open-GeoIP](#open-geoip)
	- [安装运行](#安装运行)
		- [二进制直接运行](#二进制直接运行)
		- [systemctl 托管](#systemctl-托管)
		- [数据库自动更新](#数据库自动更新)
			- [maxmind](#maxmind)
		- [编译打包](#编译打包)
		- [定制页面](#定制页面)
	- [配置说明](#配置说明)
	- [内部 IP 地理数据库](#内部-ip-地理数据库)
	- [限流方案](#限流方案)
	- [高可用与扩展性](#高可用与扩展性)
	- [API](#api)
		- [myip](#myip)
		- [mylocation](#mylocation)
		- [searchapi](#searchapi)
		- [openapi](#openapi)
	- [benchmark](#benchmark)
	- [鸣谢](#鸣谢)

## 安装运行

### 二进制直接运行
在 [release](https://github.com/ECNU/open-geoip/releases) 中下载最新的 [release] 包，解压后直接运行即可。

注意：`release` 中内置的数据库文件来自于 [ipdb-go](https://github.com/ipipdotnet/ipdb-go) 中的 `city.free.ipdb`，仅供测试使用，不保证数据的准确性。

如应用于生产环境，请获取商用授权，或者[注册](https://www.maxmind.com/en/geolite2/signup) `maxmind` 的账号后，获取免费版的 `GeoLite2-City.mmdb` 数据库文件，并更新配置文件替换数据源为 `maxmind`。

```
tar -zxvf open-geoip-0.1.0-linux-amd64.tar.gz
cd open-geoip/
./control start
```
访问你服务器的 80 端口即可使用。


### systemctl 托管
假定部署在 `/opt/open-geoip` 目录下，如果部署在其他目录修改 `open-geoip.service` 中的 `WorkingDirectory` 和 `ExecStart` 两个字段即可。
```
cp open-geoip.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable open-geoip
systemctl start open-geoip
```

### 数据库自动更新
#### maxmind
如果需要自动更新 `mmdb` 数据库，只需要在[注册](https://www.maxmind.com/en/geolite2/signup)一个 `maxmind` 的账号，获得一个 [LicenseKey](https://www.maxmind.com/en/accounts/current/license-key) ，并将他配置到 `cfg.json` 中的 `AutoDownload.MaxmindLicenseKey` 中，或者配置到系统环境变量 `MAXMIND_LICENSE_KEY` 中即可。


### 编译打包
```
git clone https://github.com/ECNU/open-geoip.git
cd open-geoip/
chmod +x control
./control pack
```

### 定制页面
修改 `templates` 目录下的 `index.html` 即可，相关资源文件在 `assets` 目录下。

## 配置说明

根据 `cfg.json.example` 文件，创建 `cfg.json` 文件，再进一步根据自己的需要修改配置。

```json
{
	"logger": {
		"dir": "logs/",
		"level": "DEBUG",
		"keepHours": 24
	},
	"redis": {
		"dsn": "127.0.0.1:6379",
		"maxIdle": 5,
		"connTimeout": 5,
		"readTimeout": 5,
		"writeTimeout": 5,
		"password": ""
	},
	"internal": {
		"source": "maxmind", 
		"enable": false, 
		"db": "GeoLite2-City.mmdb.test"
	},
	"db": {
		"maxmind": "GeoLite2-City.mmdb",
		"qqzengip": "",
		"ipdb":""
	},
	"source": {
		"ipv4": "maxmind",
		"ipv6": "maxmind"
	},
	"autoDownload":{
		"enabled":false,
		"MaxmindLicenseKey":"",
		"targetFilePath":"",
		"timeout":3,
		"interval":24
	},
	"rateLimit": {
		"enabled": false,
		"minute": 100,
		"hour": 1000,
		"day": 10000
	},
	"http": {
		"listen": "0.0.0.0:80",
		"trustProxy": ["127.0.0.1", "::1"],
		"cors":["http://localhost"],
		"x-api-key": "this-is-key"
	}
}
```

| 配置项                            | 类型     | 说明                                                                                                                  |
|--------------------------------|--------|---------------------------------------------------------------------------------------------------------------------|
| logger                         | object | 一个包含日志设置的部分                                                                                                         |
| logger.dir                     | string | 存储日志文件的目录                                                                                                           |
| logger.level                   | string | 日志的级别，比如DEBUG, INFO, WARN, 或ERROR                                                                                   |
| logger.keepHours               | number | 保留日志文件的小时数，之后删除                                                                                                     |
| redis | object | redis 配置参数，配合限流策略使用 |
| redis.dsn | string | redis 的连接地址 |
| redis.maxIdle | number | redis 的最大空闲连接数 |
| redis.connTimeout | number | redis 的连接超时时间，单位是 second |
| redis.readTimeout | number | redis 的读取超时时间，单位是 second |
| redis.writeTimeout | number | redis 的写入超时时间，单位是 second |
| redis.password | string | redis 的密码 |
| internal                       | object  | 内部数据库参数                                                                                                             | 
| internal.enabled               | bool   | 开启内部数据库                                                                                                             |
| internal.source                | string | 内部数据库来源                                                                                                             |
| internal.db                    | string | 内部数据库文件路径      
| db                             | object | 一个包含数据库设置的部分                                                                                                        |
| db.maxmind                     | string | MaxMind GeoLite2数据库的文件的路径，如果 autDownload 配置为 true，那么这里的配置不会生效                                                       |
| db.qqzengip                    | string | qqzengip数据库的文件的路径                                                                                                   |
| db.ipdb                        | string | ipip.net数据库的文件的路径                                                                                                   |
| source                         | object | 一个包含IP信息来源设置的部分                                                                                                     |
| source.ipv4                    | string | IPv4信息的来源，可配置为 [maxmind](https://www.maxmind.com)/[qqzengip](https://www.qqzeng.com/)/[ipdb](https://www.ipip.net/) |
| source.ipv6                    | string | IPv6信息的来源，可配置为 [maxmind](https://www.maxmind.com)/[qqzengip](https://www.qqzeng.com/)/[ipdb](https://www.ipip.net/) |
| autoDownload                   | object | 一个包含自动更新数据库的设置的部分                                                                                                   |
| autoDownload.enabled           | bool   | 是否启用自动更新数据库                                                                                                         |
| autoDownload.MaxmindLicenseKey | string | MaxMind License Key，用于自动更新 `MaxMind GeoLite2` 数据库，也可以配置在环境变量 `MAXMIND_LICENSE_KEY` 中。如果都没有配置，那么 `maxmind` 的自动更新会报错  |
| autoDownload.targetFilePath    | string | 自动更新数据库的目标文件路径，如果不配置此参数，默认值是 `./`，自动更新数据库会下载到这个目录                                                                   |
| autoDownload.timeout           | number | 自动更新数据库的超时时间，单位是 second，如果不配置此参数，默认值是 3                                                                             |
| autoDownload.interval          | number | 自动更新数据库的间隔时间，单位是 hour，如果不配置此参数，默认值是24                                                                               |
| rateLimit                      | object | 一个包含限流设置的部分                                                                                                         |
| rateLimit.enabled              | bool   | 是否启用限流策略                                                                                                             |
| rateLimit.minute               | number | 每分钟最多访问次, 0 表示不限制数                                                                                                           |
| rateLimit.hour                 | number | 每小时最多访问次, 0 表示不限制数                                                                                                           |
| rateLimit.day                  | number | 每天最多访问次, 0 表示不限制数                                                                                                             |
| http                           | object | 一个包含HTTP服务器设置的部分                                                                                                    |
| http.listen                    | string | HTTP服务器监听的地址和端口                                                                                                     |
| http.trustProxy                | array  | 被信任的代理的IP地址的数组，当服务被发布在反向代理后时必须正确配置，否则无法正确获取到 xff 的地址。                                                               |
| http.cors                      | array  | 允许跨域访问的域名列表,配置内的域名可以跨域访问 `/myip` 和 `/myip/format` 接口                                                                |
| http.x-api-key                 | string | 访问 openapi 接口所需的 API 密钥                                                                                             |                                                                                                     |
## 内部 IP 地理数据库

Open-GeoIP 允许以导入的方式，构建企业内部自己的 IP 地理数据库，以便于查询内部 IP 地址的物理位置。

导入的格式是 `csv`，内容如下所示，项目中已经存在一个 `internal.csv` 的示例文件，可以参考。

| continent | country | province | city | district | isp | areaCode | countryCode | countryEnglish | longitude | latitude | ip_subnet |
| --------- | ------- | -------- | ---- | -------- | --- | -------- | ----------- | -------------- | --------- | -------- | --------- |
|           |         |          |      | 保留     | 回环地址  |          |             |                |           |          | 127.0.0.0/8 |
| 亚洲      | 中国    | 上海     | 上海  | 开源教育   | 企业内网  | 310000   | CN          | China          ||| 10.0.0.0/8 |
| 亚洲      | 中国    | 上海     | 上海  | 开源教育   | 企业内网  ||| CN          || China          ||| 192.168.0.0/16 |
| 亚洲      ||| 中国    || 上海     || 上海  || 开源教育   || 企业内网  ||| 310000   || CN          || China          ||| 172.16.0.0/12 |
| 亚洲      ||| 中国    || 上海     || 上海  || 开源教育   || 企业内网  ||| 310000   || CN          || China          |||| fd00::/8 |


在启动 Open-GeoIP 之前，执行 `-csv` 命令即可导入内部数据库，此时默认会生成一个 `internal.mmdb` 文件。
```
./open-geoip -csv internal.csv
```
在配置文件中，修改 `internal.mmdb` 的相关配置，将其开启即可。

```json
        "internaldb": {
                "source": "maxmind",
                "enabled": true,
                "db": "internal.mmdb"
        }
```


## 限流方案
Open-GeoIP 通过 redis 记录每个IP地址的访问次数，当超过阈值时，对该IP进行限制访问。支持分钟，小时，天 三种颗粒的计数策略，可以通过配置文件中的 ratelimit 的部分进行配置，以下示例表示开启了限流策略，并限制了每分钟最多访问 100 次，每小时最多访问 1000 次，每天最多访问 10000 次。

```json
	"rateLimit": {
		"enabled": true,
		"minute": 100,
		"hour": 1000,
		"day": 10000
	},
```

## 高可用与扩展性
Open-GeoIP 是无状态的，因此可以任意的进行横向扩展并通过负载均衡实现高可用。在启用限流方案时，多个 Open-GeoIP 可以通过共享同一个 Redis 服务实现限流计数的一致性。

## API
### myip

myip 的接口用于返回请求者的 IP 地址，对于一些无浏览器的终端，可以使用这个接口方便的获取自身的IP地址信息（特别是 nat 后的）。

它也可以被配置了 CORS 的网站通过前端调用

提供了简单字符串与 json 格式化两种风格接口。

```
# curl http://localhost/myip
# 192.168.0.100
```

```
# curl http://localhost/myip/format
# {"errCode":0,"errMsg":"success","requestId":"0f40823e-04ce-4def-9af2-71e7e1403ec8","data":{"ip":"192.168.0.100"}}
```

### mylocation

mylocation 的接口用于返回请求者的 IP 地址对应的物理位置。

它也可以被配置了 CORS 的网站通过前端调用

提供了简单字符串与 json 格式化两种风格接口。

```
# curl http://localhost/mylocation
# 保留地址
```

```
# curl http://localhost/mylocation/format
# {"errCode":0,"errMsg":"success","requestId":"c2e8c50e-b55f-455a-a9d4-d209acd20ab9","data":{"ip":"::1","continent":"保留地址","country":"","province":"","city":"","district":"","isp":"","areaCode":"","countryEnglish":"","countryCode":"","longitude":"","latitude":""}}
```

### searchapi
searchapi 接口面向浏览器，提供了一个 IP 地址的查询接口，并输出转换好的字符串以简化前端解析。

这个接口受验证码（todo）和限流措施的保护，以防范可能的恶意爬虫

他访问的路径是 `http://localhost/ip`

### openapi
openapi 接口面向第三方应用，提供了一个 IP 地址的查询接口，通过 X-API-KEY 进行授权校验。

建议在多租户的情况下，进一步通过 API 网关进行代理封装和授权分发。

- request
```
curl -H "X-API-KEY: this-is-key" http://localhost/api/v1/network/ip?ip=2001:da8:8005:a405:250:56ff:feaf:8c28
```

- response
```json
{
	"errCode": 0,
	"errMsg": "success",
	"requestId": "7ead62f7-3f15-4822-ad1e-cf7915a8299f",
	"data": {
		"ip": "2001:da8:8005:a405:250:56ff:feaf:8c28",
		"continent": "亚洲",
		"country": "中国",
		"province": "上海",
		"city": "上海",
		"district": "",
		"isp": "",
		"areaCode": "",
		"countryEnglish": "China",
		"countryCode": "CN",
		"longitude": "121.458100",
		"latitude": "31.222200"
	}
}
```


## benchmark
基于 `maxmind` 数据库，`web` 服务性能测试
```
# go test -bench=.  -benchmem

goos: linux
goarch: amd64
pkg: github.com/ECNU/open-geoip
cpu: Intel(R) Xeon(R) Platinum 8369B CPU @ 2.70GHz
BenchmarkIndex-2             	  244190	      4271 ns/op	   10000 B/op	      15 allocs/op
BenchmarkSeachAPIForIPv4-2   	  782768	      1741 ns/op	    1904 B/op	      15 allocs/op
BenchmarkSeachAPIForIPv6-2   	  818250	      1744 ns/op	    1904 B/op	      15 allocs/op
BenchmarkOpenAPIForIPv4-2    	  394813	      3383 ns/op	    2592 B/op	      23 allocs/op
BenchmarkOpenAPIForIPv6-2    	  391868	      3378 ns/op	    2592 B/op	      23 allocs/op
PASS
ok  	github.com/ECNU/open-geoip	7.044s
```

## 鸣谢

本项目的一些主要功能使用了以下开源项目，更多的依赖详见 `go.mod` 。

感谢他们的开源精神。

- `web` 服务 —— [gin](https://github.com/gin-gonic/gin) 
- `maxmind` 解析 —— [geoip2-golang](https://github.com/oschwald/geoip2-golang)
- `maxming` 自动更新 —— [go-geoip](https://github.com/pieterclaerhout/go-geoip)
- `ipdb` 解析 —— [ipdb-go](https://github.com/ipipdotnet/ipdb-go)
- `qqzengip` 解析 —— [qqzeng-ip](https://https://github.com/zengzhan/qqzeng-ip)