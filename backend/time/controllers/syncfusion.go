package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
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
}

func (c *Syncfusion) ListKanban(ctx *gin.Context) {
	user := ctx.MustGet(middles.UserKey).(middles.UserInfo)
	kanbans, err := c.SyncfusionSvc.ListKanban(user.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "kanban get error", err)
		return
	}
	resp.Success(ctx, kanbans)
}
