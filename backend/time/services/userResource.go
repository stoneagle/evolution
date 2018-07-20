package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type UserResource struct {
	Pack ServicePackage
	structs.Service
}

func NewUserResource(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *UserResource {
	ret := UserResource{}
	ret.Init(engine, cache, log)
	return &ret
}
