package controllers

import (
	"github.com/gin-gonic/gin"
)

type Signal struct {
	Base
}

func NewSignal() *Signal {
	Signal := &Signal{}
	Signal.Prepare()
	return Signal
}

func (c *Signal) Router(router *gin.RouterGroup) {
	signal := router.Group("signal")
	signal.GET("", c.One)
}

func (c *Signal) One(ctx *gin.Context) {
}
