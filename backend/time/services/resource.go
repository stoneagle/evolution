package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Resource struct {
	structs.Service
}

func NewResource(engine *xorm.Engine, cache *redis.Client) *Resource {
	ret := Resource{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Resource) One(id int) (interface{}, error) {
	model := models.Resource{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Resource) Add(model models.Resource) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Resource) Update(id int, model models.Resource) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Resource) Delete(id int, model models.Resource) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Resource) List() (resources []models.Resource, err error) {
	resources = make([]models.Resource, 0)
	err = s.Engine.Find(&resources)
	return
}

func (s *Resource) ListWithCondition(resource *models.Resource) (resources []models.Resource, err error) {
	resources = make([]models.Resource, 0)
	condition := resource.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&resources)
	if err != nil {
		return
	}
	return
}
