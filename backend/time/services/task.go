package services

import (
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Task struct {
	Basic
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

func (s *Task) UpdateByMap(id int, task map[string]interface{}) (err error) {
	_, err = s.Engine.Table(new(models.Task)).Id(id).Update(&task)
	return
}

func (s *Task) List(tasks []models.Task) (err error) {
	err = s.Engine.Find(&tasks)
	return
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

// func (s *Task) formatJoin(itemsJoin []models.MapTaskItemJoin) (tasks []models.Task) {
// 	tasksMap := make(map[int]models.Task)
// 	for _, one := range itemsJoin {
// 		if target, ok := tasksMap[one.Task.Id]; !ok {
// 			task := one.Task
// 			task.Item = make([]models.Item, 0)
// 			item := one.Item
// 			if item.Id != 0 {
// 				task.Item = append(task.Item, item)
// 			}
// 			tasksMap[one.Task.Id] = task
// 		} else if one.Task.Id != 0 {
// 			item := one.Item
// 			target.Item = append(target.Item, item)
// 			tasksMap[one.Task.Id] = target
// 		}
// 	}

// 	tasks = make([]models.Task, 0)
// 	for _, one := range tasksMap {
// 		tasks = append(tasks, one)
// 	}
// 	return
// }
