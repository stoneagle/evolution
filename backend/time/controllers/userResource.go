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
	userResource := router.Group(c.Resource).Use(middles.OnInit(c))
	userResource.GET("/get/:id", c.One)
	userResource.GET("/list", c.List)
	userResource.POST("", c.Add)
	userResource.POST("/list", c.ListByCondition)
	userResource.PUT("/:id", c.Update)
	userResource.DELETE("/:id", c.Delete)
}

func (c *UserResource) One(ctx *gin.Context) {
	userResource := ctx.MustGet(c.Resource).(*models.UserResource)
	resp.Success(ctx, userResource)
}

func (c *UserResource) List(ctx *gin.Context) {
	userResources, err := c.UserResourceSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "userResource get error", err)
		return
	}
	resp.Success(ctx, userResources)
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

func (c *UserResource) Add(ctx *gin.Context) {
	var userResource models.UserResource
	if err := ctx.ShouldBindJSON(&userResource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.UserResourceSvc.Add(userResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "userResource insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *UserResource) Update(ctx *gin.Context) {
	var userResource models.UserResource
	if err := ctx.ShouldBindJSON(&userResource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.UserResourceSvc.Update(userResource.Id, userResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "userResource update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *UserResource) Delete(ctx *gin.Context) {
	userResource := ctx.MustGet(c.Resource).(*models.UserResource)
	err := c.UserResourceSvc.Delete(userResource.Id, userResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "userResource delete error", err)
		return
	}
	resp.Success(ctx, userResource.Id)
}
