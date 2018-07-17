package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Country struct {
	BaseController
}

func NewCountry() *Country {
	country := &Country{}
	country.Resource = ResourceCountry
	return country
}

func (c *Country) Router(router *gin.RouterGroup) {
	country := router.Group(c.Resource).Use(middles.OnInit(c))
	country.GET("/get/:id", c.One)
	country.GET("/list", c.List)
	country.POST("", c.Add)
	country.PUT("/:id", c.Update)
	country.DELETE("/:id", c.Delete)
}

func (c *Country) One(ctx *gin.Context) {
	country := ctx.MustGet(c.Resource).(*models.Country)
	resp.Success(ctx, country)
}

func (c *Country) List(ctx *gin.Context) {
	countries, err := c.CountrySvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "country get error", err)
		return
	}
	resp.Success(ctx, countries)
}

func (c *Country) Add(ctx *gin.Context) {
	var country models.Country
	if err := ctx.ShouldBindJSON(&country); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.CountrySvc.Add(country)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "country insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Country) Update(ctx *gin.Context) {
	var country models.Country
	if err := ctx.ShouldBindJSON(&country); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.CountrySvc.Update(country.Id, country)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "country update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Country) Delete(ctx *gin.Context) {
	country := ctx.MustGet(c.Resource).(*models.Country)
	err := c.CountrySvc.Delete(country.Id, country)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "country delete error", err)
		return
	}
	resp.Success(ctx, country.Id)
}
