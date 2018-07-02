package controllers

import (
	"evolution/backend/common/resp"
	"evolution/backend/time/models"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Base
	Name string
}

func NewConfig() *Config {
	Config := &Config{
		Name: "config",
	}
	Config.Prepare()
	return Config
}

func (c *Config) Router(router *gin.RouterGroup) {
	config := router.Group("config")
	config.GET("/field", c.FieldMap)
}

func (c *Config) FieldMap(ctx *gin.Context) {
	resp.Success(ctx, models.AreaFieldMap)
}
