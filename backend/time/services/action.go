package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
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

func (s *Action) Create(modelPtr structs.ModelGeneral) (err error) {
	actionPtr := modelPtr.(*models.Action)
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	_, err = session.Insert(actionPtr)
	if err != nil {
		session.Rollback()
		return
	}

	task := models.NewTask()
	err = s.Pack.TaskSvc.One(actionPtr.TaskId, task)
	if err != nil {
		session.Rollback()
		return
	}

	if task.EndDate.Before(actionPtr.EndDate) {
		task.EndDate = actionPtr.EndDate
		_, err = session.Id(actionPtr.TaskId).Update(task)
		if err != nil {
			session.Rollback()
			return
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return
}

func (s *Action) Update(id int, modelPtr structs.ModelGeneral) (err error) {
	actionPtr := modelPtr.(*models.Action)
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	_, err = session.Id(id).Update(actionPtr)
	if err != nil {
		session.Rollback()
		return
	}

	task := models.NewTask()
	err = s.Pack.TaskSvc.One(actionPtr.TaskId, task)
	if err != nil {
		session.Rollback()
		return
	}

	if task.EndDate.Before(actionPtr.EndDate) {
		task.EndDate = actionPtr.EndDate
		_, err = session.Id(actionPtr.TaskId).Update(task)
		if err != nil {
			session.Rollback()
			return
		}
	} else {
		lastAction := models.NewAction()
		has, err := session.Desc("end_date").Where("task_id = ?", actionPtr.TaskId).And("id != ?", actionPtr.Id).Get(lastAction)
		if err != nil {
			session.Rollback()
			return err
		}

		if !has {
			task.EndDate = actionPtr.EndDate
		} else {
			if lastAction.EndDate.Before(actionPtr.EndDate) {
				task.EndDate = actionPtr.EndDate
			} else {
				task.EndDate = lastAction.EndDate
			}
		}

		_, err = session.Id(actionPtr.TaskId).Update(task)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return
}

func (s *Action) Delete(id int, modelPtr structs.ModelGeneral) (err error) {
	actionPtr := modelPtr.(*models.Action)
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	task := actionPtr.Task
	lastAction := models.NewAction()
	has, err := session.Desc("end_date").Where("task_id = ?", actionPtr.TaskId).And("id != ?", actionPtr.Id).Get(lastAction)
	if err != nil {
		session.Rollback()
		return err
	}
	if has {
		if task.EndDate.After(lastAction.EndDate) || (task.EndDate == lastAction.EndDate) {
			task.EndDate = lastAction.EndDate
			_, err = session.Id(actionPtr.TaskId).Update(task)
			if err != nil {
				session.Rollback()
				return
			}
		}
	} else {
		err = s.Service.UpdateByMap(task.TableName(), task.Id, map[string]interface{}{
			"end_date": nil,
		})
		if err != nil {
			session.Rollback()
			return
		}
	}
	_, err = session.Id(actionPtr.Id).Delete(actionPtr)
	if err != nil {
		session.Rollback()
		return
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return
}
