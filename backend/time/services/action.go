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
	Base
	structs.Service
}

func NewAction(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Action {
	ret := Action{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Action) Add(model *models.Action) (err error) {
	_, err = s.Engine.Insert(model)
	return
}

func (s *Action) Update(id int, model models.Action) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Action) List() (actions []models.Action, err error) {
	actions = make([]models.Action, 0)
	session := s.Engine.Where("1 = 1")
	err = session.Find(&actions)
	if err != nil {
		sql, args := session.LastSQL()
		s.LogSql(sql, args, err)
	}
	return
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
