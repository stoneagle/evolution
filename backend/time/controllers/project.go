package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Project struct {
	Base
}

func NewProject() *Project {
	Project := &Project{}
	Project.Resource = ResourceProject
	return Project
}

func (c *Project) Router(router *gin.RouterGroup) {
	project := router.Group(c.Resource).Use(middles.OnInit(c)).Use(middles.One(c.ProjectSvc, c.Resource, models.Project{}))
	project.GET("/get/:id", c.One)
	project.GET("/list", c.List)
	project.POST("", c.Add)
	project.POST("/list", c.ListByCondition)
	project.PUT("/:id", c.Update)
	project.DELETE("/:id", c.Delete)
}

func (c *Project) One(ctx *gin.Context) {
	project := ctx.MustGet(c.Resource).(*models.Project)
	resp.Success(ctx, project)
}

func (c *Project) List(ctx *gin.Context) {
	projects, err := c.ProjectSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "project get error", err)
		return
	}
	resp.Success(ctx, projects)
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

func (c *Project) Add(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ProjectSvc.Add(project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "project insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Project) Update(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ProjectSvc.Update(project.Id, project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "project update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Project) Delete(ctx *gin.Context) {
	project := ctx.MustGet(c.Resource).(*models.Project)
	err := c.ProjectSvc.Delete(project.Id, project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "project delete error", err)
		return
	}
	resp.Success(ctx, project.Id)
}
