package controllers

import (
	"quant/backend/common"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Base
}

func NewConfig() *Config {
	Config := &Config{}
	Config.Prepare()
	return Config
}

func (c *Config) Router(router *gin.RouterGroup) {
	config := router.Group("config")
	config.GET("/type", c.Type)
}

func (c *Config) Type(ctx *gin.Context) {
	ret, err := c.Rpc.GetType()
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "get resource type error", err)
	}
	common.ResponseSuccess(ctx, ret)
}
