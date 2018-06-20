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
	itemsJoin := make([]models.MapPoolItemJoin, 0)
	sql := s.engine.Unscoped().Table("map_pool_item").Join("INNER", "pool", "pool.id = map_pool_item.pool_id").Join("INNER", "item", "item.id = map_pool_item.item_id")

	sql = sql.Where("pool.id = ?", id)
	err = sql.Find(&itemsJoin)
	if err != nil {
		return
	}

	if len(itemsJoin) == 0 {
		return models.Pool{}, nil
	}
	return s.formatJoin(itemsJoin)[0], err
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

func (s *Pool) DeleteItems(pool *models.Pool) (err error) {
	session := s.engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	if len(pool.Item) > 0 {
		for _, item := range pool.Item {
			mapJoin := models.MapPoolItem{
				PoolId: pool.Id,
				ItemId: item.Id,
			}
			_, err = session.Delete(&mapJoin)
			if err != nil {
				session.Rollback()
				return
			}
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return nil
}

func (s *Pool) AddItems(pool *models.Pool) (err error) {
	session := s.engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	if len(pool.Item) > 0 {
		for _, item := range pool.Item {
			mapJoin := models.MapPoolItem{
				PoolId: pool.Id,
				ItemId: item.Id,
			}
			has, err := session.Get(&mapJoin)
			if err != nil {
				session.Rollback()
				return err
			}
			if !has {
				_, err = session.Insert(&mapJoin)
				if err != nil {
					session.Rollback()
					return err
				}
			}
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return nil
}

func (s *Pool) formatJoin(itemsJoin []models.MapPoolItemJoin) (pools []models.Pool) {
	poolsMap := make(map[int]models.Pool)
	for _, one := range itemsJoin {
		if target, ok := poolsMap[one.Pool.Id]; !ok {
			pool := one.Pool
			pool.Item = make([]models.Item, 0)
			item := one.Item
			if item.Id != 0 {
				pool.Item = append(pool.Item, item)
			}
			poolsMap[one.Pool.Id] = pool
		} else if one.Pool.Id != 0 {
			item := one.Item
			target.Item = append(target.Item, item)
			poolsMap[one.Pool.Id] = target
		}
	}

	pools = make([]models.Pool, 0)
	for _, one := range poolsMap {
		pools = append(pools, one)
	}
	return
}
