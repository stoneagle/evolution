package services

import (
	"quant/backend/models"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Classify struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewClassify(engine *xorm.Engine, cache *redis.Client) *Classify {
	return &Classify{
		engine: engine,
		cache:  cache,
	}
}

func (s *Classify) One(id int) (classify models.Classify, err error) {
	classify = models.Classify{}
	_, err = s.engine.Where("id = ?", id).Get(&classify)
	return
}

func (s *Classify) List() (classifies []models.Classify, err error) {
	classifies = make([]models.Classify, 0)
	err = s.engine.Find(&classifies)
	return
}

func (s *Classify) BatchSave(classifyBatch []models.Classify) (err error) {
	session := s.engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}
	for _, one := range classifyBatch {
		has, err := session.Get(&one)
		if err != nil {
			session.Rollback()
			return err
		}
		if has {
			updateOne := &models.Classify{}
			updateOne.GeneralWithDeleted.UpdatedAt = time.Now()
			_, err = session.Id(one.Id).Update(updateOne)
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
	}
	// TODO 删除不再使用的分类
	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return nil
}

func (s *Classify) UpdateByMap(id int, classify map[string]interface{}) (err error) {
	_, err = s.engine.Table(new(models.Classify)).Id(id).Update(&classify)
	return
}

func (s *Classify) Delete(id int) (err error) {
	var classify models.Classify
	_, err = s.engine.Id(id).Get(&classify)
	if err == nil {
		_, err = s.engine.Id(id).Delete(&classify)
	}
	return
}
