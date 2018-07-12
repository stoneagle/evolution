package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type UserResource struct {
	structs.Controller
	UserResourceSvc *services.UserResource
}

func NewUserResource() *UserResource {
	UserResource := &UserResource{}
	UserResource.Init()
	UserResource.ProjectName = UserResource.Config.Time.System.Name
	UserResource.Name = "user-resource"
	UserResource.Prepare()
	UserResource.UserResourceSvc = services.NewUserResource(UserResource.Engine, UserResource.Cache)
	return UserResource
}

func (c *UserResource) Router(router *gin.RouterGroup) {
	userResource := router.Group(c.Name).Use(middles.One(c.UserResourceSvc, c.Name))
	userResource.GET("/get/:id", c.One)
	userResource.GET("/list", c.List)
	userResource.POST("", c.Add)
	userResource.POST("/list", c.ListByCondition)
	userResource.PUT("/:id", c.Update)
	userResource.DELETE("/:id", c.Delete)
}

func (c *UserResource) One(ctx *gin.Context) {
	userResource := ctx.MustGet(c.Name).(models.UserResource)
	resp.Success(ctx, userResource)
}

func (c *UserResource) List(ctx *gin.Context) {
	userResources, err := c.UserResourceSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "userResource get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "userResource get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "userResource insert error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "userResource update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *UserResource) Delete(ctx *gin.Context) {
	userResource := ctx.MustGet(c.Name).(models.UserResource)
	err := c.UserResourceSvc.Delete(userResource.Id, userResource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "userResource delete error", err)
		return
	}
	resp.Success(ctx, userResource.Id)
}
