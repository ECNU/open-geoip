package util

import (
	"io/ioutil"
	"net/http"

	"github.com/ECNU/open-geoip/g"
	jsoniter "github.com/json-iterator/go"
)

type UserInfo struct {
	Username string `json:"Username"`
	Nickname string `json:"Nickname"`
}

// OauthUserInfo 获取用户信息
func OauthUserInfo(accessToken string) (userInfo UserInfo, err error) {
	//fmt.Printf("%v\n", accessToken)
	client := &http.Client{}

	// 创建一个 http.NewRequest 对象
	req, err := http.NewRequest("GET", g.Config().Oauth2.UserInfoAddr, nil)
	if err != nil {
		//panic(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		//panic(err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//panic(err)
		return
	}

	userInfo.Username = getUserinfoField(responseBody, g.Config().Oauth2.UserinfoIsArray, g.Config().Oauth2.UserinfoPrefix, g.Config().Oauth2.Attributes.Username)
	userInfo.Nickname = getUserinfoField(responseBody, g.Config().Oauth2.UserinfoIsArray, g.Config().Oauth2.UserinfoPrefix, g.Config().Oauth2.Attributes.Nickname)

	return

}

func getUserinfoField(input []byte, isArray bool, prefix, field string) string {
	if prefix == "" {
		if isArray {
			return jsoniter.Get(input, 0).Get(field).ToString()
		} else {
			return jsoniter.Get(input, field).ToString()
		}
	} else {
		if isArray {
			return jsoniter.Get(input, prefix, 0).Get(field).ToString()
		} else {
			return jsoniter.Get(input, prefix).Get(field).ToString()
		}
	}
}
