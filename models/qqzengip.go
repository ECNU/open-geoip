package models

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type IpSearch struct {
	prefStart [256]uint32
	prefEnd   [256]uint32
	endArr    []uint32
	addrArr   []string
}

var instance *IpSearch
var once sync.Once

func GetInstance() *IpSearch {
	once.Do(func() {
		instance = &IpSearch{}
		var err error
		instance, err = LoadDat("./qqzeng-ip-3.0-ultimate.dat")
		if err != nil {
			log.Fatal("the IP Dat loaded failed!")
		}
	})
	return instance
}

func LoadDat(file string) (*IpSearch, error) {
	p := IpSearch{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	for k := 0; k < 256; k++ {
		i := k*8 + 4
		p.prefStart[k] = ReadLittleEndian32(data[i], data[i+1], data[i+2], data[i+3])
		p.prefEnd[k] = ReadLittleEndian32(data[i+4], data[i+5], data[i+6], data[i+7])
	}

	RecordSize := int(ReadLittleEndian32(data[0], data[1], data[2], data[3]))

	p.endArr = make([]uint32, RecordSize)
	p.addrArr = make([]string, RecordSize)
	for i := 0; i < RecordSize; i++ {
		j := 2052 + (i * 8)
		endipnum := ReadLittleEndian32(data[j], data[1+j], data[2+j], data[3+j])
		offset := ReadLittleEndian24(data[4+j], data[5+j], data[6+j])
		length := uint32(data[7+j])
		p.endArr[i] = endipnum
		p.addrArr[i] = string(data[offset:int(offset+length)])
	}
	return &p, err

}

func (p *IpSearch) Get(ip string) (ipGeo IpGeo) {
	intIP, err := ip2Long(ip)
	if err != nil {
		return
	}

	ips := strings.Split(ip, ".")
	x, _ := strconv.Atoi(ips[0])
	prefix := uint32(x)

	low := p.prefStart[prefix]
	high := p.prefEnd[prefix]

	var cur uint32
	if low == high {
		cur = low
	} else {
		cur = p.binarySearch(low, high, intIP)
	}

	ipSlice := strings.Split(p.addrArr[cur], "|")
	if len(ipSlice) != 11 {
		ipGeo.IP = ip
		return
	}

	ipGeo = IpGeo{
		IP:             ip,
		Continent:      ipSlice[0],
		Country:        ipSlice[1],
		Province:       ipSlice[2],
		City:           ipSlice[3],
		District:       ipSlice[4],
		ISP:            ipSlice[5],
		AreaCode:       ipSlice[6],
		CountryEnglish: ipSlice[7],
		CountryCode:    ipSlice[8],
		Longitude:      ipSlice[9],
		Latitude:       ipSlice[10],
	}
	return
}

func (p *IpSearch) binarySearch(low uint32, high uint32, k uint32) uint32 {
	var M uint32 = 0
	for low <= high {
		mid := (low + high) / 2
		endipNum := p.endArr[mid]
		if endipNum >= k {
			M = mid
			if mid == 0 {
				break
			}
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return M
}

func ip2Long(ipstr string) (uint32, error) {
	ip := net.ParseIP(ipstr)
	ip = ip.To4()
	if ip == nil {
		err := errors.New("qqzengip only support ipv4")
		return 0, err
	}
	return binary.BigEndian.Uint32(ip), nil
}

func ReadLittleEndian32(a, b, c, d byte) uint32 {
	return (uint32(a) & 0xFF) | ((uint32(b) << 8) & 0xFF00) | ((uint32(c) << 16) & 0xFF0000) | ((uint32(d) << 24) & 0xFF000000)
}

func ReadLittleEndian24(a, b, c byte) uint32 {
	return (uint32(a) & 0xFF) | ((uint32(b) << 8) & 0xFF00) | ((uint32(c) << 16) & 0xFF0000)
}
