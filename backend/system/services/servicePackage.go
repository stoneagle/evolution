package services

import (
	"evolution/backend/common/logger"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type ServicePackage struct {
	UserSvc *User
}

func (s *ServicePackage) PrepareService(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) {
	s.UserSvc = NewUser(engine, cache, log)
}
