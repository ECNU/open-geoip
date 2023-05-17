package controller

import (
	"net/http"

	"github.com/ECNU/open-geoip/models"

	"github.com/gin-gonic/gin"
)

func getRateLimit(c *gin.Context) {
	clientIP := c.Query("clientip")
	currentRateLimit, err := models.GetCurrentRateCount(clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorRes(InternalAPIError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessRes(currentRateLimit))
}

func clearRateLimit(c *gin.Context) {
	clientIP := c.Query("clientip")
	if err := models.ClearRateLimit(clientIP); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorRes(InternalAPIError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessRes(nil))
}
