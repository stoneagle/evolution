package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type UserResource struct {
	BaseController
}

func NewUserResource() *UserResource {
	UserResource := &UserResource{}
	UserResource.Resource = ResourceUserResource
	return UserResource
}

func (c *UserResource) Router(router *gin.RouterGroup) {
	userResource := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	userResource.GET("/get/:id", c.One)
	userResource.POST("", c.Create)
	userResource.POST("/list", c.List)
	userResource.PUT("/:id", c.Update)
	userResource.DELETE("/:id", c.Delete)
}
