package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
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
	userResource.PUT("/:id", c.Update)
	userResource.DELETE("/:id", c.Delete)
}

func (c *UserResource) List(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(c.Model); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, fmt.Sprintf("%v resource json bind error", c.Resource), err)
		return
	}

	userResource := c.Model.(*models.UserResource)
	if userResource.Resource.WithSub {
		areaIdSlice, err := c.AreaSvc.GetAllLeafId(userResource.Resource.Area.Id)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v relate area list fail", c.Resource), err)
			return
		}
		areaIdSlice = append(areaIdSlice, userResource.Resource.Area.Id)
		userResource.Resource.Area.Ids = areaIdSlice
	}
	resourcesPtr, err := c.Service.List(c.Model)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v list fail", c.Resource), err)
		return
	}
	resp.Success(ctx, resourcesPtr)
}
