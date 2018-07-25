package controllers

import (
	"evolution/backend/common/database"
	"evolution/backend/common/structs"
	"evolution/backend/system/models"
	"evolution/backend/system/services"
	"fmt"
)

const (
	ResourceUser structs.ResourceType = "user"
)

type BaseController struct {
	structs.Controller
	services.ServicePackage
	models.ModelPackage
}

func (c *BaseController) Init() {
	c.Prepare(structs.ProjectSystem)
	cache := database.GetRedis()
	engine := database.GetXorm(c.Project)
	c.PrepareService(engine, cache, c.Logger)
	c.PrepareModel()
	c.ChangeSvc(c.Resource)
	c.ChangeModel(c.Resource)
}

func (c *BaseController) ChangeSvc(resource structs.ResourceType) {
	svcMap := map[structs.ResourceType]structs.ServiceGeneral{
		ResourceUser: c.UserSvc,
	}
	svc, ok := svcMap[resource]
	if !ok {
		panic(fmt.Sprintf("%v svc not exist", c.Resource))
	}
	c.Service = svc
}

func (c *BaseController) ChangeModel(resource structs.ResourceType) {
	modelMap := map[structs.ResourceType]structs.ModelGeneral{
		ResourceUser: c.UserModel,
	}
	model, ok := modelMap[resource]
	if !ok {
		panic(fmt.Sprintf("%v model not exist", c.Resource))
	}
	c.Model = model
}
