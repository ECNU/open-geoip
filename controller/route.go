package controller

import (
	"net/http"

	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/util"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"
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

	myip := r.Group("/")
	myip.Use(CORS())
	// 仅 ip 地址
	myip.GET("/myip", getMyIP)
	//json 结构化的 ip 地址
	myip.GET("/myip/format", getMyIPFormat)
	// 仅我的地理位置
	myip.GET("/mylocation", getMyLocation)
	//json 结构化的我的地理位置
	myip.GET("/mylocation/format", getMyLocationFormat)

	rest := r.Group("/api/v1")
	rest.Use(XAPICheckMidd)
	rest.GET("/network/ip", openGetIpApi)
	rest.DELETE("/ratelimit", clearRateLimit)
	rest.GET("/ratelimit", getRateLimit)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "主页"})
	})
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, g.VERSION)
	})
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		logger.Debug(context.Request.RequestURI, " - ", context.Request.Header.Get("Origin"))
		if util.InSliceStrFuzzy(context.Request.Header.Get("Origin"), g.Config().Http.CORS) {
			context.Writer.Header().Add("Access-Control-Allow-Origin", context.Request.Header.Get("Origin"))
			context.Writer.Header().Set("Access-Control-Max-Age", "86400")
			context.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
			context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, X-API-KEY, Authorization, x-requested-with")
			context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			context.Writer.Header().Set("Cache-Control", "no-cache")
		}
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
