package controllers

import (
	"quant/backend/common"
	"quant/backend/rpc/engine"

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
	config.GET("/type/:asset", c.Type)
	config.GET("/asset", c.Asset)
}

func (c *Config) Type(ctx *gin.Context) {
	asset := ctx.Param("asset")
	atype, err := engine.AssetTypeFromString(asset)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "asset type illegal", err)
		return
	}
	ret, err := c.Rpc.GetType(atype)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "get asset type error", err)
		return
	}
	common.ResponseSuccess(ctx, ret)
}

func (c *Config) Asset(ctx *gin.Context) {
	asset := map[engine.AssetType]string{
		engine.AssetType_Stock:    "STOCK",
		engine.AssetType_Exchange: "EXCHANGE",
	}
	common.ResponseSuccess(ctx, asset)
}
