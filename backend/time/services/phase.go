package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Phase struct {
	Pack ServicePackage
	structs.Service
}

func NewPhase(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Phase {
	ret := Phase{}
	ret.Init(engine, cache, log)
	return &ret
}
