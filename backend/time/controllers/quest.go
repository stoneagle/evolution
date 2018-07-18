package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Quest struct {
	BaseController
}

func NewQuest() *Quest {
	Quest := &Quest{}
	Quest.Resource = ResourceQuest
	return Quest
}

func (c *Quest) Router(router *gin.RouterGroup) {
	quest := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	quest.GET("/get/:id", c.One)
	quest.GET("/list", c.List)
	quest.POST("", c.Create)
	quest.POST("/list", c.ListByCondition)
	quest.PUT("/:id", c.Update)
	quest.DELETE("/:id", c.Delete)
}

func (c *Quest) ListByCondition(ctx *gin.Context) {
	var quest models.Quest
	if err := ctx.ShouldBindJSON(&quest); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	quests, err := c.QuestSvc.ListWithCondition(&quest)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "quest get error", err)
		return
	}
	resp.Success(ctx, quests)
}
