package services

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Basic struct {
	Engine *xorm.Engine
	Cache  *redis.Client
}

type General interface {
	One(int) (interface{}, error)
}
