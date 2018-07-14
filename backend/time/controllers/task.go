package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	structs.Controller
	TaskSvc      *services.Task
	ProjectSvc   *services.Project
	QuestTeamSvc *services.QuestTeam
}

func NewTask() *Task {
	Task := &Task{}
	Task.Init()
	Task.ProjectName = Task.Config.Time.System.Name
	Task.Name = "task"
	Task.Prepare()
	Task.TaskSvc = services.NewTask(Task.Engine, Task.Cache)
	Task.ProjectSvc = services.NewProject(Task.Engine, Task.Cache)
	Task.QuestTeamSvc = services.NewQuestTeam(Task.Engine, Task.Cache)
	return Task
}

func (c *Task) Router(router *gin.RouterGroup) {
	task := router.Group(c.Name).Use(middles.One(c.TaskSvc, c.Name))
	task.GET("/get/:id", c.One)
	task.GET("/list", c.List)
	task.GET("/list/user/:uid", c.ListByUser)
	task.POST("", c.Add)
	task.POST("/list", c.ListByCondition)
	task.PUT("/:id", c.Update)
	task.DELETE("/:id", c.Delete)
}

func (c *Task) One(ctx *gin.Context) {
	task := ctx.MustGet(c.Name).(models.Task)
	resp.Success(ctx, task)
}

func (c *Task) List(ctx *gin.Context) {
	tasks, err := c.TaskSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task get error", err)
		return
	}
	resp.Success(ctx, tasks)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "team get error", err)
		return
	}

	questIds := make([]int, 0)
	for _, one := range teams {
		questIds = append(questIds, one.QuestId)
	}

	project := models.Project{}
	project.QuestIds = questIds
	projects, err := c.ProjectSvc.ListWithCondition(&project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task get error", err)
		return
	}
	resp.Success(ctx, tasks)
}

func (c *Task) Add(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.TaskSvc.Add(task)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Task) Update(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.TaskSvc.Update(task.Id, task)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Task) Delete(ctx *gin.Context) {
	task := ctx.MustGet(c.Name).(models.Task)
	err := c.TaskSvc.Delete(task.Id, task)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task delete error", err)
		return
	}
	resp.Success(ctx, task.Id)
}
