package controller

import (
	"net/http"
	"strings"

	"github.com/ECNU/open-geoip/g"
	"github.com/gin-contrib/sessions"

	"github.com/ECNU/open-geoip/models"

	"github.com/gin-gonic/gin"
)

func geoIpApi(c *gin.Context) {
	isAuth := false
	ipAddr := c.Query("ip")
	// 去掉左右空格
	ipAddr = strings.TrimSpace(ipAddr)
	if !models.CheckIPValid(ipAddr) {
		c.String(http.StatusOK, "不是合法的IP地址")
		return
	}
	if err := models.SetQueryRateLimit(g.Config().RateLimit.Enabled, c.ClientIP()); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if g.Config().SSO.Enabled {
		session := sessions.Default(c)
		u := session.Get("username")
		if u != nil {
			isAuth = true
		}
	}

	c.String(http.StatusOK, models.SearchIP(ipAddr, false, isAuth).ToString())
}

func getMyIP(c *gin.Context) {
	c.String(http.StatusOK, c.ClientIP())
}

func getMyIPFormat(c *gin.Context) {
	res := map[string]string{
		"ip": c.ClientIP(),
	}
	c.JSON(http.StatusOK, SuccessRes(res))
}

func getMyLocation(c *gin.Context) {
	c.String(http.StatusOK, models.SearchIP(c.ClientIP(), true, false).ToString())
}

func getMyLocationFormat(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessRes(models.SearchIP(c.ClientIP(), true, false)))
}

func openGetIpApi(c *gin.Context) {
	ipAddr := c.Query("ip")
	if !models.CheckIPValid(ipAddr) {
		c.JSON(http.StatusOK, ErrorRes(ParamValueError, "不是合法的IP地址"))
		return
	}
	res := models.SearchIP(ipAddr, true, false)
	c.JSON(http.StatusOK, SuccessRes(res))
}
