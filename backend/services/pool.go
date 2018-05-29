package services

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Pool struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewPool(engine *xorm.Engine, cache *redis.Client) *Pool {
	return &Pool{
		engine: engine,
		cache:  cache,
	}
}
