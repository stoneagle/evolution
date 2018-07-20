package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type Country struct {
	BaseController
}

func NewCountry() *Country {
	country := &Country{}
	country.Resource = ResourceCountry
	return country
}

func (c *Country) Router(router *gin.RouterGroup) {
	country := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	country.GET("/get/:id", c.One)
	country.GET("/list", c.List)
	country.POST("", c.Create)
	country.POST("/list", c.ListWithCondition)
	country.PUT("/:id", c.Update)
	country.DELETE("/:id", c.Delete)
}
