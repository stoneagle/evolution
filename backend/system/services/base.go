package services

import (
	"evolution/backend/common/logger"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Base struct {
	UserSvc *User
}

func (s *Base) PrepareService(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) {
	s.UserSvc = NewUser(engine, cache, log)
}
