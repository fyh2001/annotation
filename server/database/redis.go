package database

import (
	"github.com/redis/go-redis/v9"
	"server/config"
)

var RedisClient *redis.Client

func InitRedis() {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Global.Redis.Host + ":" + config.Global.Redis.Port,
		Password: config.Global.Redis.Password, // 没有密码，默认值
		DB:       config.Global.Redis.Db,       // 默认DB 0
		//DialTimeout:  10 * time.Second,
		//ReadTimeout:  30 * time.Second,
		//WriteTimeout: 30 * time.Second,
		//PoolSize:     10,
		//PoolTimeout:  30 * time.Second,
	})
}

func GetRedis() *redis.Client {
	return RedisClient
}
