package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type Phase struct {
	BaseController
}

func NewPhase() *Phase {
	Phase := &Phase{}
	Phase.Resource = ResourcePhase
	return Phase
}

func (c *Phase) Router(router *gin.RouterGroup) {
	phase := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	phase.GET("/get/:id", c.One)
	phase.POST("", c.Create)
	phase.POST("/list", c.List)
	phase.PUT("/:id", c.Update)
	phase.DELETE("/:id", c.Delete)
}
