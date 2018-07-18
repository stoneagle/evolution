package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	BaseController
}

func NewTask() *Task {
	Task := &Task{}
	Task.Resource = ResourceTask
	return Task
}

func (c *Task) Router(router *gin.RouterGroup) {
	task := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	task.GET("/get/:id", c.One)
	task.GET("/list", c.List)
	task.GET("/list/user/:uid", c.ListByUser)
	task.POST("", c.Create)
	task.POST("/list", c.ListByCondition)
	task.PUT("/:id", c.Update)
	task.DELETE("/:id", c.Delete)
}

func (c *Task) ListByUser(ctx *gin.Context) {
	userIdStr := ctx.Param("uid")
	questTeam := models.QuestTeam{}
	if userIdStr != "" {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "user id params error", err)
			return
		}
		questTeam.UserId = userId
	}

	teams, err := c.QuestTeamSvc.ListWithCondition(&questTeam)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "team get error", err)
		return
	}

	questIds := make([]int, 0)
	for _, one := range teams {
		questIds = append(questIds, one.QuestId)
	}

	project := models.Project{}
	project.QuestTarget.QuestIds = questIds
	projects, err := c.ProjectSvc.ListWithCondition(&project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "project get error", err)
		return
	}

	projectIds := make([]int, 0)
	for _, one := range projects {
		projectIds = append(projectIds, one.Id)
	}

	task := models.Task{}
	task.ProjectIds = projectIds
	tasks, err := c.TaskSvc.ListWithCondition(&task)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "task get error", err)
		return
	}
	resp.Success(ctx, tasks)
}

func (c *Task) ListByCondition(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	tasks, err := c.TaskSvc.ListWithCondition(&task)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "task get error", err)
		return
	}
	resp.Success(ctx, tasks)
}
