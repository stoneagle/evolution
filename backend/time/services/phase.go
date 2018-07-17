package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Phase struct {
	ServicePackage
	structs.Service
}

func NewPhase(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Phase {
	ret := Phase{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Phase) Add(model models.Phase) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Phase) Update(id int, model models.Phase) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Phase) List() (phases []models.Phase, err error) {
	phases = make([]models.Phase, 0)
	err = s.Engine.Find(&phases)
	return
}

func (s *Phase) ListWithCondition(phase *models.Phase) (phases []models.Phase, err error) {
	phases = make([]models.Phase, 0)
	sql := s.Engine.Asc("level")
	condition := phase.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&phases)
	if err != nil {
		return
	}
	return
}
