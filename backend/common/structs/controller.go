package structs

import (
	"evolution/backend/common/config"
	"evolution/backend/common/logger"
)

type Controller struct {
	Resource string
	Project  string
	Config   config.Conf
	Logger   *logger.Logger
}

type ControllerGeneral interface {
	Init()
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
		panic("error")
	}
	b.Logger.Log(logger.InfoLevel, b, nil)
	b.Logger.Project = b.Project
	b.Logger.Resource = b.Resource
}
