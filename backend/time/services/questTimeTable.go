package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestTimeTable struct {
	Base
	structs.Service
}

func NewQuestTimeTable(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *QuestTimeTable {
	ret := QuestTimeTable{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *QuestTimeTable) Add(model models.QuestTimeTable) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *QuestTimeTable) Update(id int, model models.QuestTimeTable) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *QuestTimeTable) List() (questTimeTable []models.QuestTimeTable, err error) {
	questTimeTable = make([]models.QuestTimeTable, 0)
	err = s.Engine.Find(&questTimeTable)
	return
}

func (s *QuestTimeTable) ListWithCondition(questTimeTable *models.QuestTimeTable) (questTimeTables []models.QuestTimeTable, err error) {
	questTimeTables = make([]models.QuestTimeTable, 0)
	condition := questTimeTable.BuildCondition()
	sql := s.Engine.Where(condition)
	err = sql.Find(&questTimeTables)
	if err != nil {
		return
	}
	return
}
