package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Field struct {
	Pack ServicePackage
	structs.Service
}

func NewField(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Field {
	ret := Field{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Field) Map() (fieldsMap map[int]string, err error) {
	fieldsPtr, err := s.List(&models.Field{})
	if err != nil {
		return
	}
	fields := *(fieldsPtr.(*[]models.Field))
	fieldsMap = map[int]string{}
	for _, one := range fields {
		fieldsMap[one.Id] = one.Name
	}
	return
}
