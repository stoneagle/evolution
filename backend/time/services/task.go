package services

import (
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Task struct {
	engine *xorm.Engine
	cache  *redis.Client
}

func NewTask(engine *xorm.Engine, cache *redis.Client) *Task {
	return &Task{
		engine: engine,
		cache:  cache,
	}
}

func (s *Task) One(id int) (task models.Task, err error) {
	task = models.Task{}
	_, err = s.engine.Where("id = ?", id).Get(&task)
	return
}

func (s *Task) List() (tasks []models.Task, err error) {
	tasks = make([]models.Task, 0)
	err = s.engine.Find(&tasks)
	return
}

func (s *Task) UpdateByMap(id int, task map[string]interface{}) (err error) {
	_, err = s.engine.Table(new(models.Task)).Id(id).Update(&task)
	return
}

func (s *Task) Update(id int, task *models.Task) (err error) {
	_, err = s.engine.Id(id).Update(task)
	return
}

func (s *Task) Add(task *models.Task) (err error) {
	_, err = s.engine.Insert(task)
	return
}

func (s *Task) Delete(id int) (err error) {
	var task models.Task
	_, err = s.engine.Id(id).Get(&task)
	if err == nil {
		_, err = s.engine.Id(id).Delete(&task)
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
