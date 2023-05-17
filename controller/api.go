package controller

import (
	"net/http"

	"github.com/ECNU/open-geoip/g"

	"github.com/ECNU/open-geoip/models"

	"github.com/gin-gonic/gin"
)

func geoIpApi(c *gin.Context) {
	ipAddr := c.Query("ip")
	if !models.CheckIPValid(ipAddr) {
		c.String(http.StatusOK, "不是合法的IP地址")
		return
	}
	if err := models.SetQueryRateLimit(g.Config().RateLimit.Enabled, c.ClientIP()); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, models.SearchIP(ipAddr).ToString())
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
	c.String(http.StatusOK, models.SearchIP(c.ClientIP()).ToString())
}

func getMyLocationFormat(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessRes(models.SearchIP(c.ClientIP())))
}

func openGetIpApi(c *gin.Context) {
	ipAddr := c.Query("ip")
	if !models.CheckIPValid(ipAddr) {
		c.JSON(http.StatusOK, ErrorRes(ParamValueError, "不是合法的IP地址"))
		return
	}
	res := models.SearchIP(ipAddr)
	c.JSON(http.StatusOK, SuccessRes(res))
}
