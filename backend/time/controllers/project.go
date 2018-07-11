package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Project struct {
	structs.Controller
	ProjectSvc *services.Project
}

func NewProject() *Project {
	Project := &Project{}
	Project.Init()
	Project.ProjectName = Project.Config.Time.System.Name
	Project.Name = "project"
	Project.Prepare()
	Project.ProjectSvc = services.NewProject(Project.Engine, Project.Cache)
	return Project
}

func (c *Project) Router(router *gin.RouterGroup) {
	project := router.Group(c.Name).Use(middles.One(c.ProjectSvc, c.Name))
	project.GET("/get/:id", c.One)
	project.GET("/list", c.List)
	project.POST("", c.Add)
	project.POST("/list", c.ListByCondition)
	project.PUT("/:id", c.Update)
	project.DELETE("/:id", c.Delete)
}

func (c *Project) One(ctx *gin.Context) {
	project := ctx.MustGet(c.Name).(models.Project)
	resp.Success(ctx, project)
}

func (c *Project) List(ctx *gin.Context) {
	projects, err := c.ProjectSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project insert error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Project) Delete(ctx *gin.Context) {
	project := ctx.MustGet(c.Name).(models.Project)
	err := c.ProjectSvc.Delete(project.Id, project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project delete error", err)
		return
	}
	resp.Success(ctx, project.Id)
}
