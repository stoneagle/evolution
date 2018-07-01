package controllers

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/common/logger"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"go.uber.org/zap"
)

type Base struct {
	Config config.Conf
	Engine *xorm.Engine
	Cache  *redis.Client
	Logger *zap.SugaredLogger
}

func (b *Base) Prepare() {
	b.Config = *config.Get()
	b.Cache = database.GetRedis()
	b.Engine = database.GetXorm(b.Config.Time.Database.Name)
	b.Logger = logger.Get()
}
