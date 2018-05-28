package common

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

var onceRedis *redis.Client

var mapXormEngine map[string]*xorm.Engine = make(map[string]*xorm.Engine)

func SetEngine(key string, e *xorm.Engine) {
	mapXormEngine[key] = e
}

func GetEngine(key string) *xorm.Engine {
	e, ok := mapXormEngine[key]
	if !ok {
		panic("engine not exist")
	}
	return e
}

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
