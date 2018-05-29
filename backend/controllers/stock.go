package controllers

import (
	"quant/backend/common"

	"github.com/gin-gonic/gin"
)

type Stock struct {
	Base
}

func NewStock() *Stock {
	Stock := &Stock{}
	Stock.Prepare()
	return Stock
}

func (c *Stock) Router(router *gin.RouterGroup) {
	stock := router.Group("stock")
	stock.GET("", c.One)
}

func (c *Stock) One(ctx *gin.Context) {
	common.ResponseSuccess(ctx, struct{}{})
}
