package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"fmt"

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
	userResource.POST("/list/time", c.ListTime)
	userResource.POST("/count", c.Count)
	userResource.PUT("/:id", c.Update)
	userResource.DELETE("/:id", c.Delete)
}

func (c *UserResource) ListTime(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(c.Model); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, fmt.Sprintf("%v resource json bind error", c.Resource), err)
		return
	}
	userResourcePtr := c.Model.(*models.UserResource)
	userResourcePtr.Page = &structs.Page{}
	userResourceGeneralPtr := c.Model.SlicePtr()
	err := c.Service.List(userResourcePtr, userResourceGeneralPtr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v list fail", c.Resource), err)
		return
	}

	userResourcesPtr := userResourcePtr.Transfer(userResourceGeneralPtr)
	sumTime := 0
	for _, one := range *userResourcesPtr {
		sumTime += one.Time
	}

	resp.Success(ctx, sumTime)
}
