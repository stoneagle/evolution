package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Phase struct {
	BaseController
}

func NewPhase() *Phase {
	Phase := &Phase{}
	Phase.Resource = ResourcePhase
	return Phase
}

func (c *Phase) Router(router *gin.RouterGroup) {
	phase := router.Group(c.Resource).Use(middles.OnInit(c))
	phase.GET("/get/:id", c.One)
	phase.GET("/list", c.List)
	phase.POST("", c.Add)
	phase.POST("/list", c.ListByCondition)
	phase.PUT("/:id", c.Update)
	phase.DELETE("/:id", c.Delete)
}

func (c *Phase) One(ctx *gin.Context) {
	phase := ctx.MustGet(c.Resource).(*models.Phase)
	resp.Success(ctx, phase)
}

func (c *Phase) List(ctx *gin.Context) {
	phases, err := c.PhaseSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "phase get error", err)
		return
	}
	resp.Success(ctx, phases)
}

func (c *Phase) ListByCondition(ctx *gin.Context) {
	var phase models.Phase
	if err := ctx.ShouldBindJSON(&phase); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	phases, err := c.PhaseSvc.ListWithCondition(&phase)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "phase get error", err)
		return
	}
	resp.Success(ctx, phases)
}

func (c *Phase) Add(ctx *gin.Context) {
	var phase models.Phase
	if err := ctx.ShouldBindJSON(&phase); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.PhaseSvc.Add(phase)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "phase insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Phase) Update(ctx *gin.Context) {
	var phase models.Phase
	if err := ctx.ShouldBindJSON(&phase); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.PhaseSvc.Update(phase.Id, phase)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "phase update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Phase) Delete(ctx *gin.Context) {
	phase := ctx.MustGet(c.Resource).(*models.Phase)
	err := c.PhaseSvc.Delete(phase.Id, phase)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "phase delete error", err)
		return
	}
	resp.Success(ctx, phase.Id)
}
