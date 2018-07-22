package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Task struct {
	Pack ServicePackage
	structs.Service
}

func NewTask(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Task {
	ret := Task{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Task) Update(id int, modelPtr structs.ModelGeneral) (err error) {
	taskPtr := modelPtr.(*models.Task)
	err = s.Service.Update(id, taskPtr)
	if err != nil {
		return
	}

	if taskPtr.StartDateReset {
		err = s.Service.UpdateByMap(taskPtr.TableName(), id, map[string]interface{}{
			"start_date": nil,
		})
		if err != nil {
			return
		}
	}

	if taskPtr.EndDateReset {
		err = s.Service.UpdateByMap(taskPtr.TableName(), id, map[string]interface{}{
			"end_date": nil,
		})
		if err != nil {
			return
		}
	}
	return
}
