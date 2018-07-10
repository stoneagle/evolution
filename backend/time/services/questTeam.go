package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestTeam struct {
	structs.Service
}

func NewQuestTeam(engine *xorm.Engine, cache *redis.Client) *QuestTeam {
	ret := QuestTeam{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *QuestTeam) One(id int) (interface{}, error) {
	model := models.QuestTeam{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *QuestTeam) Add(model models.QuestTeam) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *QuestTeam) Update(id int, model models.QuestTeam) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *QuestTeam) Delete(id int, model models.QuestTeam) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *QuestTeam) List() (questTeam []models.QuestTeam, err error) {
	questTeam = make([]models.QuestTeam, 0)
	err = s.Engine.Find(&questTeam)
	return
}

func (s *QuestTeam) ListWithCondition(questTeam *models.QuestTeam) (questTeams []models.QuestTeam, err error) {
	questTeams = make([]models.QuestTeam, 0)
	condition := questTeam.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&questTeams)
	if err != nil {
		return
	}
	return
}
