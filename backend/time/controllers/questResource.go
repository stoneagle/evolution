package controllers

import (
	"evolution/backend/common/middles"

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
	questResource.POST("", c.Create)
	questResource.POST("/list", c.List)
	questResource.PUT("/:id", c.Update)
	questResource.DELETE("/:id", c.Delete)
}
