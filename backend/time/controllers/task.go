package controllers

import (
	"evolution/backend/common/resp"
	"evolution/backend/time/middles"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Task struct {
	Base
	Name    string
	TaskSvc *services.Task
}

func NewTask() *Task {
	Task := &Task{
		Name: "task",
	}
	Task.Prepare()
	Task.TaskSvc = services.NewTask(Task.Engine, Task.Cache)
	return Task
}

func (c *Task) Router(router *gin.RouterGroup) {
	task := router.Group("task").Use(middles.One(c.TaskSvc, c.Name))
	task.GET("/get/:id", c.One)
	task.GET("/list", c.List)
	task.POST("", c.Add)
	task.PUT("/:id", c.Update)
	task.DELETE("/:id", c.Delete)
}

func (c *Task) One(ctx *gin.Context) {
	task := ctx.MustGet("task").(models.Task)
	resp.Success(ctx, task)
}

func (c *Task) List(ctx *gin.Context) {
	tasks := make([]models.Task, 0)
	err := c.TaskSvc.List(tasks)
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
	task := ctx.MustGet("task").(models.Task)
	err := c.TaskSvc.Delete(task.Id, task)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task delete error", err)
		return
	}
	resp.Success(ctx, task.Id)
}
