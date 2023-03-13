package controller

import (
	"net/http"

	"github.com/ECNU/go-geoip/g"
	"github.com/gin-gonic/gin"
)

func InitGin(listen string) (httpServer *http.Server) {
	if g.Config().Logger.Level == "DEBUG" {
		gin.SetMode((gin.DebugMode))
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	if g.Config().Logger.Level == "DEBUG" {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.SetTrustedProxies(g.Config().Http.TrustProxy)

	Routes(r)

	httpServer = &http.Server{
		Addr:    g.Config().Http.Listen,
		Handler: r,
	}
	return
}

func Routes(r *gin.Engine) {
	r.LoadHTMLFiles("templates/index.html")
	r.Static("/assets", "assets")

	r.GET("/ip", geoIpApi)
	r.GET("/myip", getMyIP)
	//json 结构化的 ip 地址
	r.GET("/myip/format", getMyIPFormat)

	rest := r.Group("/api/v1")
	rest.Use(XAPICheckMidd)
	rest.GET("/network/ip", openGetIpApi)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "主页"})
	})
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, g.VERSION)
	})

}
