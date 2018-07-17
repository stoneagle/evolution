package services

import (
	"evolution/backend/common/logger"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type ServicePackage struct {
	ActionSvc         *Action
	AreaSvc           *Area
	CountrySvc        *Country
	FieldSvc          *Field
	PhaseSvc          *Phase
	ProjectSvc        *Project
	QuestSvc          *Quest
	QuestResourceSvc  *QuestResource
	QuestTargetSvc    *QuestTarget
	QuestTeamSvc      *QuestTeam
	QuestTimeTableSvc *QuestTimeTable
	ResourceSvc       *Resource
	SyncfusionSvc     *Syncfusion
	TaskSvc           *Task
	UserResourceSvc   *UserResource
}

func (s *ServicePackage) PrepareService(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) {
	s.ActionSvc = NewAction(engine, cache, log)
	s.AreaSvc = NewArea(engine, cache, log)
	s.CountrySvc = NewCountry(engine, cache, log)
	s.FieldSvc = NewField(engine, cache, log)
	s.PhaseSvc = NewPhase(engine, cache, log)
	s.ProjectSvc = NewProject(engine, cache, log)
	s.QuestSvc = NewQuest(engine, cache, log)
	s.QuestResourceSvc = NewQuestResource(engine, cache, log)
	s.QuestTargetSvc = NewQuestTarget(engine, cache, log)
	s.QuestTeamSvc = NewQuestTeam(engine, cache, log)
	s.QuestTimeTableSvc = NewQuestTimeTable(engine, cache, log)
	s.ResourceSvc = NewResource(engine, cache, log)
	s.SyncfusionSvc = NewSyncfusion(engine, cache, log)
	s.TaskSvc = NewTask(engine, cache, log)
	s.UserResourceSvc = NewUserResource(engine, cache, log)
}
