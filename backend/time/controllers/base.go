package controllers

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

const (
	ResourceAction         string = "action"
	ResourceArea           string = "area"
	ResourceCountry        string = "country"
	ResourceField          string = "field"
	ResourcePhase          string = "phase"
	ResourceProject        string = "project"
	ResourceQuest          string = "quest"
	ResourceQuestResource  string = "quest-resource"
	ResourceQuestTarget    string = "quest-target"
	ResourceQuestTeam      string = "quest-team"
	ResourceQuestTimeTable string = "quest-timetable"
	ResourceResource       string = "resource"
	ResourceSyncfusion     string = "syncfusion"
	ResourceTask           string = "task"
	ResourceUserResource   string = "user-resource"
)

type BaseController struct {
	structs.Controller
	services.ServicePackage
	models.ModelPackage
}

func (c *BaseController) Init() {
	c.Prepare(config.ProjectTime)
	cache := database.GetRedis()
	engine := database.GetXorm(c.Project)
	c.PrepareService(engine, cache, c.Logger)
	c.PrepareModel()
	c.ResourceSvcMap = map[string]structs.ServiceGeneral{
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
	c.ResourceModelMap = map[string]structs.ModelGeneral{
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
}

func (c *BaseController) One(ctx *gin.Context) {
	one := ctx.MustGet(c.Resource)
	resp.Success(ctx, one)
}
