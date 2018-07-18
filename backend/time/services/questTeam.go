package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestTeam struct {
	ServicePackage
	structs.Service
}

func NewQuestTeam(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *QuestTeam {
	ret := QuestTeam{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *QuestTeam) ListWithCondition(questTeam *models.QuestTeam) (questTeams []models.QuestTeam, err error) {
	questTeamsJoin := make([]models.QuestTeamJoin, 0)
	sql := s.Engine.Unscoped().Table("quest_team").Join("INNER", "quest", "quest.id = quest_team.quest_id")

	condition := questTeam.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&questTeamsJoin)
	if err != nil {
		return
	}
	questTeams = make([]models.QuestTeam, 0)
	for _, one := range questTeamsJoin {
		one.QuestTeam.Quest = one.Quest
		questTeams = append(questTeams, one.QuestTeam)
	}
	return
}
