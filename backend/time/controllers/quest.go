package controllers

import (
	"evolution/backend/common/middles"

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
	quest.POST("/list", c.ListWithCondition)
	quest.PUT("/:id", c.Update)
	quest.DELETE("/:id", c.Delete)
}
