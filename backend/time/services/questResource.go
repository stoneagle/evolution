package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestResource struct {
	structs.Service
}

func NewQuestResource(engine *xorm.Engine, cache *redis.Client) *QuestResource {
	ret := QuestResource{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *QuestResource) One(id int) (interface{}, error) {
	model := models.QuestResource{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *QuestResource) Add(model models.QuestResource) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *QuestResource) Update(id int, model models.QuestResource) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *QuestResource) Delete(id int, model models.QuestResource) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *QuestResource) List() (questResources []models.QuestResource, err error) {
	questResources = make([]models.QuestResource, 0)
	err = s.Engine.Find(&questResources)
	return
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
