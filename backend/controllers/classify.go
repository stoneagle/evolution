package controllers

import (
	"quant/backend/common"

	"github.com/gin-gonic/gin"
)

type Classify struct {
	Base
}

func NewClassify() *Classify {
	Classify := &Classify{}
	Classify.Prepare()
	return Classify
}

func (c *Classify) Router(router *gin.RouterGroup) {
	stock := router.Group("classify")
	stock.GET("", c.One)
}

func (c *Classify) One(ctx *gin.Context) {
	common.ResponseSuccess(ctx, struct{}{})
}
