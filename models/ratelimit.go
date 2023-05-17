package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ECNU/open-geoip/g"
	"github.com/gomodule/redigo/redis"
	"github.com/toolkits/pkg/logger"
)

const (
	Minute    = "minute"
	MinuteTTL = 60
	Hour      = "hour"
	HourTTL   = 3600
	Day       = "day"
	DayTTL    = 86400
	NameSpace = "open-geoip:rate:limit:"
)

func SetQueryRateLimit(enabled bool, clientIP string) error {
	if !enabled {
		return nil
	}
	redisPool := g.ConnectRedis()
	var count int
	var locked bool
	if g.Config().RateLimit.Day > 0 {
		count, locked = setQueryRateDay(redisPool, clientIP, g.Config().RateLimit.Day)
		if locked == true {
			err := fmt.Errorf("每日请求计数值达到 %d, %s 已经被锁定", count, clientIP)
			return err
		}
	}
	if g.Config().RateLimit.Hour > 0 {
		count, locked = setQueryRateHour(redisPool, clientIP, g.Config().RateLimit.Hour)
		if locked == true {
			err := fmt.Errorf("每小时请求计数值达到 %d, %s 已经被锁定", count, clientIP)
			return err
		}
	}
	if g.Config().RateLimit.Minute > 0 {
		count, locked = setQueryRateMinute(redisPool, clientIP, g.Config().RateLimit.Minute)
		if locked == true {
			err := fmt.Errorf("每分钟请求计数值达到 %d, %s 已经被锁定", count, clientIP)
			return err
		}
	}
	return nil
}

func ClearRateLimit(clientIP string) (err error) {
	conn := g.ConnectRedis().Get()
	minKey := NameSpace + Minute + ":" + clientIP
	_, err = conn.Do("DEL", minKey)
	if err != nil {
		return
	}
	hourKey := NameSpace + Hour + ":" + clientIP
	_, err = conn.Do("DEL", hourKey)
	if err != nil {
		return
	}
	dayKey := NameSpace + Day + ":" + clientIP
	_, err = conn.Do("DEL", dayKey)
	if err != nil {
		return
	}
	return
}

type CurrentRateCount struct {
	MinuteCount RateCount `json:"minuteCount"`
	HourCount   RateCount `json:"hourCount"`
	DayCount    RateCount `json:"dayCount"`
}

type RateCount struct {
	Count  int `json:"count"`
	Expire int `json:"expire"`
}

func GetCurrentRateCount(clientIP string) (res CurrentRateCount, err error) {
	redisPool := g.ConnectRedis()
	res.MinuteCount.Count, res.MinuteCount.Expire, err = getCurrentRateCount(redisPool, Minute, clientIP)
	if err != nil {
		return
	}
	res.HourCount.Count, res.HourCount.Expire, err = getCurrentRateCount(redisPool, Hour, clientIP)
	if err != nil {
		return
	}
	res.DayCount.Count, res.DayCount.Expire, err = getCurrentRateCount(redisPool, Day, clientIP)
	if err != nil {
		return
	}
	return
}

func getCurrentRateCount(redisPool *redis.Pool, mode, clientIP string) (count, expire int, err error) {
	if mode != Minute && mode != Hour && mode != Day {
		err = errors.New("不是合法的 mode")
		return
	}
	conn := redisPool.Get()
	defer conn.Close()
	redisKey := NameSpace + mode + ":" + clientIP
	res, err := conn.Do("GET", redisKey)
	if err != nil {
		return
	}
	if res == nil {
		return
	}
	value := string(res.([]byte))
	count, err = strconv.Atoi(value)
	if err != nil {
		logger.Error(err)
	}
	ttl, err := conn.Do("TTL", redisKey)
	if err != nil {
		return
	}
	expire = int(ttl.(int64))
	return
}

func setQueryRateMinute(redisPool *redis.Pool, clientIP string, limit int) (count int, locked bool) {
	conn := redisPool.Get()
	defer conn.Close()
	redisKey := NameSpace + Minute + ":" + clientIP

	// 使用 INCR 命令对键进行原子递增操作，并返回当前的计数值
	count, err := redis.Int(conn.Do("INCR", redisKey))
	if err != nil {
		logger.Error(err)
		return
	}
	// 如果这个键是第一次被创建，就使用 EXPIRE 命令为其设置过期时间为 60 秒
	if count == 1 {
		_, err := conn.Do("EXPIRE", redisKey, MinuteTTL)
		if err != nil {
			logger.Error(err)
			return
		}
	}
	// 根据计数值和限制值进行比较，如果超过了限制值，就返回 locked = true，否则返回 locked = false
	if count > limit {
		locked = true
	}
	return
}

func setQueryRateHour(redisPool *redis.Pool, clientIP string, limit int) (count int, locked bool) {
	conn := redisPool.Get()
	defer conn.Close()
	redisKey := NameSpace + Hour + ":" + clientIP

	// 使用 INCR 命令对键进行原子递增操作，并返回当前的计数值
	count, err := redis.Int(conn.Do("INCR", redisKey))
	if err != nil {
		logger.Error(err)
		return
	}
	// 如果这个键是第一次被创建，就使用 EXPIRE 命令为其设置过期时间为 3600 秒
	if count == 1 {
		_, err := conn.Do("EXPIRE", redisKey, HourTTL)
		if err != nil {
			logger.Error(err)
			return
		}
	}
	// 根据计数值和限制值进行比较，如果超过了限制值，就返回 locked = true，否则返回 locked = false
	if count > limit {
		locked = true
	}
	return
}

func setQueryRateDay(redisPool *redis.Pool, clientIP string, limit int) (count int, locked bool) {
	conn := redisPool.Get()
	defer conn.Close()
	redisKey := NameSpace + Day + ":" + clientIP

	// 使用 INCR 命令对键进行原子递增操作，并返回当前的计数值
	count, err := redis.Int(conn.Do("INCR", redisKey))
	if err != nil {
		logger.Error(err)
		return
	}
	// 如果这个键是第一次被创建，就使用 EXPIRE 命令为其设置过期时间为 3600 秒
	if count == 1 {
		_, err := conn.Do("EXPIRE", redisKey, DayTTL)
		if err != nil {
			logger.Error(err)
			return
		}
	}
	// 根据计数值和限制值进行比较，如果超过了限制值，就返回 locked = true，否则返回 locked = false
	if count > limit {
		locked = true
	}
	return
}
