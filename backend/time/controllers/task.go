package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

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
	task.GET("/list/user/:uid", c.ListByUser)
	task.POST("", c.Create)
	task.POST("/list", c.List)
	task.PUT("/:id", c.Update)
	task.DELETE("/:id", c.Delete)
}

func (c *Task) ListByUser(ctx *gin.Context) {
	userIdStr := ctx.Param("uid")
	if userIdStr != "" {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "user id params error", err)
			return
		}
		c.QuestTeamModel.UserId = userId
	}

	questTeamsGeneralPtr := c.QuestTeamModel.SlicePtr()
	err := c.QuestTeamSvc.List(c.QuestTeamModel, questTeamsGeneralPtr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "team get error", err)
		return
	}
	questTeamsPtr := c.QuestTeamModel.Transfer(questTeamsGeneralPtr)
	questIds := make([]int, 0)
	for _, one := range *questTeamsPtr {
		questIds = append(questIds, one.QuestId)
	}

	c.TaskModel.QuestTarget.QuestIds = questIds
	tasksPtr := c.TaskModel.SlicePtr()
	err = c.TaskSvc.List(c.TaskModel, tasksPtr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "task get error", err)
		return
	}
	resp.Success(ctx, tasksPtr)
}
