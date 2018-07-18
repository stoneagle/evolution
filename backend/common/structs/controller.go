package structs

import (
	"evolution/backend/common/config"
	"evolution/backend/common/logger"
	"evolution/backend/common/resp"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ProjectType int8
type ResourceType string

func (rtype ResourceType) String() string {
	return string(rtype)
}

const (
	ProjectQuant ProjectType = iota + 1
	ProjectTime
	ProjectSystem
)

type Controller struct {
	Resource ResourceType
	Project  string
	Config   config.Conf
	Logger   *logger.Logger
	Service  ServiceGeneral
	Model    ModelGeneral
}

type ControllerGeneral interface {
	Prepare(ProjectType)
	GetResource() string
	GetService() ServiceGeneral
	GetModel() ModelGeneral
	SetModel(ModelGeneral)
}

type GinBaseController interface {
	ControllerGeneral
	Init()
	One(ctx *gin.Context)
	List(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func (b *Controller) Prepare(ptype ProjectType) {
	b.Logger = logger.Get()
	b.Config = *config.Get()
	switch ptype {
	case ProjectQuant:
		b.Project = b.Config.Quant.System.Name
	case ProjectTime:
		b.Project = b.Config.Time.System.Name
	case ProjectSystem:
		b.Project = b.Config.System.System.Name
	default:
		panic("project type not exist")
	}
	b.Logger.Project = b.Project
	b.Logger.Resource = b.Resource.String()
}

func (b *Controller) GetService() ServiceGeneral {
	return b.Service
}

func (b *Controller) GetModel() ModelGeneral {
	return b.Model
}

func (b *Controller) GetResource() string {
	return b.Resource.String()
}

func (b *Controller) SetModel(model ModelGeneral) {
	b.Model = model
}

func (c *Controller) One(ctx *gin.Context) {
	resp.Success(ctx, c.Model)
}

func (c *Controller) List(ctx *gin.Context) {
	resourcesPtr := c.Model.SlicePtr()
	err := c.Service.List(resourcesPtr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v list fail", c.Resource), err)
		return
	}
	resp.Success(ctx, resourcesPtr)
}

func (c *Controller) Delete(ctx *gin.Context) {
	id := ctx.MustGet("id").(int)
	err := c.Service.Delete(id, c.Model)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v delete error", c.Resource), err)
		return
	}
	resp.Success(ctx, c.Model)
}

func (c *Controller) Create(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(c.Model); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, fmt.Sprintf("%v resource json bind error", c.Resource), err)
		return
	}

	err := c.Service.Create(c.Model)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v create error", c.Resource), err)
		return
	}
	resp.Success(ctx, c.Model)
}

func (c *Controller) Update(ctx *gin.Context) {
	id := ctx.MustGet("id").(int)
	if err := ctx.ShouldBindJSON(c.Model); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, fmt.Sprintf("%v resource json bind error", c.Resource), err)
		return
	}
	err := c.Service.Update(id, c.Model)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("%v update error", c.Resource), err)
		return
	}
	resp.Success(ctx, c.Model)
}
