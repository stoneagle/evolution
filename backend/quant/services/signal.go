package services

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Signal struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewSignal(engine *xorm.Engine, cache *redis.Client) *Signal {
	return &Signal{
		engine: engine,
		cache:  cache,
	}
}
