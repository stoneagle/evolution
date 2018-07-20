package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Quest struct {
	Pack ServicePackage
	structs.Service
}

func NewQuest(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Quest {
	ret := Quest{}
	ret.Init(engine, cache, log)
	return &ret
}
