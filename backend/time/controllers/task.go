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
	task.GET("/list/user/:uid", c.ListByUser)
	task.POST("", c.Create)
	task.POST("/list", c.List)
	task.POST("/count", c.Count)
	task.PUT("/:id", c.Update)
	task.DELETE("/:id", c.Delete)
}

func (c *Task) ListByUser(ctx *gin.Context) {
	userIdStr := ctx.Param("uid")
	var userId int
	var err error
	if userIdStr != "" {
		userId, err = strconv.Atoi(userIdStr)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "user id params error", err)
			return
		}
	}

	questsPtr, err := c.QuestTeamSvc.GetQuestsByUser(userId, models.QuestStatusExec)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "user relate quests get error", err)
		return
	}
	questIds := make([]int, 0)
	for _, one := range *questsPtr {
		questIds = append(questIds, one.Id)
	}

	c.TaskModel.QuestTarget.QuestIds = questIds
	c.TaskModel.Project.Status = models.ProjectStatusWait
	tasksPtr := c.TaskModel.SlicePtr()
	err = c.TaskSvc.List(c.TaskModel, tasksPtr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "task get error", err)
		return
	}
	resp.Success(ctx, tasksPtr)
}
