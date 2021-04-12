package db

import (
	"fmt"
	"game_server/core/base"

	"github.com/go-redis/redis"
)

var RedisGame *redis.Client

func InitRedis() {
	RedisGame = redis.NewClient(&redis.Options{
		Addr:     base.Setting.Redis.Host + ":" + base.Setting.Redis.Port,
		Password: base.Setting.Redis.Password,
		DB:       base.Setting.Redis.DbName,
		PoolSize: base.Setting.Redis.PoolSize,
	})

	_, err := RedisGame.Set("ip_limits", "*", 0).Result()
	if err != nil {
		fmt.Printf("redis set key ip_limits failed")
	}
}
