package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type QuestResource struct {
	structs.Controller
	QuestResourceSvc *services.QuestResource
}

func NewQuestResource() *QuestResource {
	QuestResource := &QuestResource{}
	QuestResource.Init()
	QuestResource.ProjectName = QuestResource.Config.Time.System.Name
	QuestResource.Name = "quest-resource"
	QuestResource.Prepare()
	QuestResource.QuestResourceSvc = services.NewQuestResource(QuestResource.Engine, QuestResource.Cache)
	return QuestResource
}

func (c *QuestResource) Router(router *gin.RouterGroup) {
	questResource := router.Group(c.Name).Use(middles.One(c.QuestResourceSvc, c.Name))
	questResource.GET("/get/:id", c.One)
	questResource.GET("/list", c.List)
	questResource.POST("", c.Add)
	questResource.POST("/list", c.ListByCondition)
	questResource.PUT("/:id", c.Update)
	questResource.DELETE("/:id", c.Delete)
}

func (c *QuestResource) One(ctx *gin.Context) {
	questResource := ctx.MustGet(c.Name).(models.QuestResource)
	resp.Success(ctx, questResource)
}

func (c *QuestResource) List(ctx *gin.Context) {
	questResources, err := c.QuestResourceSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questResource get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questResource get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questResource insert error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questResource update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestResource) Delete(ctx *gin.Context) {
	questResource := ctx.MustGet(c.Name).(models.QuestResource)
	err := c.QuestResourceSvc.Delete(questResource.Id, questResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questResource delete error", err)
		return
	}
	resp.Success(ctx, questResource.Id)
}
