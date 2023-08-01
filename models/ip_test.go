package models

import (
	"log"
	"testing"

	"github.com/ECNU/open-geoip/g"
	"github.com/toolkits/file"

	"github.com/stretchr/testify/assert"
)

func init() {
	g.ParseConfig("cfg.json.test")
	if file.IsExist("qqzeng-ip-3.0-ultimate.dat") {
		g.Config().DB.Qqzengip = "qqzeng-ip-3.0-ultimate.dat"
	}
	if file.IsExist("city.free.ipdb") {
		g.Config().DB.Ipdb = "city.free.ipdb"
	}
	err := InitReader()
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

	res, _ := GetIP("192.168.0.1", g.Config().Source, false, false)
	assert.Equal(t, res.Country, "中国")
	res, _ = GetIP("fd12:3456:789a:bcde::1", g.Config().Source, false, false)
	assert.Equal(t, res.Country, "中国")

	if file.IsExist("city.free.ipdb") {
		g.Config().Source.IPv4 = "ipdb"
		res, _ = GetIP("114.114.114.114", g.Config().Source, false, false)
		t.Log(res)
		assert.Equal(t, res.Country, "114DNS.COM")
	}
	if file.IsExist("qqzeng-ip-3.0-ultimate.dat") {
		g.Config().Source.IPv4 = "qqzengip"
		res, _ = GetIP("202.120.92.60", g.Config().Source, false, false)
		t.Log(res)
		assert.Equal(t, res.ISP, "教育网")
	}
	//非法的请求
	_, err := GetIP("201..1", g.Config().Source, false, false)
	assert.NotNil(t, err)
	_, err = GetIP("2001:da8:::::::::", g.Config().Source, false, false)
	assert.NotNil(t, err)
}

func TestSearchIP(t *testing.T) {
	g.Config().Source.IPv4 = "maxmind"
	res := SearchIP("192.168.0.1", false, false)
	assert.Equal(t, res.Province, "上海")
	assert.Equal(t, res.IP, "192.168.0.1")
	res = SearchIP("fd12:3456:789a:bcde::1", false, false)
	assert.Equal(t, res.Country, "中国")
	assert.Equal(t, res.IP, "fd12:3456:789a:bcde::1")

	//非法的请求
	res = SearchIP("202..1", false, false)
	assert.Equal(t, res.IP, "202..1")
	assert.Equal(t, res.Country, "")
	res = SearchIP("202:da8:::::::", false, false)
	assert.Equal(t, res.IP, "202:da8:::::::")
	assert.Equal(t, res.Country, "")
}

func Benchmark_maxmind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SearchIP("202.120.92.60", false, false)
	}
}

func Benchmark_Ipdb(b *testing.B) {
	if !file.IsExist("city.free.ipdb") {
		return
	}
	g.Config().Source.IPv4 = "ipdb"
	for i := 0; i < b.N; i++ {
		SearchIP("202.120.92.60", false, false)
	}
}

func Benchmark_qqzengip(b *testing.B) {
	if !file.IsExist("qqzeng-ip-3.0-ultimate.dat") {
		return
	}
	g.Config().Source.IPv4 = "qqzengip"
	for i := 0; i < b.N; i++ {
		SearchIP("202.120.92.60", false, false)
	}
}
