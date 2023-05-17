package g

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var redisConnPool *redis.Pool

func ConnectRedis() *redis.Pool {
	return redisConnPool
}

func CloseRedis() (err error) {
	err = redisConnPool.Close()
	return
}

func InitRedisConnPool() {
	dsn := Config().Redis.Dsn
	maxIdle := Config().Redis.MaxIdle
	idleTimeout := 240 * time.Second
	connTimeout := time.Duration(Config().Redis.ConnTimeout) * time.Second
	readTimeout := time.Duration(Config().Redis.ReadTimeout) * time.Second
	writeTimeout := time.Duration(Config().Redis.WriteTimeout) * time.Second
	password := Config().Redis.Password

	redisConnPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", dsn,
				redis.DialConnectTimeout(connTimeout),
				redis.DialReadTimeout(readTimeout),
				redis.DialWriteTimeout(writeTimeout),
				redis.DialPassword(password))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: pingRedis,
	}
}

func pingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Println("[ERROR] ping redis fail", err)
	}
	return err
}
