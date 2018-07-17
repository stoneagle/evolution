package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Country struct {
	ServicePackage
	structs.Service
}

func NewCountry(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Country {
	ret := Country{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Country) Add(model models.Country) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Country) Update(id int, model models.Country) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Country) List() (countries []models.Country, err error) {
	countries = make([]models.Country, 0)
	err = s.Engine.Find(&countries)
	return
}
