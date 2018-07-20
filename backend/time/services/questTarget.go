package services

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"

	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/common/utils"
	"evolution/backend/time/models"
)

type QuestTarget struct {
	Pack ServicePackage
	structs.Service
}

func NewQuestTarget(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *QuestTarget {
	ret := QuestTarget{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *QuestTarget) BatchSave(targets []models.QuestTarget) (err error) {
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	batchQuestTargetsMap := map[int][]models.QuestTarget{}
	for _, questTarget := range targets {
		questTargets, ok := batchQuestTargetsMap[questTarget.QuestId]
		if !ok {
			batchQuestTargetsMap[questTarget.QuestId] = make([]models.QuestTarget, 0)
		}
		questTargets = append(questTargets, questTarget)
		batchQuestTargetsMap[questTarget.QuestId] = questTargets
	}

	for questId, questTargets := range batchQuestTargetsMap {
		oldTargets := make([]models.QuestTarget, 0)
		err = session.Where("quest_id = ?", questId).Find(&oldTargets)
		if err != nil {
			session.Rollback()
			return
		}
		oldAreaIdsSlice := make([]interface{}, 0)
		oldQuestTargetsMap := map[int]models.QuestTarget{}
		newAreaIdsSlice := make([]interface{}, 0)
		newQuestTargetsMap := map[int]models.QuestTarget{}
		for _, old := range oldTargets {
			oldAreaIdsSlice = append(oldAreaIdsSlice, old.AreaId)
			oldQuestTargetsMap[old.AreaId] = old
		}

		for _, questTarget := range questTargets {
			newAreaIdsSlice = append(newAreaIdsSlice, questTarget.AreaId)
			newQuestTargetsMap[questTarget.AreaId] = questTarget
		}
		needAddSlice := utils.SliceDiff(newAreaIdsSlice, oldAreaIdsSlice)
		needDeleteSlice := utils.SliceDiff(oldAreaIdsSlice, newAreaIdsSlice)

		for _, areaId := range needAddSlice {
			questTarget, ok := newQuestTargetsMap[areaId.(int)]
			if !ok {
				err = errors.New("quest target not exist")
				session.Rollback()
				return
			}
			_, err = session.Insert(&questTarget)
			if err != nil {
				session.Rollback()
				return
			}
		}

		for _, areaId := range needDeleteSlice {
			// TODO relate project check and deny delete
			questTarget := models.QuestTarget{}
			questTarget.AreaId = areaId.(int)
			questTarget.QuestId = questId
			_, err = session.Delete(&questTarget)
			if err != nil {
				session.Rollback()
				return
			}
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return
}

func (s *QuestTarget) ListWithCondition(questTarget *models.QuestTarget) (questTargets []models.QuestTarget, err error) {
	questTargetsJoin := make([]models.QuestTargetJoin, 0)
	sql := s.Engine.Unscoped().Table("quest_target").Join("INNER", "area", "area.id = quest_target.area_id")

	condition := questTarget.BuildCondition()
	sql = sql.Where(condition)
	err = sql.Find(&questTargetsJoin)
	if err != nil {
		return
	}

	questTargets = make([]models.QuestTarget, 0)
	for _, one := range questTargetsJoin {
		one.QuestTarget.Area = one.Area
		questTargets = append(questTargets, one.QuestTarget)
	}
	return
}
