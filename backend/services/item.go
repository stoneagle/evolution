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
	itemsJoin := make([]models.MapClassifyItemJoin, 0)
	sql := s.engine.Unscoped().Table("map_classify_item").Join("INNER", "classify", "classify.id = map_classify_item.classify_id").Join("INNER", "item", "item.id = map_classify_item.item_id")

	sql = sql.Where("item.id = ?", id)
	err = sql.Find(&itemsJoin)
	if err != nil {
		return
	}

	if len(itemsJoin) == 0 {
		item = models.Item{}
		_, err = s.engine.Where("id = ?", id).Get(&item)
		return
	}
	return s.formatJoin(itemsJoin)[0], err
}

func (s *Item) List(item *models.Item) (items []models.Item, err error) {
	itemsJoin := make([]models.MapClassifyItemJoin, 0)
	// be careful with the sequense of join table and struct relationship
	sql := s.engine.Unscoped().Table("map_classify_item").Join("INNER", "classify", "classify.id = map_classify_item.classify_id").Join("INNER", "item", "item.id = map_classify_item.item_id")

	condition := item.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&itemsJoin)
	if err != nil {
		return
	}

	return s.formatJoin(itemsJoin), err
}

func (s *Item) formatJoin(itemsJoin []models.MapClassifyItemJoin) (items []models.Item) {
	itemsMap := make(map[int]models.Item)
	for _, one := range itemsJoin {
		if target, ok := itemsMap[one.Item.Id]; !ok {
			item := one.Item
			item.Classify = make([]models.Classify, 0)
			classify := one.Classify
			if classify.Id != 0 {
				item.Classify = append(item.Classify, classify)
			}
			itemsMap[one.Item.Id] = item
		} else if one.Item.Id != 0 {
			classify := one.Classify
			target.Classify = append(target.Classify, classify)
			itemsMap[one.Item.Id] = target
		}
	}

	items = make([]models.Item, 0)
	for _, one := range itemsMap {
		items = append(items, one)
	}
	return
}
