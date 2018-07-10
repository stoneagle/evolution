package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestTarget struct {
	structs.Service
}

func NewQuestTarget(engine *xorm.Engine, cache *redis.Client) *QuestTarget {
	ret := QuestTarget{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *QuestTarget) One(id int) (interface{}, error) {
	model := models.QuestTarget{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *QuestTarget) Add(model models.QuestTarget) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *QuestTarget) Update(id int, model models.QuestTarget) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *QuestTarget) Delete(id int, model models.QuestTarget) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *QuestTarget) List() (questTargets []models.QuestTarget, err error) {
	questTargets = make([]models.QuestTarget, 0)
	err = s.Engine.Find(&questTargets)
	return
}

func (s *QuestTarget) ListWithCondition(questTarget *models.QuestTarget) (questTargets []models.QuestTarget, err error) {
	questTargets = make([]models.QuestTarget, 0)
	condition := questTarget.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&questTargets)
	if err != nil {
		return
	}
	return
}
