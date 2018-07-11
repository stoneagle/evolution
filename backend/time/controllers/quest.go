package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Quest struct {
	structs.Controller
	QuestSvc *services.Quest
}

func NewQuest() *Quest {
	Quest := &Quest{}
	Quest.Init()
	Quest.ProjectName = Quest.Config.Time.System.Name
	Quest.Name = "quest"
	Quest.Prepare()
	Quest.QuestSvc = services.NewQuest(Quest.Engine, Quest.Cache)
	return Quest
}

func (c *Quest) Router(router *gin.RouterGroup) {
	quest := router.Group(c.Name).Use(middles.One(c.QuestSvc, c.Name))
	quest.GET("/get/:id", c.One)
	quest.GET("/list", c.List)
	quest.POST("", c.Add)
	quest.POST("/list", c.ListByCondition)
	quest.PUT("/:id", c.Update)
	quest.DELETE("/:id", c.Delete)
}

func (c *Quest) One(ctx *gin.Context) {
	quest := ctx.MustGet(c.Name).(models.Quest)
	resp.Success(ctx, quest)
}

func (c *Quest) List(ctx *gin.Context) {
	quests, err := c.QuestSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest insert error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Quest) Delete(ctx *gin.Context) {
	quest := ctx.MustGet(c.Name).(models.Quest)
	err := c.QuestSvc.Delete(quest.Id, quest)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest delete error", err)
		return
	}
	resp.Success(ctx, quest.Id)
}
