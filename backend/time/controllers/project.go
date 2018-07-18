package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

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
	project.POST("/list", c.ListByCondition)
	project.PUT("/:id", c.Update)
	project.DELETE("/:id", c.Delete)
}

func (c *Project) ListByCondition(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	projects, err := c.ProjectSvc.ListWithCondition(&project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "project get error", err)
		return
	}
	resp.Success(ctx, projects)
}
