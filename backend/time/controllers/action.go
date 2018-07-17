package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Action struct {
	Base
}

func NewAction() *Action {
	Action := &Action{}
	Action.Resource = ResourceAction
	return Action
}

func (c *Action) Router(router *gin.RouterGroup) {
	action := router.Group(c.Resource).Use(middles.OnInit(c)).Use(middles.One(c.ActionSvc, c.Resource, &models.Action{}))
	action.GET("/get/:id", c.One)
	action.GET("/list", c.List)
	action.POST("", c.Add)
	action.POST("/list", c.ListByCondition)
	action.PUT("/:id", c.Update)
	action.DELETE("/:id", c.Delete)
}

func (c *Action) One(ctx *gin.Context) {
	one := ctx.MustGet(c.Resource).(models.Action)
	resp.Success(ctx, one)
}

func (c *Action) List(ctx *gin.Context) {
	actions, err := c.ActionSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "action get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "action get error", err)
		return
	}
	resp.Success(ctx, actions)
}

func (c *Action) Add(ctx *gin.Context) {
	var action models.Action
	if err := ctx.ShouldBindJSON(&action); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ActionSvc.Add(&action)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "action insert error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "action update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Action) Delete(ctx *gin.Context) {
	action := ctx.MustGet(c.Resource).(*models.Action)
	err := c.ActionSvc.Delete(action.Id, action)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "action delete error", err)
		return
	}
	resp.Success(ctx, action.Id)
}
