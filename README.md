# Go-GeoIP
Go-GeoIP: 简单且高性能的 IP 地址地理信息查询服务

## 安装运行
### 编译打包
```
chmod +x control
./control build
./control pack
```

### 启动服务
根据 `cfg.example.json` 创建 `cfg.json` 配置文件，并根据自己的实际情况进行修改 
```
mv cfg.example.json cfg.json
./control start
```
### systemctl 托管
假定部署在 `/opt/go-geoip` 目录下，如果部署在其他目录修改 `go-geoip.service` 中的 `WorkingDirectory` 和 `ExecStart` 两个字段即可。
```
cp go-geoip.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable go-geoip
systemctl start go-geoip
```

### 定制页面
修改 `templates` 目录下的 `index.html` 即可，相关资源文件在 `assets` 目录下。

## 配置说明

```json
{
	"logger": {
		"dir": "logs/",
		"level": "DEBUG",
		"keepHours": 24
	},
	"campus": {
		"continent": "亚洲",
		"country": "中国",
		"province": "上海",
		"city": "上海",
		"district": "华东师范大学",
		"isp": "校园网",
		"areaCode": "310000",
		"countryEnglish": "China",
		"countryCode": "CN",
		"longitude": "",
		"latitude": "",
		"ips": [
			"10.0.0.0/8",
			"192.168.0.0/16",
			"172.16.0.0/12"
		]
	},
	"db": {
		"maxmind": "GeoLite2-City.mmdb",
		"qqzengip": "qqzeng-ip-3.0-ultimate.dat"
	},
	"source": {
		"ipv4": "maxmind",
		"ipv6": "maxmind"
	},
	"http": {
		"listen": "0.0.0.0:80",
		"trustProxy": ["127.0.0.1", "::1"],
		"x-api-key": "this-is-key"
	}
}
```

|配置项|类型|说明|
|---|---|---|
|logger|object|一个包含日志设置的部分|
|logger.dir|string|存储日志文件的目录|
|logger.level|string|日志的级别，比如DEBUG, INFO, WARN, 或ERROR|
|logger.keepHours|number|保留日志文件的小时数，之后删除|
|campus|object|一个包含园区内网信息的部分|
|campus.continent|string|园区所在的洲|
|campus.country|string|园区所在的国家|
|campus.province|string|园区所在的省份|
|campus.city|string|园区所在的城市|
|campus.district|string|园区所在的区县（行政区）|
|campus.isp|string|园区的ISP运营商|
|campus.areaCode|string|园区所在的行政区划代（国内部分）|
|campus.countryEnglish|string|园区所在国家的英文名|
|campus.countryCode|string|园区所在国家的国家代码|
|campus.longitude|string|园区的经度|
|campus.latitude|string|园区的纬度|
|campus.ips|array|属于该园区的IP范围的数组，命中这部分的IP地址，将使用配置文件中的内容进行返回|
|db|object|一个包含数据库设置的部分|
|db.maxmind|string|MaxMind GeoLite2数据库的文件的路径|
|db.qqzengip|string|qqzengip数据库的文件的路径|
|source|object|一个包含IP信息来源设置的部分|
|source.ipv4|string|IPv4信息的来源，目前支持maxmind或qqzengip|
|source.ipv6|string|IPv6信息的来源，目前支持maxmind|
|http|object|一个包含HTTP服务器设置的部分|
|http.listen|string|HTTP服务器监听的地址和端口|
|http.trustProxy|array|被信任的代理的IP地址的数组，当服务被发布在反向代理后时必须正确配置，否则无法正确获取到 xff 的地址。|
|http.x-api-key|string|访问openapi接口所需的API密钥|

## 数据库自动更新
### maxmind
对于 `maxmind` 的数据源，可以使用其官方提供的 `GeoIP Update` 程序，它是一个命令行工具，可以定期下载和安装最新的数据库。你可以在这里找到更多关于`GeoIP Update` 的信息 [Updating GeoIP and GeoLite Databases](https://dev.maxmind.com/geoip/updating-databases?lang=en)

## API
### myip

myip 的接口用于返回请求者的 IP 地址，对于一些无浏览器的终端，可以使用这个接口方便的获取自身的IP地址信息（特别是 nat 后的）。

它也可以被配置了 CORS 的网站通过前端调用（ToDO）

提供了简单字符串与 json 格式化两种风格接口。

```
# curl http://go-geoip/myip
# 192.168.0.100
```

```
# curl http://go-geoip/myip/format
# {"errCode":0,"errMsg":"success","requestId":"0f40823e-04ce-4def-9af2-71e7e1403ec8","data":{"ip":"192.168.0.100"}}
```

### searchapi
searchapi 接口面向浏览器，提供了一个 IP 地址的查询接口，并输出转换好的字符串以简化前端解析。

这个接口受验证码和限流措施的保护，以防范可能的恶意爬虫（ToDo）

他访问的路径是 `http://go-geoip/ip`

### openapi
openapi 接口面向第三方应用，提供了一个 IP 地址的查询接口，通过 X-API-KEY 进行授权校验。

建议在多租户的情况下，进一步通过 API 网关进行代理封装和授权分发。

- request
```
curl -H "X-API-KEY: this-is-key" http://go-geoip/api/v1/network/ip?ip=2001:da8:8005:a405:250:56ff:feaf:8c28
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


## brenchmark
配置|内容
--|--
环境|阿里云
CPU|2核(vCPU)
内存|8 GiB
操作系统|Anolis OS 8.6 RHCK 64位

```
# go test -bench=.  -benchmem

goos: linux
goarch: amd64
pkg: github.com/ECNU/go-geoip
cpu: Intel(R) Xeon(R) Platinum 8369B CPU @ 2.70GHz
BenchmarkIndex-2             	  244190	      4271 ns/op	   10000 B/op	      15 allocs/op
BenchmarkSeachAPIForIPv4-2   	  782768	      1741 ns/op	    1904 B/op	      15 allocs/op
BenchmarkSeachAPIForIPv6-2   	  818250	      1744 ns/op	    1904 B/op	      15 allocs/op
BenchmarkOpenAPIForIPv4-2    	  394813	      3383 ns/op	    2592 B/op	      23 allocs/op
BenchmarkOpenAPIForIPv6-2    	  391868	      3378 ns/op	    2592 B/op	      23 allocs/op
PASS
ok  	github.com/ECNU/go-geoip	7.044s
```