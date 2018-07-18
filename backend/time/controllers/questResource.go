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
	questResource := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	questResource.GET("/get/:id", c.One)
	questResource.GET("/list", c.List)
	questResource.POST("", c.Create)
	questResource.POST("/list", c.ListByCondition)
	questResource.PUT("/:id", c.Update)
	questResource.DELETE("/:id", c.Delete)
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
