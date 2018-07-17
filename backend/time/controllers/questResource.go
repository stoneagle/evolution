package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type QuestResource struct {
	BaseController
}

func NewQuestResource() *QuestResource {
	QuestResource := &QuestResource{}
	QuestResource.Resource = ResourceQuestResource
	return QuestResource
}

func (c *QuestResource) Router(router *gin.RouterGroup) {
	questResource := router.Group(c.Resource).Use(middles.OnInit(c))
	questResource.GET("/get/:id", c.One)
	questResource.GET("/list", c.List)
	questResource.POST("", c.Add)
	questResource.POST("/list", c.ListByCondition)
	questResource.PUT("/:id", c.Update)
	questResource.DELETE("/:id", c.Delete)
}

func (c *QuestResource) One(ctx *gin.Context) {
	questResource := ctx.MustGet(c.Resource).(*models.QuestResource)
	resp.Success(ctx, questResource)
}

func (c *QuestResource) List(ctx *gin.Context) {
	questResources, err := c.QuestResourceSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questResource get error", err)
		return
	}
	resp.Success(ctx, questResources)
}

func (c *QuestResource) ListByCondition(ctx *gin.Context) {
	var questResource models.QuestResource
	if err := ctx.ShouldBindJSON(&questResource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	questResources, err := c.QuestResourceSvc.ListWithCondition(&questResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questResource get error", err)
		return
	}
	resp.Success(ctx, questResources)
}

func (c *QuestResource) Add(ctx *gin.Context) {
	var questResource models.QuestResource
	if err := ctx.ShouldBindJSON(&questResource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestResourceSvc.Add(questResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questResource insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestResource) Update(ctx *gin.Context) {
	var questResource models.QuestResource
	if err := ctx.ShouldBindJSON(&questResource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestResourceSvc.Update(questResource.Id, questResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questResource update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestResource) Delete(ctx *gin.Context) {
	questResource := ctx.MustGet(c.Resource).(*models.QuestResource)
	err := c.QuestResourceSvc.Delete(questResource.Id, questResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questResource delete error", err)
		return
	}
	resp.Success(ctx, questResource.Id)
}
