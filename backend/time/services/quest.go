package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

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

func (s *Quest) ListWithCondition(quest *models.Quest) (quests []models.Quest, err error) {
	quests = make([]models.Quest, 0)
	condition := quest.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&quests)
	if err != nil {
		return
	}
	return
}
