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

func (s *Field) Map() (fieldsMap map[int]models.Field, err error) {
	field := models.NewField()
	fieldsGeneralPtr := field.SlicePtr()
	err = s.List(field, fieldsGeneralPtr)
	if err != nil {
		return
	}
	fieldsMap = map[int]models.Field{}
	fieldsPtr := field.Transfer(fieldsGeneralPtr)
	for _, one := range *fieldsPtr {
		fieldsMap[one.Id] = one
	}
	return
}
