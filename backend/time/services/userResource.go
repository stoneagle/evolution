package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type UserResource struct {
	structs.Service
}

func NewUserResource(engine *xorm.Engine, cache *redis.Client) *UserResource {
	ret := UserResource{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *UserResource) One(id int) (interface{}, error) {
	model := models.UserResource{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *UserResource) Add(model models.UserResource) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *UserResource) Update(id int, model models.UserResource) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *UserResource) Delete(id int, model models.UserResource) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *UserResource) List() (userResources []models.UserResource, err error) {
	userResources = make([]models.UserResource, 0)
	err = s.Engine.Find(&userResources)
	return
}

func (s *UserResource) ListWithCondition(resource *models.UserResource) (userResources []models.UserResource, err error) {
	userResources = make([]models.UserResource, 0)
	condition := resource.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&userResources)
	if err != nil {
		return
	}
	return
}
