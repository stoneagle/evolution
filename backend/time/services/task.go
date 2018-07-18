package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"math"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Task struct {
	ServicePackage
	structs.Service
}

func NewTask(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Task {
	ret := Task{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Task) Update(id int, modelPtr interface{}) (err error) {
	taskPtr := modelPtr.(*models.Task)
	emptyTime := time.Time{}
	if taskPtr.EndDate != emptyTime {
		task := models.Task{}
		err := s.One(id, &task)
		if err != nil {
			return err
		}
		diffHours := taskPtr.EndDate.Sub(task.StartDate).Hours()
		taskPtr.Duration = int(math.Ceil(diffHours / 24))
	}
	_, err = s.Engine.Id(id).Update(taskPtr)
	if taskPtr.StartDateReset {
		err = s.UpdateByMap(id, map[string]interface{}{
			"start_date": nil,
		})
		if err != nil {
			return
		}
	}
	if taskPtr.EndDateReset {
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

func (s *Task) ListWithCondition(task *models.Task) (tasks []models.Task, err error) {
	tasksJoin := make([]models.TaskJoin, 0)
	sql := s.Engine.Unscoped().Table("task").Join("INNER", "resource", "resource.id = task.resource_id").Join("INNER", "map_area_resource", "map_area_resource.resource_id = resource.id")

	condition := task.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&tasksJoin)
	if err != nil {
		return
	}

	tasks = make([]models.Task, 0)
	for _, one := range tasksJoin {
		one.Task.Resource = one.Resource
		one.Task.Resource.Area.Id = one.MapAreaResource.AreaId
		tasks = append(tasks, one.Task)
	}
	return
}
