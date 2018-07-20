package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type Project struct {
	BaseController
}

func NewProject() *Project {
	Project := &Project{}
	Project.Resource = ResourceProject
	return Project
}

func (c *Project) Router(router *gin.RouterGroup) {
	project := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	project.GET("/get/:id", c.One)
	project.GET("/list", c.List)
	project.POST("", c.Create)
	project.POST("/list", c.ListWithCondition)
	project.PUT("/:id", c.Update)
	project.DELETE("/:id", c.Delete)
}
