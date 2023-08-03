package controller

import (
	"net/http"

	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
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
	if g.Config().SSO.Enabled {
		store, err := redis.NewStore(10, "tcp", g.Config().Redis.Dsn, g.Config().Redis.Password, []byte("open-geoip"))
		if err != nil {
			panic(err)
		}
		store.Options(sessions.Options{
			Path:     g.Config().Http.SessionOptions.Path,
			Domain:   g.Config().Http.SessionOptions.Domain,
			MaxAge:   g.Config().Http.SessionOptions.MaxAge,
			Secure:   g.Config().Http.SessionOptions.Secure,
			HttpOnly: g.Config().Http.SessionOptions.HttpOnly,
		})
		r.Use(sessions.Sessions("mysession", store))
	}

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, g.VERSION)
	})
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
	r.GET("/", index)

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

	// sso认证回调
	sso := r.Group("/sso")
	sso.Use(NoCache())
	// logout
	sso.GET("/logout", ssoLogout)
	// oauth
	sso.GET("/callback/oauth2", OauthAuth)
}

func index(c *gin.Context) {
	if g.Config().SSO.Enabled {
		var username, nickname string
		session := sessions.Default(c)
		u := session.Get("username")
		n := session.Get("nickname")
		if u != nil {
			username = u.(string)
		}
		if n != nil {
			nickname = n.(string)
		}
		if g.Config().Oauth2.Enabled {
			authCodeURL := g.Oauth2Config.AuthCodeURL(util.RandStringRunes(16))
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "主页", "sso": g.Config().SSO, "oauth2": g.Config().Oauth2, "authCodeURL": authCodeURL, "username": username, "nickname": nickname})
			return
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "主页"})
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

func NoCache() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Add("Cache-Control", "no-store")
		context.Writer.Header().Add("Pragma", "no-cache")
		context.Next()
	}
}
