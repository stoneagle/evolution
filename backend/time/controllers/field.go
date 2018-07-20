package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type Field struct {
	BaseController
}

func NewField() *Field {
	Field := &Field{}
	Field.Resource = ResourceField
	return Field
}

func (c *Field) Router(router *gin.RouterGroup) {
	field := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	field.GET("/get/:id", c.One)
	field.GET("/list", c.List)
	field.POST("", c.Create)
	field.POST("/list", c.ListWithCondition)
	field.PUT("/:id", c.Update)
	field.DELETE("/:id", c.Delete)
}
