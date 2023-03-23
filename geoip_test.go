package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ECNU/open-geoip/controller"
	"github.com/ECNU/open-geoip/g"
	"github.com/ECNU/open-geoip/models"
)

func init() {
	g.ParseConfig("cfg.json.test")
	err := models.InitReader()
	if err != nil {
		log.Fatalf("load geo db failed, %v", err)
	}
}

func TestIndex(t *testing.T) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 创建一个测试请求
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// 创建一个响应记录器
	w := httptest.NewRecorder()
	// 调用测试服务器的处理函数
	httpServer.Handler.ServeHTTP(w, req)
	// 检查响应状态码是否为 200
	assert.Equal(t, 200, w.Code)
	// 检查响应内容是否包含 open-geoip 地址
	assert.Contains(t, w.Body.String(), "open-geoip")
}

func TestSeachAPI(t *testing.T) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 创建一个测试请求
	req, err := http.NewRequest("GET", "/ip?ip=202.120.92.60", nil)
	if err != nil {
		t.Fatal(err)
	}
	// 创建一个响应记录器
	w := httptest.NewRecorder()
	// 调用测试服务器的处理函数
	httpServer.Handler.ServeHTTP(w, req)
	// 检查响应状态码是否为 200
	assert.Equal(t, 200, w.Code)
	// 检查响应内容是否包含 IP 地址
	assert.Contains(t, w.Body.String(), "中国")
}

func TestOpenAPI(t *testing.T) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 创建一个测试请求，带有 X-API-KEY 头
	req, err := http.NewRequest("GET", "/api/v1/network/ip?ip=202.120.92.60", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-API-KEY", "this-is-key")
	// 创建一个响应记录器
	w := httptest.NewRecorder()
	// 调用测试服务器的处理函数
	httpServer.Handler.ServeHTTP(w, req)
	// 检查响应状态码是否为 200
	assert.Equal(t, 200, w.Code)
	// 检查响应内容是否包含 IP 地址
	assert.Contains(t, w.Body.String(), "中国")
}

func BenchmarkIndex(b *testing.B) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 重置计时器
	b.ResetTimer()
	// 循环执行 b.N 次测试
	for i := 0; i < b.N; i++ {
		// 创建一个测试请求
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			b.Fatal(err)
		}
		// 创建一个响应记录器
		w := httptest.NewRecorder()
		// 调用测试服务器的处理函数
		httpServer.Handler.ServeHTTP(w, req)
	}
}

func BenchmarkSeachAPIForIPv4(b *testing.B) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 重置计时器
	b.ResetTimer()
	// 循环执行 b.N 次测试
	for i := 0; i < b.N; i++ {
		// 创建一个测试请求
		req, err := http.NewRequest("GET", "/ip?202.120.92.60", nil)
		if err != nil {
			b.Fatal(err)
		}
		// 创建一个响应记录器
		w := httptest.NewRecorder()
		// 调用测试服务器的处理函数
		httpServer.Handler.ServeHTTP(w, req)
	}
}

func BenchmarkSeachAPIForIPv6(b *testing.B) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 重置计时器
	b.ResetTimer()
	// 循环执行 b.N 次测试
	for i := 0; i < b.N; i++ {
		// 创建一个测试请求
		req, err := http.NewRequest("GET", "/ip?2001:da8:8005::1", nil)
		if err != nil {
			b.Fatal(err)
		}
		// 创建一个响应记录器
		w := httptest.NewRecorder()
		// 调用测试服务器的处理函数
		httpServer.Handler.ServeHTTP(w, req)
	}
}

func BenchmarkOpenAPIForIPv4(b *testing.B) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 重置计时器
	b.ResetTimer()
	// 循环执行 b.N 次测试
	for i := 0; i < b.N; i++ {
		// 创建一个测试请求
		req, err := http.NewRequest("GET", "/api/v1/network/ip?202.120.92.60", nil)
		if err != nil {
			b.Fatal(err)
		}
		req.Header.Set("X-API-KEY", "this-is-key")
		// 创建一个响应记录器
		w := httptest.NewRecorder()
		// 调用测试服务器的处理函数
		httpServer.Handler.ServeHTTP(w, req)
	}
}

func BenchmarkOpenAPIForIPv6(b *testing.B) {
	// 创建一个测试服务器
	httpServer := controller.InitGin(":8080")
	// 重置计时器
	b.ResetTimer()
	// 循环执行 b.N 次测试
	for i := 0; i < b.N; i++ {
		// 创建一个测试请求
		req, err := http.NewRequest("GET", "/api/v1/network/ip?2001:da8:8005::1", nil)
		if err != nil {
			b.Fatal(err)
		}
		req.Header.Set("X-API-KEY", "this-is-key")
		// 创建一个响应记录器
		w := httptest.NewRecorder()
		// 调用测试服务器的处理函数
		httpServer.Handler.ServeHTTP(w, req)
	}
}
