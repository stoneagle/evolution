package services

import (
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type EntityLife struct {
	Basic
}

func NewEntityLife(engine *xorm.Engine, cache *redis.Client) *EntityLife {
	ret := EntityLife{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *EntityLife) One(id int) (interface{}, error) {
	model := models.EntityLife{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *EntityLife) Add(model models.EntityLife) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *EntityLife) Update(id int, model models.EntityLife) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *EntityLife) Delete(id int, model models.EntityLife) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *EntityLife) List() (entityLifes []models.EntityLife, err error) {
	entityLifes = make([]models.EntityLife, 0)
	err = s.Engine.Find(&entityLifes)
	return
}
