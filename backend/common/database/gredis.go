package database

import (
	"evolution/backend/common/config"

	"github.com/go-redis/redis"
)

var onceRedis *redis.Client

func SetRedis(redisConf config.RedisConf) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Host + ":" + redisConf.Port,
		Password: redisConf.Password,
		DB:       redisConf.Db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	onceRedis = client
}

func GetRedis() *redis.Client {
	return onceRedis
}
