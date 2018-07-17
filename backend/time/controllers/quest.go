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
	quest := router.Group(c.Resource).Use(middles.OnInit(c))
	quest.GET("/get/:id", c.One)
	quest.GET("/list", c.List)
	quest.POST("", c.Add)
	quest.POST("/list", c.ListByCondition)
	quest.PUT("/:id", c.Update)
	quest.DELETE("/:id", c.Delete)
}

func (c *Quest) One(ctx *gin.Context) {
	quest := ctx.MustGet(c.Resource).(*models.Quest)
	resp.Success(ctx, quest)
}

func (c *Quest) List(ctx *gin.Context) {
	quests, err := c.QuestSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "quest get error", err)
		return
	}
	resp.Success(ctx, quests)
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

func (c *Quest) Add(ctx *gin.Context) {
	var quest models.Quest
	if err := ctx.ShouldBindJSON(&quest); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestSvc.Add(&quest)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "quest insert error", err)
		return
	}
	resp.Success(ctx, quest)
}

func (c *Quest) Update(ctx *gin.Context) {
	var quest models.Quest
	if err := ctx.ShouldBindJSON(&quest); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestSvc.Update(quest.Id, quest)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "quest update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Quest) Delete(ctx *gin.Context) {
	quest := ctx.MustGet(c.Resource).(*models.Quest)
	err := c.QuestSvc.Delete(quest.Id, quest)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "quest delete error", err)
		return
	}
	resp.Success(ctx, quest.Id)
}
