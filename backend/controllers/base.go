package controllers

import (
	"quant/backend/common"

	"github.com/go-xorm/xorm"
)

type Base struct {
	Config common.Conf
	Engine *xorm.Engine
	Cache  *redis.Client
}

func (b *Base) Prepare(ptype common.ProjectType) {
	b.Config = *common.GetConfig()
	b.Cache = common.GetRedis()
	b.Engine = common.GetEngine(b.Config.Quant.Database.Name)
}
