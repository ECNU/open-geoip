package models

import (
	"testing"

	"log"

	"github.com/ECNU/go-geoip/g"
	"github.com/stretchr/testify/assert"
)

func init() {
	g.ParseConfig("cfg.json")
	err := InitReader(g.Config().DB)
	if err != nil {
		log.Fatalf("load geo db failed, %v", err)
	}
}

func Test_IpCheck(t *testing.T) {

	var campusIPs = []string{
		"10.10.10.0/8",
		"192.168.0.0/24",
		"172.16.0.0/12",
		"192.168.10.1-192.168.10.250",
		"192.168.100.3",
		"2001:da8:8005::/48",
	}

	var ipCheckList = []struct {
		ip  string
		exp bool
	}{
		{"10.10.3.1", true},
		{"192.168.12.3", false},
		{"192.168.100.3", true},
		{"192.168.10.100", true},
		{"172.20.3.10", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"2001:da8:8005:abcd:1234::8888", true},
		{"2001:da8:8000:abcd:1234::8888", false},
	}

	for _, r := range ipCheckList {
		out := IPCheck(r.ip, campusIPs)
		assert.Equal(t, out, r.exp)
	}
}

func Test_GetIP(t *testing.T) {
	config := g.Config().Source

	// ipv4 和 ipv6 均采用 maxmind 数据源
	config.IPv4 = "maxmind"
	config.IPv6 = "maxmind"
	res, _ := GetIP("202.120.92.60", config)
	assert.Equal(t, res.Country, "中国")
	res, _ = GetIP("2001:da8:8005:a492::60", config)
	assert.Equal(t, res.Country, "中国")

	// ipv4 切换为 qqzengip
	config.IPv4 = "qqzengip"
	config.IPv6 = "qqzengip"
	res, _ = GetIP("202.120.92.60", config)
	assert.Equal(t, res.Country, "中国")
	// qqzengip 不支持 ipv6，此时应该查不到
	res, _ = GetIP("2001:da8:8005:a492::60", config)
	assert.Equal(t, res.Country, "")

	//非法的请求
	_, err := GetIP("201..1", config)
	assert.NotNil(t, err)
	_, err = GetIP("2001:da8:::::::::", config)
	assert.NotNil(t, err)

}

func TestSearchIP(t *testing.T) {
	res := SearchIP("192.168.100.1")
	assert.Equal(t, res.ISP, "校园网")
	assert.Equal(t, res.IP, "192.168.100.1")
	res = SearchIP("2001:da8:8005::1")
	assert.Equal(t, res.Country, "中国")
	assert.Equal(t, res.IP, "2001:da8:8005::1")

	//非法的请求
	res = SearchIP("202..1")
	assert.Equal(t, res.IP, "202..1")
	assert.Equal(t, res.Country, "")
	res = SearchIP("202:da8:::::::")
	assert.Equal(t, res.IP, "202:da8:::::::")
	assert.Equal(t, res.Country, "")

}

func Benchmark_IpFindv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SearchIP("202.120.92.60")
	}
}

func Benchmark_IpFindv6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SearchIP("2001:da8:8005:abcd:1234::8888")
	}
}
