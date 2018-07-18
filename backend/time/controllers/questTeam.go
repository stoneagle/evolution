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
	questTeam := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	questTeam.GET("/get/:id", c.One)
	questTeam.GET("/list", c.List)
	questTeam.POST("", c.Create)
	questTeam.POST("/list", c.ListByCondition)
	questTeam.PUT("/:id", c.Update)
	questTeam.DELETE("/:id", c.Delete)
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
