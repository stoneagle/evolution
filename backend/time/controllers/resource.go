package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Resource struct {
	structs.Controller
	ResourceSvc *services.Resource
}

func NewResource() *Resource {
	Resource := &Resource{}
	Resource.Init()
	Resource.ProjectName = Resource.Config.Time.System.Name
	Resource.Name = "resource"
	Resource.Prepare()
	Resource.ResourceSvc = services.NewResource(Resource.Engine, Resource.Cache)
	return Resource
}

func (c *Resource) Router(router *gin.RouterGroup) {
	resource := router.Group(c.Name).Use(middles.One(c.ResourceSvc, c.Name))
	resource.GET("/get/:id", c.One)
	resource.GET("/list", c.List)
	resource.GET("/list/areas/:id", c.ListAreas)
	resource.POST("", c.Add)
	resource.POST("/list", c.ListByCondition)
	resource.POST("/list/leaf", c.ListGroupByLeaf)
	resource.PUT("/:id", c.Update)
	resource.DELETE("/:id", c.Delete)
}

func (c *Resource) One(ctx *gin.Context) {
	resource := ctx.MustGet(c.Name).(models.Resource)
	resp.Success(ctx, resource)
}

func (c *Resource) List(ctx *gin.Context) {
	resources, err := c.ResourceSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "resource get error", err)
		return
	}
	resp.Success(ctx, resources)
}

func (c *Resource) ListAreas(ctx *gin.Context) {
	resource := ctx.MustGet(c.Name).(models.Resource)
	areas, err := c.ResourceSvc.ListAreas(resource.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area get error", err)
		return
	}
	resp.Success(ctx, areas)
}

func (c *Resource) ListByCondition(ctx *gin.Context) {
	var resource models.Resource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	resources, err := c.ResourceSvc.ListWithCondition(&resource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "resource get error", err)
		return
	}
	resp.Success(ctx, resources)
}

func (c *Resource) ListGroupByLeaf(ctx *gin.Context) {
	var resource models.Resource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	resources, err := c.ResourceSvc.ListWithCondition(&resource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "resource get error", err)
		return
	}
	areas := c.ResourceSvc.GroupByArea(resources)
	resp.Success(ctx, areas)
}

func (c *Resource) Add(ctx *gin.Context) {
	var resource models.Resource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ResourceSvc.Add(resource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "resource insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Resource) Update(ctx *gin.Context) {
	var resource models.Resource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ResourceSvc.Update(resource.Id, resource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "resource update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Resource) Delete(ctx *gin.Context) {
	resource := ctx.MustGet(c.Name).(models.Resource)
	err := c.ResourceSvc.Delete(resource.Id, resource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "resource delete error", err)
		return
	}
	resp.Success(ctx, resource.Id)
}
