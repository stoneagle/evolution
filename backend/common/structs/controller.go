package structs

import (
	"evolution/backend/common/config"
	"evolution/backend/common/logger"
)

type Controller struct {
	Resource         string
	Project          string
	Config           config.Conf
	Logger           *logger.Logger
	ResourceSvcMap   map[string]ServiceGeneral
	ResourceModelMap map[string]ModelGeneral
}

type ControllerGeneral interface {
	Prepare(config.ProjectType)
	Init()
	GetResourceServiceMap() map[string]ServiceGeneral
	GetResourceModelMap() map[string]ModelGeneral
	GetResource() string
}

func (b *Controller) Prepare(ptype config.ProjectType) {
	b.Logger = logger.Get()
	b.Config = *config.Get()
	switch ptype {
	case config.ProjectQuant:
		b.Project = b.Config.Quant.System.Name
	case config.ProjectTime:
		b.Project = b.Config.Time.System.Name
	case config.ProjectSystem:
		b.Project = b.Config.System.System.Name
	default:
		panic("project type not exist")
	}
	b.Logger.Project = b.Project
	b.Logger.Resource = b.Resource
}

func (b *Controller) GetResourceServiceMap() map[string]ServiceGeneral {
	return b.ResourceSvcMap
}

func (b *Controller) GetResourceModelMap() map[string]ModelGeneral {
	return b.ResourceModelMap
}

func (b *Controller) GetResource() string {
	return b.Resource
}
