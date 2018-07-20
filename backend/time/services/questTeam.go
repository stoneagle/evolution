package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestTeam struct {
	Pack ServicePackage
	structs.Service
}

func NewQuestTeam(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *QuestTeam {
	ret := QuestTeam{}
	ret.Init(engine, cache, log)
	return &ret
}
