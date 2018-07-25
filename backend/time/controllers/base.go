package controllers

import (
	"evolution/backend/common/database"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"
	"fmt"
)

const (
	ResourceAction         structs.ResourceType = "action"
	ResourceArea           structs.ResourceType = "area"
	ResourceCountry        structs.ResourceType = "country"
	ResourceField          structs.ResourceType = "field"
	ResourcePhase          structs.ResourceType = "phase"
	ResourceProject        structs.ResourceType = "project"
	ResourceQuest          structs.ResourceType = "quest"
	ResourceQuestResource  structs.ResourceType = "quest-resource"
	ResourceQuestTarget    structs.ResourceType = "quest-target"
	ResourceQuestTeam      structs.ResourceType = "quest-team"
	ResourceQuestTimeTable structs.ResourceType = "quest-timetable"
	ResourceResource       structs.ResourceType = "resource"
	ResourceSyncfusion     structs.ResourceType = "syncfusion"
	ResourceTask           structs.ResourceType = "task"
	ResourceUserResource   structs.ResourceType = "user-resource"
)

type BaseController struct {
	structs.Controller
	services.ServicePackage
	models.ModelPackage
}

func (c *BaseController) Init() {
	c.Prepare(structs.ProjectTime)
	cache := database.GetRedis()
	engine := database.GetXorm(c.Project)
	c.PrepareService(engine, cache, c.Logger)
	c.PrepareModel()
	if c.Resource != ResourceSyncfusion {
		c.ChangeSvc(c.Resource)
		c.ChangeModel(c.Resource)
	}
}

func (c *BaseController) ChangeSvc(resource structs.ResourceType) {
	svcMap := map[structs.ResourceType]structs.ServiceGeneral{
		ResourceAction:         c.ActionSvc,
		ResourceArea:           c.AreaSvc,
		ResourceCountry:        c.CountrySvc,
		ResourceField:          c.FieldSvc,
		ResourcePhase:          c.PhaseSvc,
		ResourceProject:        c.ProjectSvc,
		ResourceQuest:          c.QuestSvc,
		ResourceQuestResource:  c.QuestResourceSvc,
		ResourceQuestTarget:    c.QuestTargetSvc,
		ResourceQuestTeam:      c.QuestTeamSvc,
		ResourceQuestTimeTable: c.QuestTimeTableSvc,
		ResourceResource:       c.ResourceSvc,
		ResourceSyncfusion:     c.SyncfusionSvc,
		ResourceTask:           c.TaskSvc,
		ResourceUserResource:   c.UserResourceSvc,
	}
	svc, ok := svcMap[resource]
	if !ok {
		panic(fmt.Sprintf("%v svc not exist", c.Resource))
	}
	c.Service = svc
}

func (c *BaseController) ChangeModel(resource structs.ResourceType) {
	modelMap := map[structs.ResourceType]structs.ModelGeneral{
		ResourceAction:         c.ActionModel,
		ResourceArea:           c.AreaModel,
		ResourceCountry:        c.CountryModel,
		ResourceField:          c.FieldModel,
		ResourcePhase:          c.PhaseModel,
		ResourceProject:        c.ProjectModel,
		ResourceQuest:          c.QuestModel,
		ResourceQuestResource:  c.QuestResourceModel,
		ResourceQuestTarget:    c.QuestTargetModel,
		ResourceQuestTeam:      c.QuestTeamModel,
		ResourceQuestTimeTable: c.QuestTimeTableModel,
		ResourceResource:       c.ResourceModel,
		ResourceTask:           c.TaskModel,
		ResourceUserResource:   c.UserResourceModel,
	}
	model, ok := modelMap[resource]
	if !ok {
		panic(fmt.Sprintf("%v model not exist", c.Resource))
	}
	c.Model = model
}
