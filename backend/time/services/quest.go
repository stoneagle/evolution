package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Quest struct {
	structs.Service
}

func NewQuest(engine *xorm.Engine, cache *redis.Client) *Quest {
	ret := Quest{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Quest) One(id int) (interface{}, error) {
	model := models.Quest{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Quest) Add(model models.Quest) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Quest) Update(id int, model models.Quest) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Quest) Delete(id int, model models.Quest) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Quest) List() (quests []models.Quest, err error) {
	quests = make([]models.Quest, 0)
	err = s.Engine.Find(&quests)
	return
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
