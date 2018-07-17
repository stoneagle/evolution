package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type UserResource struct {
	ServicePackage
	structs.Service
}

func NewUserResource(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *UserResource {
	ret := UserResource{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *UserResource) Add(model models.UserResource) (err error) {
	_, err = s.Engine.Insert(&model)
	return
}

func (s *UserResource) Update(id int, model models.UserResource) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	return
}

func (s *UserResource) List() (userResources []models.UserResource, err error) {
	userResources = make([]models.UserResource, 0)
	err = s.Engine.Find(&userResources)
	return
}

func (s *UserResource) ListWithCondition(userResource *models.UserResource) (userResources []models.UserResource, err error) {
	userResourcesJoin := make([]models.UserResourceJoin, 0)
	sql := s.Engine.Unscoped().Table("user_resource").Join("INNER", "resource", "resource.id = user_resource.resource_id").Join("INNER", "map_area_resource", "map_area_resource.resource_id = resource.id").Join("INNER", "area", "area.id = map_area_resource.area_id")

	condition := userResource.BuildCondition()
	if userResource.Resource.WithSub {
		areaIdSlice, err := s.AreaSvc.GetAllLeafId(userResource.Resource.Area.Id)
		if err != nil {
			return userResources, err
		}
		areaIdSlice = append(areaIdSlice, userResource.Resource.Area.Id)
		condition["area.id"] = areaIdSlice
	}

	sql = sql.Where(condition)
	err = sql.Find(&userResourcesJoin)
	if err != nil {
		return
	}

	userResources = make([]models.UserResource, 0)
	for _, one := range userResourcesJoin {
		one.UserResource.Resource = one.Resource
		one.UserResource.Resource.Area = one.Area
		userResources = append(userResources, one.UserResource)
	}
	return
}
