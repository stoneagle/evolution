package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Resource struct {
	BaseController
}

func NewResource() *Resource {
	Resource := &Resource{}
	Resource.Resource = ResourceResource
	return Resource
}

func (c *Resource) Router(router *gin.RouterGroup) {
	resource := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	resource.GET("/get/:id", c.One)
	resource.GET("/list", c.List)
	resource.GET("/list/areas/:id", c.ListAreas)
	resource.POST("", c.Create)
	resource.POST("/list", c.ListByCondition)
	resource.POST("/list/leaf", c.ListGroupByLeaf)
	resource.PUT("/:id", c.Update)
	resource.DELETE("/:id", c.Delete)
}

func (c *Resource) ListAreas(ctx *gin.Context) {
	resource := ctx.MustGet(c.Resource.String()).(*models.Resource)
	areas, err := c.ResourceSvc.ListAreas(resource.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "area get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "resource get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "resource get error", err)
		return
	}
	areas := c.ResourceSvc.GroupByArea(resources)
	resp.Success(ctx, areas)
}
