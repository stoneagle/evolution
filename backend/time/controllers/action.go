package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Action struct {
	BaseController
}

func NewAction() *Action {
	Action := &Action{}
	Action.Resource = ResourceAction
	return Action
}

func (c *Action) Router(router *gin.RouterGroup) {
	action := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	action.GET("/get/:id", c.One)
	action.GET("/list", c.List)
	action.POST("", c.Create)
	action.POST("/list", c.ListByCondition)
	action.PUT("/:id", c.Update)
	action.DELETE("/:id", c.Delete)
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
