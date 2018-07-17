package controllers

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/common/structs"
	"evolution/backend/system/services"
)

const (
	ResourceUser = "user"
)

type Base struct {
	structs.Controller
	services.Base
}

func (c *Base) Init() {
	cache := database.GetRedis()
	engine := database.GetXorm(c.Project)
	c.Prepare(config.ProjectSystem)
	c.PrepareService(engine, cache, c.Logger)
}
