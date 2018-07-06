package structs

import (
	"evolution/backend/common/config"
	"evolution/backend/common/database"
	"evolution/backend/common/logger"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"go.uber.org/zap"
)

type Controller struct {
	Name        string
	ProjectName string
	Config      config.Conf
	Engine      *xorm.Engine
	Cache       *redis.Client
	Logger      *zap.SugaredLogger
}

func (b *Controller) Init() {
	b.Config = *config.Get()
	b.Logger = logger.Get()
}

func (b *Controller) Prepare() {
	b.Cache = database.GetRedis()
	b.Engine = database.GetXorm(b.ProjectName)
}
