package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

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

func (s *QuestResource) ListWithCondition(questResource *models.QuestResource) (questResources []models.QuestResource, err error) {
	questResources = make([]models.QuestResource, 0)
	condition := questResource.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&questResources)
	if err != nil {
		return
	}
	return
}
