package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"math"
	"time"

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
	emptyTime := time.Time{}
	if model.EndDate != emptyTime {
		task, err := s.One(id)
		if err != nil {
			return err
		}
		diffHours := model.EndDate.Sub(task.(models.Task).StartDate).Hours()
		model.Duration = int(math.Ceil(diffHours / 24))
	}
	_, err = s.Engine.Id(id).Update(&model)
	if model.StartDateReset {
		err = s.UpdateByMap(id, map[string]interface{}{
			"start_date": nil,
		})
		if err != nil {
			return
		}
	}
	if model.EndDateReset {
		err = s.UpdateByMap(id, map[string]interface{}{
			"end_date": nil,
			"duration": 0,
		})
		if err != nil {
			return
		}
	}
	return
}

func (s *Task) UpdateByMap(id int, model map[string]interface{}) (err error) {
	_, err = s.Engine.Table(new(models.Task)).Id(id).Update(&model)
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
