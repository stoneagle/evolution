package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

type Action struct {
	structs.Service
}

func NewAction(engine *xorm.Engine, cache *redis.Client) *Action {
	ret := Action{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Action) One(id int) (interface{}, error) {
	model := models.Action{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Action) Add(model *models.Action) (err error) {
	_, err = s.Engine.Insert(model)
	return
}

func (s *Action) Update(id int, model models.Action) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Action) Delete(id int, model models.Action) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Action) List() (actions []models.Action, err error) {
	actions = make([]models.Action, 0)
	err = s.Engine.Find(&actions)
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
