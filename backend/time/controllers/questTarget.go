package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type QuestTarget struct {
	structs.Controller
	QuestTargetSvc *services.QuestTarget
}

func NewQuestTarget() *QuestTarget {
	QuestTarget := &QuestTarget{}
	QuestTarget.Init()
	QuestTarget.ProjectName = QuestTarget.Config.Time.System.Name
	QuestTarget.Name = "quest-target"
	QuestTarget.Prepare()
	QuestTarget.QuestTargetSvc = services.NewQuestTarget(QuestTarget.Engine, QuestTarget.Cache)
	return QuestTarget
}

func (c *QuestTarget) Router(router *gin.RouterGroup) {
	questTarget := router.Group(c.Name).Use(middles.One(c.QuestTargetSvc, c.Name))
	questTarget.GET("/get/:id", c.One)
	questTarget.GET("/list", c.List)
	questTarget.POST("", c.Add)
	questTarget.POST("/batch", c.BatchAdd)
	questTarget.POST("/list", c.ListByCondition)
	questTarget.PUT("/:id", c.Update)
	questTarget.DELETE("/:id", c.Delete)
}

func (c *QuestTarget) One(ctx *gin.Context) {
	questTarget := ctx.MustGet(c.Name).(models.QuestTarget)
	resp.Success(ctx, questTarget)
}

func (c *QuestTarget) List(ctx *gin.Context) {
	questTargets, err := c.QuestTargetSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTarget get error", err)
		return
	}
	resp.Success(ctx, questTargets)
}

func (c *QuestTarget) ListByCondition(ctx *gin.Context) {
	var questTarget models.QuestTarget
	if err := ctx.ShouldBindJSON(&questTarget); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	questTargets, err := c.QuestTargetSvc.ListWithCondition(&questTarget)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTarget get error", err)
		return
	}
	resp.Success(ctx, questTargets)
}

func (c *QuestTarget) Add(ctx *gin.Context) {
	var questTarget models.QuestTarget
	if err := ctx.ShouldBindJSON(&questTarget); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTargetSvc.Add(&questTarget)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTarget insert error", err)
		return
	}
	resp.Success(ctx, questTarget)
}

func (c *QuestTarget) BatchAdd(ctx *gin.Context) {
	batchQuestTarget := make([]models.QuestTarget, 0)
	if err := ctx.ShouldBindJSON(&batchQuestTarget); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTargetSvc.BatchAdd(batchQuestTarget)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTarget batch insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestTarget) Update(ctx *gin.Context) {
	var questTarget models.QuestTarget
	if err := ctx.ShouldBindJSON(&questTarget); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTargetSvc.Update(questTarget.Id, questTarget)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTarget update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestTarget) Delete(ctx *gin.Context) {
	questTarget := ctx.MustGet(c.Name).(models.QuestTarget)
	err := c.QuestTargetSvc.Delete(questTarget.Id, questTarget)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTarget delete error", err)
		return
	}
	resp.Success(ctx, questTarget.Id)
}
