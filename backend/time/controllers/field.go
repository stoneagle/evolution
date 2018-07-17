package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Field struct {
	Base
}

func NewField() *Field {
	Field := &Field{}
	Field.Resource = ResourceField
	return Field
}

func (c *Field) Router(router *gin.RouterGroup) {
	field := router.Group(c.Resource).Use(middles.OnInit(c)).Use(middles.One(c.FieldSvc, c.Resource, models.Field{}))
	field.GET("/get/:id", c.One)
	field.GET("/list", c.List)
	field.POST("", c.Add)
	field.PUT("/:id", c.Update)
	field.DELETE("/:id", c.Delete)
}

func (c *Field) One(ctx *gin.Context) {
	field := ctx.MustGet(c.Resource).(*models.Field)
	resp.Success(ctx, field)
}

func (c *Field) List(ctx *gin.Context) {
	fields, err := c.FieldSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "field get error", err)
		return
	}
	resp.Success(ctx, fields)
}

func (c *Field) Add(ctx *gin.Context) {
	var field models.Field
	if err := ctx.ShouldBindJSON(&field); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.FieldSvc.Add(field)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "field insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Field) Update(ctx *gin.Context) {
	var field models.Field
	if err := ctx.ShouldBindJSON(&field); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.FieldSvc.Update(field.Id, field)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "field update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Field) Delete(ctx *gin.Context) {
	field := ctx.MustGet(c.Resource).(*models.Field)
	err := c.FieldSvc.Delete(field.Id, field)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "field delete error", err)
		return
	}
	resp.Success(ctx, field.Id)
}
