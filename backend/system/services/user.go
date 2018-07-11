package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/system/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type User struct {
	structs.Service
}

func NewUser(engine *xorm.Engine, cache *redis.Client) *User {
	ret := User{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *User) One(id int) (interface{}, error) {
	model := models.User{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *User) OneByCondition(user *models.User) (models.User, error) {
	model := models.User{}
	condition := user.BuildCondition()
	_, err := s.Engine.Where(condition).Get(&model)
	return model, err
}

func (s *User) Add(model models.User) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *User) Update(id int, model models.User) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *User) Delete(id int, model models.User) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *User) List() (users []models.User, err error) {
	users = make([]models.User, 0)
	err = s.Engine.Find(&users)
	return
}

func (s *User) ListWithCondition(user *models.User) (users []models.User, err error) {
	users = make([]models.User, 0)
	condition := user.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&users)
	if err != nil {
		return
	}
	return
}
