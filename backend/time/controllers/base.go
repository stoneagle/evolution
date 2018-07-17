package controllers

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
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

type Base struct {
	structs.Controller
	services.Base
}

func (c *Base) Init() {
	cache := database.GetRedis()
	engine := database.GetXorm(c.Project)
	c.Prepare(config.ProjectTime)
	c.PrepareService(engine, cache, c.Logger)
}

func (c *Base) One(ctx *gin.Context) {
	one := ctx.MustGet(c.Resource)
	resp.Success(ctx, one)
}
