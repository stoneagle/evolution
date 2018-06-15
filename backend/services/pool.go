package services

import (
	"quant/backend/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Pool struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewPool(engine *xorm.Engine, cache *redis.Client) *Pool {
	return &Pool{
		engine: engine,
		cache:  cache,
	}
}

func (s *Pool) One(id int) (pool models.Pool, err error) {
	pool = models.Pool{}
	_, err = s.engine.Where("id = ?", id).Get(&pool)
	return
}

func (s *Pool) List() (pools []models.Pool, err error) {
	pools = make([]models.Pool, 0)
	err = s.engine.Find(&pools)
	return
}

func (s *Pool) UpdateByMap(id int, pool map[string]interface{}) (err error) {
	_, err = s.engine.Table(new(models.Pool)).Id(id).Update(&pool)
	return
}

func (s *Pool) Update(id int, pool *models.Pool) (err error) {
	_, err = s.engine.Id(id).Update(pool)
	return
}

func (s *Pool) Add(pool *models.Pool) (err error) {
	_, err = s.engine.Insert(pool)
	return
}

func (s *Pool) Delete(id int) (err error) {
	var pool models.Pool
	_, err = s.engine.Id(id).Get(&pool)
	if err == nil {
		_, err = s.engine.Id(id).Delete(&pool)
	}
	return
}
