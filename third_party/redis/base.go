package redis

import (
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

var pool *redis.Pool

func init() {
	redisConfig := GetConfig()
	duration, _ := time.ParseDuration(redisConfig.Timeout)
	pool = &redis.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		MaxActive:   redisConfig.MaxActive,
		IdleTimeout: duration,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConfig.Conn, redis.DialDatabase(redisConfig.DbNum), redis.DialPassword(redisConfig.Password))
		},
	}
}

func Set(key string, value interface{}, expire int) {
	conn := pool.Get()
	conn.Do("Set", key, value)
	conn.Do("expire", key, expire)
}

func Get(key string) string {
	conn := pool.Get()
	res, err := redis.String(conn.Do("Get", key))
	if err != nil {
		logs.Error(err)
		return ""
	}
	return res
}

func Delete(key string) {
	conn := pool.Get()
	conn.Do("Del", key)
}
