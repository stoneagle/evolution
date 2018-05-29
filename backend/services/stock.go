package services

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Stock struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewStock(engine *xorm.Engine, cache *redis.Client) *Stock {
	return &Stock{
		engine: engine,
		cache:  cache,
	}
}
