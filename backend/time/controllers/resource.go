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
	resource.GET("/list/areas/:id", c.ListAreas)
	resource.POST("", c.Create)
	resource.POST("/list", c.List)
	resource.POST("/count", c.Count)
	resource.POST("/list/leaf", c.ListGroupByLeaf)
	resource.PUT("/:id", c.Update)
	resource.DELETE("/:id", c.Delete)
}

func (c *Resource) ListAreas(ctx *gin.Context) {
	resourcePtr := c.Model.(*models.Resource)
	areas, err := c.ResourceSvc.ListAreas(resourcePtr.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "area get error", err)
		return
	}
	resp.Success(ctx, areas)
}

func (c *Resource) ListGroupByLeaf(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(c.Model); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	resourcesGeneralPtr := c.Model.SlicePtr()
	err := c.ResourceSvc.List(c.Model, resourcesGeneralPtr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "resource get error", err)
		return
	}
	resourcesPtr := c.ResourceModel.Transfer(resourcesGeneralPtr)
	areas := c.ResourceSvc.GroupByArea(resourcesPtr)
	resp.Success(ctx, areas)
}
