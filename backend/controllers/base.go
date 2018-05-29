package controllers

import (
	"quant/backend/common"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Base struct {
	Config common.Conf
	Engine *xorm.Engine
	Cache  *redis.Client
}

func (b *Base) Prepare() {
	b.Config = *common.GetConfig()
	b.Cache = common.GetRedis()
	b.Engine = common.GetEngine(b.Config.Quant.Database.Name)
}
