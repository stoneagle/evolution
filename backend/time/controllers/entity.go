package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Entity struct {
	structs.Controller
	EntitySvc *services.Entity
}

func NewEntity() *Entity {
	Entity := &Entity{}
	Entity.Init()
	Entity.ProjectName = Entity.Config.Time.System.Name
	Entity.Name = "entity"
	Entity.Prepare()
	Entity.EntitySvc = services.NewEntity(Entity.Engine, Entity.Cache)
	return Entity
}

func (c *Entity) Router(router *gin.RouterGroup) {
	entity := router.Group(c.Name).Use(middles.One(c.EntitySvc, c.Name))
	entity.GET("/get/:id", c.One)
	entity.GET("/list", c.List)
	entity.POST("", c.Add)
	entity.POST("/list", c.ListByCondition)
	entity.POST("/list/leaf", c.ListGroupByLeaf)
	entity.PUT("/:id", c.Update)
	entity.DELETE("/:id", c.Delete)
}

func (c *Entity) One(ctx *gin.Context) {
	entity := ctx.MustGet(c.Name).(models.Entity)
	resp.Success(ctx, entity)
}

func (c *Entity) List(ctx *gin.Context) {
	entities, err := c.EntitySvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entity get error", err)
		return
	}
	resp.Success(ctx, entities)
}

func (c *Entity) ListByCondition(ctx *gin.Context) {
	var entity models.Entity
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	entities, err := c.EntitySvc.ListWithCondition(&entity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entity get error", err)
		return
	}
	resp.Success(ctx, entities)
}

func (c *Entity) ListGroupByLeaf(ctx *gin.Context) {
	var entity models.Entity
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	entities, err := c.EntitySvc.ListWithCondition(&entity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entity get error", err)
		return
	}
	areas := c.EntitySvc.GroupByArea(entities)
	resp.Success(ctx, areas)
}

func (c *Entity) Add(ctx *gin.Context) {
	var entity models.Entity
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.EntitySvc.Add(entity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entity insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Entity) Update(ctx *gin.Context) {
	var entity models.Entity
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.EntitySvc.Update(entity.Id, entity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entity update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Entity) Delete(ctx *gin.Context) {
	entity := ctx.MustGet(c.Name).(models.Entity)
	err := c.EntitySvc.Delete(entity.Id, entity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "entity delete error", err)
		return
	}
	resp.Success(ctx, entity.Id)
}
