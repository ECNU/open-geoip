package controller

import (
	"fmt"
	"net/http"
	"strings"

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
		username, _ := c.Cookie("username")
		nickname, _ := c.Cookie("nickname")
		oauthUrl := fmt.Sprintf("%s?scope=%s&redirect_uri=%s&response_type=code&client_id=%s", g.Config().Oauth.SsoAddr, strings.Join(g.Config().Oauth.Scopes, " "), g.Config().Oauth.RedirectURL, g.Config().Oauth.ClientId)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "主页", "sso": g.Config().SSO, "oauth": g.Config().Oauth, "oauthUrl": oauthUrl, "username": username, "nickname": nickname})
	})
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, g.VERSION)
	})

	r.GET("/logout", func(c *gin.Context) {

		username, _ := c.Cookie("username")
		nickname, _ := c.Cookie("nickname")

		c.SetCookie("nickname", nickname, -1, "/", "", false, true)
		c.SetCookie("username", username, -1, "/", "", false, true)
		c.Redirect(http.StatusMovedPermanently, g.Config().Oauth.LogoutAddr)
	})

	// sso认证回调
	ssoCallback := r.Group("/auth")
	// oauth
	ssoCallback.GET("/callback/oauth", OauthAuth)
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
