package controllers

import (
	"evolution/backend/common/resp"
	"evolution/backend/time/middles"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type EntityLife struct {
	Base
	Name          string
	EntityLifeSvc *services.EntityLife
}

func NewEntityLife() *EntityLife {
	EntityLife := &EntityLife{
		Name: "entityLife",
	}
	EntityLife.Prepare()
	EntityLife.EntityLifeSvc = services.NewEntityLife(EntityLife.Engine, EntityLife.Cache)
	return EntityLife
}

func (c *EntityLife) Router(router *gin.RouterGroup) {
	entityLife := router.Group("life").Use(middles.One(c.EntityLifeSvc, c.Name))
	entityLife.GET("/get/:id", c.One)
	entityLife.GET("/list", c.List)
	entityLife.POST("", c.Add)
	entityLife.PUT("/:id", c.Update)
	entityLife.DELETE("/:id", c.Delete)
}

func (c *EntityLife) One(ctx *gin.Context) {
	entityLife := ctx.MustGet("entityLife").(models.EntityLife)
	resp.Success(ctx, entityLife)
}

func (c *EntityLife) List(ctx *gin.Context) {
	entityLifes, err := c.EntityLifeSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entityLife get error", err)
		return
	}
	resp.Success(ctx, entityLifes)
}

func (c *EntityLife) Add(ctx *gin.Context) {
	var entityLife models.EntityLife
	if err := ctx.ShouldBindJSON(&entityLife); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.EntityLifeSvc.Add(entityLife)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entityLife insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *EntityLife) Update(ctx *gin.Context) {
	var entityLife models.EntityLife
	if err := ctx.ShouldBindJSON(&entityLife); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.EntityLifeSvc.Update(entityLife.Id, entityLife)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entityLife update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *EntityLife) Delete(ctx *gin.Context) {
	entityLife := ctx.MustGet("entityLife").(models.EntityLife)
	err := c.EntityLifeSvc.Delete(entityLife.Id, entityLife)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entityLife delete error", err)
		return
	}
	resp.Success(ctx, entityLife.Id)
}
