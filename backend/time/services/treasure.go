package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Treasure struct {
	structs.Service
}

func NewTreasure(engine *xorm.Engine, cache *redis.Client) *Treasure {
	ret := Treasure{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Treasure) One(id int) (interface{}, error) {
	model := models.Treasure{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Treasure) Add(model models.Treasure) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Treasure) Update(id int, model models.Treasure) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Treasure) Delete(id int, model models.Treasure) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Treasure) List() (treasures []models.Treasure, err error) {
	treasures = make([]models.Treasure, 0)
	err = s.Engine.Find(&treasures)
	return
}

func (s *Treasure) ListWithCondition(treasure *models.Treasure) (treasures []models.Treasure, err error) {
	treasures = make([]models.Treasure, 0)
	condition := treasure.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&treasures)
	if err != nil {
		return
	}
	return
}
