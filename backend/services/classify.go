package services

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Classify struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewClassify(engine *xorm.Engine, cache *redis.Client) *Classify {
	return &Classify{
		engine: engine,
		cache:  cache,
	}
}
