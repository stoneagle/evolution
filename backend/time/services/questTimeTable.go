package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestTimeTable struct {
	Pack ServicePackage
	structs.Service
}

func NewQuestTimeTable(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *QuestTimeTable {
	ret := QuestTimeTable{}
	ret.Init(engine, cache, log)
	return &ret
}
