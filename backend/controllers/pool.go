package controllers

import (
	"quant/backend/common"

	"github.com/gin-gonic/gin"
)

type Pool struct {
	Base
}

func NewPool() *Pool {
	Pool := &Pool{}
	Pool.Prepare()
	return Pool
}

func (c *Pool) Router(router *gin.RouterGroup) {
	pool := router.Group("pool")
	pool.GET("", c.One)
}

func (c *Pool) One(ctx *gin.Context) {
	common.ResponseSuccess(ctx, struct{}{})
}
