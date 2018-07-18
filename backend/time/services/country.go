package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Country struct {
	ServicePackage
	structs.Service
}

func NewCountry(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Country {
	ret := Country{}
	ret.Init(engine, cache, log)
	return &ret
}
