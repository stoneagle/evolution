package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Project struct {
	Pack ServicePackage
	structs.Service
}

func NewProject(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Project {
	ret := Project{}
	ret.Init(engine, cache, log)
	return &ret
}
