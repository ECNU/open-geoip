package models

import (
	"fmt"

	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func Test_RateLimit(t *testing.T) {
	// 启动一个 Redis 服务器
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	clientIP := "192.168.0.1" // 模拟客户端 IP
	limit := 10               // 设置每分钟的请求限制为 10 次

	for i := 0; i < 15; i++ {
		// 创建一个 Redis 连接
		conn, err := redis.Dial("tcp", s.Addr())
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		rc := &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return conn, nil
			},
		}
		count, locked := setQueryRateMinute(rc, clientIP, limit)
		fmt.Printf("第 %d 次请求，计数值为 %d，是否被锁定：%v\n", i+1, count, locked)
		if i < 10 {
			assert.Equal(t, locked, false)
		}
		if i >= 10 {
			assert.Equal(t, locked, true)
		}
	}
}
