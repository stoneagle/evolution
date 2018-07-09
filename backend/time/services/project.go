package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Project struct {
	structs.Service
}

func NewProject(engine *xorm.Engine, cache *redis.Client) *Project {
	ret := Project{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Project) One(id int) (interface{}, error) {
	model := models.Project{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Project) Add(model models.Project) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *Project) Update(id int, model models.Project) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *Project) Delete(id int, model models.Project) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Project) List() (projects []models.Project, err error) {
	projects = make([]models.Project, 0)
	err = s.Engine.Find(&projects)
	return
}

func (s *Project) ListWithCondition(project *models.Project) (projects []models.Project, err error) {
	projects = make([]models.Project, 0)
	sql := s.Engine.Asc("level")
	condition := project.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&projects)
	if err != nil {
		return
	}
	return
}
