package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Field struct {
	structs.Service
}

func NewField(engine *xorm.Engine, cache *redis.Client) *Field {
	ret := Field{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Field) One(id int) (interface{}, error) {
	model := models.Field{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Field) Add(model models.Field) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Field) Update(id int, model models.Field) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Field) Delete(id int, model models.Field) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Field) List() (fields []models.Field, err error) {
	fields = make([]models.Field, 0)
	err = s.Engine.Find(&fields)
	return
}

func (s *Field) Map() (fieldsMap map[int]string, err error) {
	fields := make([]models.Field, 0)
	err = s.Engine.Find(&fields)
	if err != nil {
		return
	}
	fieldsMap = map[int]string{}
	for _, one := range fields {
		fieldsMap[one.Id] = one.Name
	}
	return
}
