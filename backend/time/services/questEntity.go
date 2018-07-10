package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestEntity struct {
	structs.Service
}

func NewQuestEntity(engine *xorm.Engine, cache *redis.Client) *QuestEntity {
	ret := QuestEntity{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *QuestEntity) One(id int) (interface{}, error) {
	model := models.QuestEntity{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *QuestEntity) Add(model models.QuestEntity) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *QuestEntity) Update(id int, model models.QuestEntity) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *QuestEntity) Delete(id int, model models.QuestEntity) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *QuestEntity) List() (questEntitys []models.QuestEntity, err error) {
	questEntitys = make([]models.QuestEntity, 0)
	err = s.Engine.Find(&questEntitys)
	return
}

func (s *QuestEntity) ListWithCondition(questEntity *models.QuestEntity) (questEntitys []models.QuestEntity, err error) {
	questEntitys = make([]models.QuestEntity, 0)
	condition := questEntity.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&questEntitys)
	if err != nil {
		return
	}
	return
}
