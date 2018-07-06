package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Country struct {
	structs.Service
}

func NewCountry(engine *xorm.Engine, cache *redis.Client) *Country {
	ret := Country{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Country) One(id int) (interface{}, error) {
	model := models.Country{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Country) Add(model models.Country) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Country) Update(id int, model models.Country) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Country) Delete(id int, model models.Country) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Country) List() (countries []models.Country, err error) {
	countries = make([]models.Country, 0)
	err = s.Engine.Find(&countries)
	return
}
