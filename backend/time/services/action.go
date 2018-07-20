package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Action struct {
	Pack ServicePackage
	structs.Service
}

func NewAction(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Action {
	ret := Action{}
	ret.Init(engine, cache, log)
	return &ret
}
