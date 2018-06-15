package controllers

import (
	"quant/backend/common"
	"quant/backend/models"
	"quant/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Item struct {
	Base
	ItemSvc *services.Item
}

func NewItem() *Item {
	Item := &Item{}
	Item.Prepare()
	Item.ItemSvc = services.NewItem(Item.Engine, Item.Cache)
	return Item
}

func (c *Item) Router(router *gin.RouterGroup) {
	item := router.Group("item")
	item.GET("/:id", initItem(c.ItemSvc), c.One)
	item.GET("", c.List)
}

func initItem(svc *services.Item) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorParams, "id params error", err)
		}

		item, err := svc.One(id)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorMysql, "get item error", err)
		}

		ctx.Set("item", item)
		return
	}
}

func (c *Item) One(ctx *gin.Context) {
	item := ctx.MustGet("item").(models.Item)
	common.ResponseSuccess(ctx, item)
}

func (c *Item) List(ctx *gin.Context) {
	items, err := c.ItemSvc.List()
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "item get error", err)
		return
	}
	common.ResponseSuccess(ctx, items)
}
