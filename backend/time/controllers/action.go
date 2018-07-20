package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type Action struct {
	BaseController
}

func NewAction() *Action {
	Action := &Action{}
	Action.Resource = ResourceAction
	return Action
}

func (c *Action) Router(router *gin.RouterGroup) {
	action := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	action.GET("/get/:id", c.One)
	action.POST("", c.Create)
	action.POST("/list", c.List)
	action.PUT("/:id", c.Update)
	action.DELETE("/:id", c.Delete)
}
