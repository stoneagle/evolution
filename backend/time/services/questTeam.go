package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type QuestTeam struct {
	Pack ServicePackage
	structs.Service
}

func NewQuestTeam(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *QuestTeam {
	ret := QuestTeam{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *QuestTeam) GetQuestsByUser(userId int, questStatus int) (questsPtr *[]models.Quest, err error) {
	questTeam := models.NewQuestTeam()
	questTeam.UserId = userId
	questTeam.Quest.Status = models.QuestStatusExec
	questTeamsGeneralPtr := questTeam.SlicePtr()
	err = s.List(questTeam, questTeamsGeneralPtr)
	if err != nil {
		return
	}
	questTeamsPtr := questTeam.Transfer(questTeamsGeneralPtr)
	quests := make([]models.Quest, 0)
	for _, one := range *questTeamsPtr {
		quests = append(quests, *one.Quest)
	}
	questsPtr = &quests
	return
}
