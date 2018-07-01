package controllers

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/common/logger"
	"evolution/backend/quant/rpc"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"go.uber.org/zap"
)

type Base struct {
	Config config.Conf
	Engine *xorm.Engine
	Cache  *redis.Client
	Logger *zap.SugaredLogger
	Rpc    *rpc.Rpc
}

func (b *Base) Prepare() {
	b.Config = *config.Get()
	b.Cache = database.GetRedis()
	b.Engine = database.GetXorm(b.Config.Quant.Database.Name)
	b.Rpc = rpc.NewRpc(b.Config.Quant.Rpc.Host, b.Config.Quant.Rpc.Port)
	b.Logger = logger.Get()
}
