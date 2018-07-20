package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/time/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

type QuestResource struct {
	BaseController
}

func NewQuestResource() *QuestResource {
	QuestResource := &QuestResource{}
	QuestResource.Resource = ResourceQuestResource
	return QuestResource
}

func (c *QuestResource) Router(router *gin.RouterGroup) {
	questResource := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	questResource.GET("/get/:id", c.One)
	questResource.POST("", c.Create)
	questResource.POST("/list", c.List)
	questResource.PUT("/:id", c.Update)
	questResource.DELETE("/:id", c.Delete)
}

func (c *QuestResource) ListWithCondition(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(c.Model); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, fmt.Sprintf("%v resource json bind error", c.Resource), err)
		return
	}

	resource := c.Model.(*models.Resource)
	if resource.WithSub {
		areaIdSlice, err := c.AreaSvc.GetAllLeafId(resource.Area.Id)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v relate area list fail", c.Resource), err)
			return
		}
		areaIdSlice = append(areaIdSlice, resource.Area.Id)
		resource.Area.Ids = areaIdSlice
	}
	resourcesPtr, err := c.Service.List(resource)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v list fail", c.Resource), err)
		return
	}
	resp.Success(ctx, resourcesPtr)
}
