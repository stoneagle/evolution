package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type User struct {
	ServicePackage
	structs.Service
}

func NewUser(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *User {
	ret := User{}
	ret.Init(engine, cache, log)
	return &ret
}
