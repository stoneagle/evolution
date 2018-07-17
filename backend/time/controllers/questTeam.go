package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type QuestTeam struct {
	BaseController
}

func NewQuestTeam() *QuestTeam {
	QuestTeam := &QuestTeam{}
	QuestTeam.Resource = ResourceQuestTeam
	return QuestTeam
}

func (c *QuestTeam) Router(router *gin.RouterGroup) {
	questTeam := router.Group(c.Resource).Use(middles.OnInit(c))
	questTeam.GET("/get/:id", c.One)
	questTeam.GET("/list", c.List)
	questTeam.POST("", c.Add)
	questTeam.POST("/list", c.ListByCondition)
	questTeam.PUT("/:id", c.Update)
	questTeam.DELETE("/:id", c.Delete)
}

func (c *QuestTeam) One(ctx *gin.Context) {
	questTeam := ctx.MustGet(c.Resource).(*models.QuestTeam)
	resp.Success(ctx, questTeam)
}

func (c *QuestTeam) List(ctx *gin.Context) {
	questTeams, err := c.QuestTeamSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questTeam get error", err)
		return
	}
	resp.Success(ctx, questTeams)
}

func (c *QuestTeam) ListByCondition(ctx *gin.Context) {
	var questTeam models.QuestTeam
	if err := ctx.ShouldBindJSON(&questTeam); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	questTeams, err := c.QuestTeamSvc.ListWithCondition(&questTeam)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questTeam get error", err)
		return
	}
	resp.Success(ctx, questTeams)
}

func (c *QuestTeam) Add(ctx *gin.Context) {
	var questTeam models.QuestTeam
	if err := ctx.ShouldBindJSON(&questTeam); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTeamSvc.Add(questTeam)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questTeam insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestTeam) Update(ctx *gin.Context) {
	var questTeam models.QuestTeam
	if err := ctx.ShouldBindJSON(&questTeam); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTeamSvc.Update(questTeam.Id, questTeam)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questTeam update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestTeam) Delete(ctx *gin.Context) {
	questTeam := ctx.MustGet(c.Resource).(*models.QuestTeam)
	err := c.QuestTeamSvc.Delete(questTeam.Id, questTeam)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questTeam delete error", err)
		return
	}
	resp.Success(ctx, questTeam.Id)
}
