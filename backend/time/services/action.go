package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

type Action struct {
	Pack ServicePackage
	structs.Service
}

func NewAction(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Action {
	ret := Action{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Action) ListWithCondition(action *models.Action) (actions []models.Action, err error) {
	actionsJoin := make([]models.ActionJoin, 0)
	sql := s.Engine.Unscoped().Table("action").Join("INNER", "task", "task.id = action.task_id")

	condition := action.BuildCondition()
	sql = sql.Where(condition)
	emptyTime := time.Time{}
	if action.StartDate != emptyTime {
		sql = sql.And(builder.Gte{"action.start_date": action.StartDate})
	}
	if action.EndDate != emptyTime {
		sql = sql.And(builder.Lte{"action.start_date": action.EndDate})
	}

	err = sql.Find(&actionsJoin)
	if err != nil {
		return
	}

	actions = make([]models.Action, 0)
	for _, one := range actionsJoin {
		one.Action.Task = one.Task
		actions = append(actions, one.Action)
	}
	return
}
