package services

import (
	"quant/backend/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Item struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewItem(engine *xorm.Engine, cache *redis.Client) *Item {
	return &Item{
		engine: engine,
		cache:  cache,
	}
}

func (s *Item) One(id int) (item models.Item, err error) {
	item = models.Item{}
	_, err = s.engine.Where("id = ?", id).Get(&item)
	return
}

func (s *Item) List() (items []models.Item, err error) {
	items = make([]models.Item, 0)
	err = s.engine.Find(&items)
	return
}
