package g

import (
	"github.com/gocarina/gocsv"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/toolkits/file"
	"log"
	"net"
	"os"
)

type Ips struct {
	IPSubnet       string `csv:"ip_subnet"`      //IP地址
	Continent      string `csv:"continent"`      //州
	Country        string `csv:"country"`        //国家
	Province       string `csv:"province"`       //省
	City           string `csv:"city"`           //城市
	District       string `csv:"district"`       //区县(行政区）
	ISP            string `csv:"isp"`            //运营商
	AreaCode       string `csv:"areaCode"`       //行政区划代码（国内）
	CountryEnglish string `csv:"countryEnglish"` //国家英文名
	CountryCode    string `csv:"countryCode"`    //国家英文代码
	Longitude      string `csv:"longitude"`      //经度
	Latitude       string `csv:"latitude"`       //纬度
}

func InitInternalDB(csvFile string) {

	if !file.IsExist(csvFile) {
		log.Fatalln("config file:", csvFile, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ips, err := loadDBFIle(csvFile)
	if err != nil {
		log.Println("cannot init internal db:", err)
		os.Exit(1)
	}
	//fmt.Printf("dasdasdasda\n")
	err = saveMMDB(ips)
	if err != nil {
		log.Println("cannot init internal db:", err)
		os.Exit(1)
	}

}

func loadDBFIle(csvFile string) (ips []*Ips, err error) {

	DBFile, err := os.OpenFile(csvFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer DBFile.Close()

	err = gocsv.UnmarshalFile(DBFile, &ips)

	if err != nil {
		return
	}

	return
}

func saveMMDB(ips []*Ips) (err error) {

	writer, err := mmdbwriter.New(mmdbwriter.Options{DatabaseType: "GeoIP2-City", IncludeReservedNetworks: true})

	if err != nil {
		return
	}

	for _, ip := range ips {

		_, sreNet, err := net.ParseCIDR(ip.IPSubnet)

		if err != nil {
			return err
		}

		sreData := mmdbtype.Map{
			"subdivisions": mmdbtype.Slice{
				mmdbtype.Map{
					"names": mmdbtype.Map{
						"zh-CN": mmdbtype.String(ip.Province),
					},
				},
			},
			"city": mmdbtype.Map{
				//"code":       mmdbtype.String("AS"),
				//"geoname_id": mmdbtype.Uint64(1808926),
				"names": mmdbtype.Map{
					//"de":    mmdbtype.String("China"),
					//"en":    mmdbtype.String("China"),
					//"es":    mmdbtype.String("China"),
					//"fr":    mmdbtype.String("Chine"),
					//"ja":    mmdbtype.String("中国"),
					//"pt-BR": mmdbtype.String("China"),
					//"ru":    mmdbtype.String("Китай"),
					"zh-CN": mmdbtype.String(ip.City),
				},
			},
			"continent": mmdbtype.Map{
				"names": mmdbtype.Map{
					"zh-CN": mmdbtype.String(ip.Continent),
				},
			},
			"country": mmdbtype.Map{
				"names": mmdbtype.Map{
					"zh-CN": mmdbtype.String(ip.Country),
					"en":    mmdbtype.String(ip.CountryEnglish),
				},
				"iso_code": mmdbtype.String(ip.CountryCode),
			},
			"registered_country": mmdbtype.Map{
				"names": mmdbtype.Map{
					"zh-CN": mmdbtype.String(ip.Country),
				},
			},
			"internal": mmdbtype.Map{
				"isp": mmdbtype.Map{
					"zh-CN": mmdbtype.String(ip.ISP),
				},
				"district": mmdbtype.Map{
					"zh-CN": mmdbtype.String(ip.District),
				},
				"areaCode": mmdbtype.String(ip.AreaCode),
			},
		}

		if err := writer.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
			return err
		}

	}

	//
	fh, err := os.Create("internal.mmdb")
	if err != nil {
		return
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		return
	}
	return
}
