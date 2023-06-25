package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ECNU/open-geoip/g"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
)

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	Redirect    string `json:"redirect"`
	Msg         string `json:"msg"`
	AccessToken string `json:"accessToken"`
	Username    string `json:"Username"`
	Nickname    string `json:"Nickname"`
	Phone       string `yaml:"Phone"`
	Email       string `yaml:"Email"`
}

// OauthAuthReq 请求token
func OauthAuthReq(code string) (oauthToken OauthToken, err error) {

	client := &http.Client{}

	data := map[string]string{
		"code":          code,
		"redirect_uri":  g.Config().Oauth.RedirectURL,
		"grant_type":    "authorization_code",
		"client_id":     g.Config().Oauth.ClientId,
		"client_secret": g.Config().Oauth.ClientSecret,
	}

	requestBody, err := json.Marshal(data)
	if err != nil {
		//panic(err)
		return
	}
	req, err := http.NewRequest("POST", g.Config().Oauth.TokenAddr, bytes.NewBuffer(requestBody))
	if err != nil {
		//panic(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		//panic(err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		//panic(err)
		return
	}

	err = json.Unmarshal(responseBody, &oauthToken)
	if err != nil {
		//panic(err)
		return
	}
	fmt.Printf("%v\n", oauthToken.AccessToken)
	return

	//fmt.Printf("%v\n", oauthToken.AccessToken)
}

// OauthUserInfo 获取用户信息
func OauthUserInfo(accessToken string) (userInfo UserInfo, err error) {
	//fmt.Printf("%v\n", accessToken)
	client := &http.Client{}

	// 创建一个 http.NewRequest 对象
	req, err := http.NewRequest("GET", g.Config().Oauth.UserInfoAddr, nil)
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

	userInfo.Username = getUserinfoField(responseBody, g.Config().Oauth.UserinfoIsArray, g.Config().Oauth.UserinfoPrefix, g.Config().Oauth.Attributes.Username)
	userInfo.Nickname = getUserinfoField(responseBody, g.Config().Oauth.UserinfoIsArray, g.Config().Oauth.UserinfoPrefix, g.Config().Oauth.Attributes.Nickname)
	userInfo.Email = getUserinfoField(responseBody, g.Config().Oauth.UserinfoIsArray, g.Config().Oauth.UserinfoPrefix, g.Config().Oauth.Attributes.Email)
	userInfo.Phone = getUserinfoField(responseBody, g.Config().Oauth.UserinfoIsArray, g.Config().Oauth.UserinfoPrefix, g.Config().Oauth.Attributes.Phone)
	userInfo.AccessToken = accessToken

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
