package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

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
	userResource.GET("/list", c.List)
	userResource.POST("", c.Create)
	userResource.POST("/list", c.ListByCondition)
	userResource.PUT("/:id", c.Update)
	userResource.DELETE("/:id", c.Delete)
}

func (c *UserResource) ListByCondition(ctx *gin.Context) {
	var userResource models.UserResource
	if err := ctx.ShouldBindJSON(&userResource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	userResources, err := c.UserResourceSvc.ListWithCondition(&userResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "userResource get error", err)
		return
	}
	resp.Success(ctx, userResources)
}
