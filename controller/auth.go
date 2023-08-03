package controller

import (
	"context"
	"net/http"

	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func OauthAuth(c *gin.Context) {
	code := c.Request.FormValue("code")
	token, err := g.Oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "内部错误")
	}

	userInfo, err := util.OauthUserInfo(token.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "内部错误")
	}

	session := sessions.Default(c)
	session.Set("username", userInfo.Username)
	session.Set("nickname", userInfo.Nickname)
	session.Save()

	c.Redirect(http.StatusMovedPermanently, "/")
}

func ssoLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	c.Redirect(http.StatusMovedPermanently, g.Config().Oauth2.LogoutAddr)
}
