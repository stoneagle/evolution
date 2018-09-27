package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Quest struct {
	Pack ServicePackage
	structs.Service
}

func NewQuest(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Quest {
	ret := Quest{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Quest) Update(id int, modelPtr structs.ModelGeneral) (err error) {
	questPtr := modelPtr.(*models.Quest)
	err = s.Service.Update(id, questPtr)
	if err != nil {
		return
	}

	if questPtr.StartDateReset {
		err = s.Service.UpdateByMap(questPtr.TableName(), id, map[string]interface{}{
			"start_date": nil,
		})
		if err != nil {
			return
		}
	}

	if questPtr.EndDateReset {
		err = s.Service.UpdateByMap(questPtr.TableName(), id, map[string]interface{}{
			"end_date": nil,
		})
		if err != nil {
			return
		}
	}
	return
}
