package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

type Syncfusion struct {
	structs.Controller
	SyncfusionSvc *services.Syncfusion
}

func NewSyncfusion() *Syncfusion {
	Syncfusion := &Syncfusion{}
	Syncfusion.Init()
	Syncfusion.ProjectName = Syncfusion.Config.Time.System.Name
	Syncfusion.Name = "syncfusion"
	Syncfusion.Prepare()
	Syncfusion.SyncfusionSvc = services.NewSyncfusion(Syncfusion.Engine, Syncfusion.Cache)
	return Syncfusion
}

func (c *Syncfusion) Router(router *gin.RouterGroup) {
	syncfusion := router.Group(c.Name)
	syncfusion.GET("/list/kanban", c.ListKanban)
	syncfusion.GET("/list/gantt", c.ListGantt)
	syncfusion.GET("/list/schedule/", c.ListSchedule)
	syncfusion.GET("/list/treegrid/:fieldId/", c.ListTreeGrid)
}

func (c *Syncfusion) ListKanban(ctx *gin.Context) {
	user := ctx.MustGet(middles.UserKey).(middles.UserInfo)
	kanbans, err := c.SyncfusionSvc.ListKanban(user.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "task kanban get error", err)
		return
	}
	resp.Success(ctx, kanbans)
}

func (c *Syncfusion) ListSchedule(ctx *gin.Context) {
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
	currentTime, err := dateparse.ParseLocal(currentDate)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "current date get error", err)
		return
	}

	user := ctx.MustGet(middles.UserKey).(middles.UserInfo)
	var startDate time.Time
	var endDate time.Time
	switch currentView {
	case models.SyncfusionScheduleViewAgenda:
		fallthrough
	case models.SyncfusionScheduleViewDay:
		startDate = now.New(currentTime).BeginningOfDay()
		endDate = now.New(currentTime).EndOfDay()
		break
	case models.SyncfusionScheduleViewWorkWeek:
		fallthrough
	case models.SyncfusionScheduleViewWeek:
		startDate = now.New(currentTime).BeginningOfWeek()
		endDate = now.New(currentTime).EndOfWeek()
		break
	case models.SyncfusionScheduleViewMonth:
		startDate = now.New(currentTime).BeginningOfMonth()
		endDate = now.New(currentTime).EndOfMonth()
		break
	default:
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "current view not match", nil)
		return
	}
	schedules, err := c.SyncfusionSvc.ListSchedule(user.Id, startDate, endDate)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "action schedules get error", nil)
		return
	}
	resp.CustomSuccess(ctx, schedules)
}

func (c *Syncfusion) ListTreeGrid(ctx *gin.Context) {
	fieldIdStr := ctx.Param("fieldId")
	var fieldId int
	var err error
	if fieldIdStr != "" {
		fieldId, err = strconv.Atoi(fieldIdStr)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "field id params error", err)
			return
		}
	}

	parentId := 0
	filter := ctx.Query("$filter")
	if len(filter) > 0 {
		s := strings.Split(filter, " eq ")
		if len(s) < 2 {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "parent id params error", nil)
			return
		}
		parentIdStr := s[1]
		if parentIdStr != "null" {
			parentId, err = strconv.Atoi(parentIdStr)
			if err != nil {
				resp.ErrorBusiness(ctx, resp.ErrorParams, "parent id params transfer error", err)
				return
			}
		}
	}

	treeGrids, err := c.SyncfusionSvc.ListTreeGrid(fieldId, parentId)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area treegrid get error", err)
		return
	}
	res := map[string]interface{}{}
	res["result"] = treeGrids
	res["__count"] = len(treeGrids)
	resp.CustomSuccess(ctx, res)
}

func (c *Syncfusion) ListGantt(ctx *gin.Context) {
	user := ctx.MustGet(middles.UserKey).(middles.UserInfo)
	gantts, err := c.SyncfusionSvc.ListGantt(user.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest to task gantt get error", err)
		return
	}
	resp.Success(ctx, gantts)
}
