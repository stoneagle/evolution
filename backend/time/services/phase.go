package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Phase struct {
	structs.Service
}

func NewPhase(engine *xorm.Engine, cache *redis.Client) *Phase {
	ret := Phase{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Phase) One(id int) (interface{}, error) {
	model := models.Phase{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Phase) Add(model models.Phase) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Phase) Update(id int, model models.Phase) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Phase) Delete(id int, model models.Phase) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Phase) List() (phases []models.Phase, err error) {
	phases = make([]models.Phase, 0)
	err = s.Engine.Find(&phases)
	return
}

func (s *Phase) ListWithCondition(phase *models.Phase) (phases []models.Phase, err error) {
	phases = make([]models.Phase, 0)
	sql := s.Engine.Asc("level")
	condition := phase.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&phases)
	if err != nil {
		return
	}
	return
}
