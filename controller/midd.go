package controller

import (
	"net/http"

	"github.com/ECNU/open-geoip/g"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/toolkits/pkg/logger"
)

// APIResult api 接口的数据结构
type APIResult struct {
	ErrCode   int64       `json:"errCode"`
	ErrMsg    string      `json:"errMsg"`
	RequestId string      `json:"requestId"`
	Data      interface{} `json:"data"`
}

const (
	Success          = 0
	ParamFormatError = 4001
	ParamValueError  = 4002
	ParamMissError   = 4003
	InternalAPIError = 5000
)

var codeMsg = map[int64]string{
	Success:          "success",
	ParamFormatError: "参数校验错误",
	ParamValueError:  "参数取值错误",
	ParamMissError:   "缺失参数",
	InternalAPIError: "服务器内部错误",
}

// ErrorRes 请求异常时的返回
func ErrorRes(code int64, msg string) (res APIResult) {
	res.ErrCode = code
	if msg == "" {
		res.ErrMsg = codeMsg[code]
	} else {
		res.ErrMsg = msg
	}
	res.RequestId = uuid.NewV4().String()
	return
}

// SuccessRes 请求正常时的返回
func SuccessRes(data interface{}) (apiResult APIResult) {
	apiResult.ErrCode = 0
	apiResult.ErrMsg = "success"
	apiResult.RequestId = uuid.NewV4().String()
	apiResult.Data = data
	return
}

// XAPICheckMidd 校验X-API-KEY，供API网关代理这个接口
func XAPICheckMidd(c *gin.Context) {
	key := c.Request.Header.Get("X-API-KEY")
	if !checkXApiKey(key) {
		logger.Warning(key, g.Config().Http.XAPIKey)
		c.JSON(http.StatusUnauthorized, ErrorRes(InternalAPIError, ""))
		c.Abort()
		return
	}
	c.Next()
}

func checkXApiKey(key string) bool {
	return key == g.Config().Http.XAPIKey
}
