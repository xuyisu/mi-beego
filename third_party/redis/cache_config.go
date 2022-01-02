package redis

import (
	beego "github.com/beego/beego/v2/server/web"
)

type RedisConfig struct {
	Key       string `json:"key"`
	Conn      string `json:"conn"`
	DbNum     int    `json:"dbNum"`
	Password  string `json:"password"`
	MaxIdle   int    `json:"maxIdle"`
	Timeout   string `json:"timeout"`
	MaxActive int    `json:"maxActive"`
}

func GetConfig() RedisConfig {
	key, _ := beego.AppConfig.String("redis_key")
	conn, _ := beego.AppConfig.String("redis_conn")
	dbNum, _ := beego.AppConfig.Int("redis_dbNum")
	password, _ := beego.AppConfig.String("redis_password")
	maxIdle, _ := beego.AppConfig.Int("redis_maxIdle")
	maxActive, _ := beego.AppConfig.Int("redis_maxActive")
	timeout, _ := beego.AppConfig.String("redis_timeout")
	return RedisConfig{
		Key:       key,
		Conn:      conn,
		DbNum:     dbNum,
		Password:  password,
		MaxIdle:   maxIdle,
		Timeout:   timeout,
		MaxActive: maxActive,
	}
}
