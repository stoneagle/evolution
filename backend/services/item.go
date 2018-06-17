package services

import (
	"quant/backend/models"
	"time"

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

func (s *Item) BatchSave(classify models.Classify, itemBatch []models.Item) (err error) {
	session := s.engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	for _, one := range itemBatch {
		has, err := session.Get(&one)
		if err != nil {
			session.Rollback()
			return err
		}

		// save data
		if has {
			updateOne := models.Item{}
			updateOne.GeneralWithDeleted.UpdatedAt = time.Now()
			_, err = session.Id(one.Id).Update(&updateOne)
			if err != nil {
				session.Rollback()
				return err
			}
		} else {
			_, err = session.Insert(&one)
			if err != nil {
				session.Rollback()
				return err
			}
		}

		// save join
		join := models.MapClassifyItem{}
		join.ItemId = one.Id
		join.ClassifyId = classify.Id
		hasJoin, err := session.Get(&join)
		if err != nil {
			session.Rollback()
			return err
		}
		if !hasJoin {
			_, err = session.Insert(&join)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}

	// TODO 更新股票状态
	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return nil
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
