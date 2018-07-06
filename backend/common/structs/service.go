package structs

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Service struct {
	Engine *xorm.Engine
	Cache  *redis.Client
}

type ServiceGeneral interface {
	One(int) (interface{}, error)
}
