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

func (m *Pool) One(id int) (pool models.Pool, err error) {
	pool = models.Pool{}
	_, err = m.engine.Where("id = ?", id).Get(&pool)
	return
}

func (m *Pool) List() (pools []models.Pool, err error) {
	pools = make([]models.Pool, 0)
	err = m.engine.Find(&pools)
	return
}

func (m *Pool) UpdateByMap(id int, pool map[string]interface{}) (err error) {
	_, err = m.engine.Table(new(models.Pool)).Id(id).Update(&pool)
	return
}

func (m *Pool) Update(id int, pool *models.Pool) (err error) {
	_, err = m.engine.Id(id).Update(pool)
	return
}

func (m *Pool) Add(pool *models.Pool) (err error) {
	_, err = m.engine.Insert(pool)
	return
}

func (m *Pool) Delete(id int) (err error) {
	var pool models.Pool
	_, err = m.engine.Id(id).Get(&pool)
	if err == nil {
		_, err = m.engine.Id(id).Delete(&pool)
	}
	return
}
