package controller

import (
	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OauthAuth(c *gin.Context) {
	code := c.Query("code")
	oauthToken, err := util.OauthAuthReq(code)
	if err != nil {
		c.String(http.StatusInternalServerError, "内部错误")

	}

	userInfo, err := util.OauthUserInfo(oauthToken.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "内部错误")
	}

	c.SetCookie("nickname", userInfo.Nickname, g.Config().Oauth.AuthExpire, "/", "", false, true)
	c.SetCookie("username", userInfo.Username, g.Config().Oauth.AuthExpire, "/", "", false, true)
	c.Redirect(http.StatusMovedPermanently, "/")
}
