package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

type Action struct {
	structs.Controller
	ActionSvc *services.Action
}

func NewAction() *Action {
	Action := &Action{}
	Action.Init()
	Action.ProjectName = Action.Config.Time.System.Name
	Action.Name = "action"
	Action.Prepare()
	Action.ActionSvc = services.NewAction(Action.Engine, Action.Cache)
	return Action
}

func (c *Action) Router(router *gin.RouterGroup) {
	action := router.Group(c.Name).Use(middles.One(c.ActionSvc, c.Name))
	action.GET("/get/:id", c.One)
	action.GET("/list", c.List)
	action.GET("/list/syncfusion/schedule/", c.ListSyncfusionScheduleFormat)
	action.POST("", c.Add)
	action.POST("/list", c.ListByCondition)
	action.PUT("/:id", c.Update)
	action.DELETE("/:id", c.Delete)
}

func (c *Action) One(ctx *gin.Context) {
	action := ctx.MustGet(c.Name).(models.Action)
	resp.Success(ctx, action)
}

func (c *Action) List(ctx *gin.Context) {
	actions, err := c.ActionSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "action get error", err)
		return
	}
	resp.Success(ctx, actions)
}

func (c *Action) ListByCondition(ctx *gin.Context) {
	var action models.Action
	if err := ctx.ShouldBindJSON(&action); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	actions, err := c.ActionSvc.ListWithCondition(&action)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "action get error", err)
		return
	}
	resp.Success(ctx, actions)
}

func (c *Action) ListSyncfusionScheduleFormat(ctx *gin.Context) {
	currentDate := ctx.Query("CurrentDate")
	if len(currentDate) == 0 {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "current date not exist", nil)
		return
	}
	currentView := ctx.Query("CurrentView")
	if len(currentView) == 0 {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "current view not exist", nil)
		return
	}
	time, err := dateparse.ParseLocal(currentDate)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "current date get error", err)
		return
	}

	user := ctx.MustGet(middles.UserKey).(middles.UserInfo)
	action := models.Action{}
	action.UserId = user.Id
	switch currentView {
	case models.SyncfusionScheduleViewAgenda:
		fallthrough
	case models.SyncfusionScheduleViewDay:
		action.StartDate = now.New(time).BeginningOfDay()
		action.EndDate = now.New(time).EndOfDay()
		break
	case models.SyncfusionScheduleViewWorkWeek:
		fallthrough
	case models.SyncfusionScheduleViewWeek:
		action.StartDate = now.New(time).BeginningOfWeek()
		action.EndDate = now.New(time).EndOfWeek()
		break
	case models.SyncfusionScheduleViewMonth:
		action.StartDate = now.New(time).BeginningOfMonth()
		action.EndDate = now.New(time).EndOfMonth()
		break
	default:
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "current view not match", nil)
		return
	}
	actions, err := c.ActionSvc.ListWithCondition(&action)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "action get error", err)
		return
	}

	schedule := models.SyncfusionSchedule{}
	result := schedule.BuildActionSlice(actions)
	resp.CustomSuccess(ctx, result)
}

func (c *Action) Add(ctx *gin.Context) {
	var action models.Action
	if err := ctx.ShouldBindJSON(&action); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ActionSvc.Add(&action)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "action insert error", err)
		return
	}
	resp.Success(ctx, action)
}

func (c *Action) Update(ctx *gin.Context) {
	var action models.Action
	if err := ctx.ShouldBindJSON(&action); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ActionSvc.Update(action.Id, action)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "action update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Action) Delete(ctx *gin.Context) {
	action := ctx.MustGet(c.Name).(models.Action)
	err := c.ActionSvc.Delete(action.Id, action)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "action delete error", err)
		return
	}
	resp.Success(ctx, action.Id)
}
