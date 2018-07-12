package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Task struct {
	structs.Service
}

func NewTask(engine *xorm.Engine, cache *redis.Client) *Task {
	ret := Task{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Task) One(id int) (interface{}, error) {
	model := models.Task{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Task) Add(model models.Task) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Task) Update(id int, model models.Task) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Task) Delete(id int, model models.Task) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Task) List() (tasks []models.Task, err error) {
	tasks = make([]models.Task, 0)
	err = s.Engine.Find(&tasks)
	return
}

func (s *Task) ListWithCondition(task *models.Task) (tasks []models.Task, err error) {
	tasksJoin := make([]models.TaskJoin, 0)
	sql := s.Engine.Unscoped().Table("task").Join("INNER", "resource", "resource.id = task.resource_id")

	condition := task.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&tasksJoin)
	if err != nil {
		return
	}

	tasks = make([]models.Task, 0)
	for _, one := range tasksJoin {
		one.Task.Resource = one.Resource
		tasks = append(tasks, one.Task)
	}
	return
}
