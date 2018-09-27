package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type QuestTarget struct {
	BaseController
}

func NewQuestTarget() *QuestTarget {
	QuestTarget := &QuestTarget{}
	QuestTarget.Resource = ResourceQuestTarget
	return QuestTarget
}

func (c *QuestTarget) Router(router *gin.RouterGroup) {
	questTarget := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	questTarget.GET("/get/:id", c.One)
	questTarget.POST("", c.Create)
	questTarget.POST("/batch", c.BatchSave)
	questTarget.POST("/list", c.List)
	questTarget.POST("/count", c.Count)
	questTarget.PUT("/:id", c.Update)
	questTarget.DELETE("/:id", c.Delete)
}

func (c *QuestTarget) BatchSave(ctx *gin.Context) {
	batchQuestTarget := make([]models.QuestTarget, 0)
	if err := ctx.ShouldBindJSON(&batchQuestTarget); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTargetSvc.BatchSave(batchQuestTarget)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questTarget batch insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}
