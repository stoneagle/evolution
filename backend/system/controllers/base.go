package controllers

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/common/structs"
	"evolution/backend/system/models"
	"evolution/backend/system/services"
)

const (
	ResourceUser = "user"
)

type BaseController struct {
	structs.Controller
	services.ServicePackage
	models.ModelPackage
}

func (c *BaseController) Init() {
	c.Prepare(config.ProjectSystem)
	cache := database.GetRedis()
	engine := database.GetXorm(c.Project)
	c.PrepareService(engine, cache, c.Logger)
	c.PrepareModel()
	c.ResourceSvcMap = map[string]structs.ServiceGeneral{
		ResourceUser: c.UserSvc,
	}
	c.ResourceModelMap = map[string]structs.ModelGeneral{
		ResourceUser: c.UserModel,
	}
}
