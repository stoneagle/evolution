package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Treasure struct {
	structs.Controller
	TreasureSvc *services.Treasure
}

func NewTreasure() *Treasure {
	Treasure := &Treasure{}
	Treasure.Init()
	Treasure.ProjectName = Treasure.Config.Time.System.Name
	Treasure.Name = "treasure"
	Treasure.Prepare()
	Treasure.TreasureSvc = services.NewTreasure(Treasure.Engine, Treasure.Cache)
	return Treasure
}

func (c *Treasure) Router(router *gin.RouterGroup) {
	// TODO 暂时由前端控制
	// var treasure gin.IRoutes
	// switch c.Config.System.Auth.Type {
	// case middles.TypeBasicAuth:
	// 	treasure = router.Group(c.Name).Use(middles.One(c.TreasureSvc, c.Name)).Use(middles.UserFromSession(c.Config.System.Auth.Session))
	// default:
	// 	treasure = router.Group(c.Name).Use(middles.One(c.TreasureSvc, c.Name))
	// }
	treasure := router.Group(c.Name).Use(middles.One(c.TreasureSvc, c.Name))
	treasure.GET("/get/:id", c.One)
	treasure.GET("/list", c.List)
	treasure.POST("", c.Add)
	treasure.POST("/list", c.ListByCondition)
	treasure.PUT("/:id", c.Update)
	treasure.DELETE("/:id", c.Delete)
}

func (c *Treasure) One(ctx *gin.Context) {
	treasure := ctx.MustGet(c.Name).(models.Treasure)
	resp.Success(ctx, treasure)
}

func (c *Treasure) List(ctx *gin.Context) {
	treasures, err := c.TreasureSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "treasure get error", err)
		return
	}
	resp.Success(ctx, treasures)
}

func (c *Treasure) ListByCondition(ctx *gin.Context) {
	var treasure models.Treasure
	if err := ctx.ShouldBindJSON(&treasure); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	treasures, err := c.TreasureSvc.ListWithCondition(&treasure)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "treasure get error", err)
		return
	}
	resp.Success(ctx, treasures)
}

func (c *Treasure) Add(ctx *gin.Context) {
	var treasure models.Treasure
	if err := ctx.ShouldBindJSON(&treasure); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.TreasureSvc.Add(treasure)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "treasure insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Treasure) Update(ctx *gin.Context) {
	var treasure models.Treasure
	if err := ctx.ShouldBindJSON(&treasure); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.TreasureSvc.Update(treasure.Id, treasure)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "treasure update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Treasure) Delete(ctx *gin.Context) {
	treasure := ctx.MustGet(c.Name).(models.Treasure)
	err := c.TreasureSvc.Delete(treasure.Id, treasure)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "treasure delete error", err)
		return
	}
	resp.Success(ctx, treasure.Id)
}
