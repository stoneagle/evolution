package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Field struct {
	ServicePackage
	structs.Service
}

func NewField(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Field {
	ret := Field{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Field) Map() (fieldsMap map[int]string, err error) {
	fields := make([]models.Field, 0)
	err = s.Engine.Find(&fields)
	if err != nil {
		return
	}
	fieldsMap = map[int]string{}
	for _, one := range fields {
		fieldsMap[one.Id] = one.Name
	}
	return
}
