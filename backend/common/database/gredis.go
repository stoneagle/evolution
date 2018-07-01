package database

import (
	"github.com/go-redis/redis"
)

var onceRedis *redis.Client

func SetRedis(host, port, password string, db int) {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
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
