package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Field struct {
	structs.Controller
	FieldSvc *services.Field
}

func NewField() *Field {
	Field := &Field{}
	Field.Init()
	Field.ProjectName = Field.Config.Time.System.Name
	Field.Name = "field"
	Field.Prepare()
	Field.FieldSvc = services.NewField(Field.Engine, Field.Cache)
	return Field
}

func (c *Field) Router(router *gin.RouterGroup) {
	field := router.Group(c.Name).Use(middles.One(c.FieldSvc, c.Name))
	field.GET("/get/:id", c.One)
	field.GET("/list", c.List)
	field.POST("", c.Add)
	field.PUT("/:id", c.Update)
	field.DELETE("/:id", c.Delete)
}

func (c *Field) One(ctx *gin.Context) {
	field := ctx.MustGet(c.Name).(models.Field)
	resp.Success(ctx, field)
}

func (c *Field) List(ctx *gin.Context) {
	fields, err := c.FieldSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "field get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "field insert error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "field update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Field) Delete(ctx *gin.Context) {
	field := ctx.MustGet(c.Name).(models.Field)
	err := c.FieldSvc.Delete(field.Id, field)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "field delete error", err)
		return
	}
	resp.Success(ctx, field.Id)
}
