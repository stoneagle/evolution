package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestResource struct {
	Pack ServicePackage
	structs.Service
}

func NewQuestResource(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *QuestResource {
	ret := QuestResource{}
	ret.Init(engine, cache, log)
	return &ret
}
