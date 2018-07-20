package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

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

func (s *Project) One(id int, modelPtr structs.ModelGeneral) error {
	projectJoin := models.ProjectJoin{}
	sql := s.Engine.Unscoped().Table("project").Join("INNER", "quest_target", "quest_target.id = project.quest_target_id").Join("INNER", "area", "area.id = quest_target.area_id")
	_, err := sql.Where("project.id = ?", id).Get(&projectJoin)
	projectPtr := modelPtr.(*models.Project)
	projectPtr.Area = projectJoin.Area
	projectPtr.QuestTarget = projectJoin.QuestTarget
	return err
}
