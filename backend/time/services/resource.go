package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Resource struct {
	structs.Service
}

func NewResource(engine *xorm.Engine, cache *redis.Client) *Resource {
	ret := Resource{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Resource) One(id int) (interface{}, error) {
	model := models.Resource{}
	_, err := s.Engine.Where("id = ?", id).Get(&model)
	return model, err
}

func (s *Resource) Add(model models.Resource) (err error) {
	_, err = s.Engine.Insert(&model)
	if err != nil {
		return
	}
	mapAreaResource := models.MapAreaResource{}
	mapAreaResource.AreaId = model.Area.Id
	mapAreaResource.ResourceId = model.Id
	_, err = s.Engine.Insert(&mapAreaResource)
	return
}

func (s *Resource) Update(id int, model models.Resource) (err error) {
	_, err = s.Engine.Id(id).Update(&model)
	if err != nil {
		return
	}
	if model.Area.Id != 0 {
		mapAreaResource := models.MapAreaResource{}
		mapAreaResource.AreaId = model.Area.Id
		_, err = s.Engine.Where("resource_id = ?", id).Update(&mapAreaResource)
		if err != nil {
			return
		}
	}
	return
}

func (s *Resource) Delete(id int, model models.Resource) (err error) {
	_, err = s.Engine.Id(id).Get(&model)
	if err == nil {
		_, err = s.Engine.Id(id).Delete(&model)
	}
	return
}

func (s *Resource) ListWithCondition(resource *models.Resource) (resources []models.Resource, err error) {
	resourcesJoin := make([]models.ResourceJoin, 0)
	sql := s.Engine.Unscoped().Table("resource").Join("INNER", "map_area_resource", "map_area_resource.resource_id = resource.id").Join("INNER", "area", "area.id = map_area_resource.area_id")

	condition := resource.BuildCondition()
	if resource.WithSub {
		areaIdSlice, err := NewArea(s.Engine, s.Cache).GetAllLeafId(resource.Area.Id)
		if err != nil {
			return resources, err
		}
		areaIdSlice = append(areaIdSlice, resource.Area.Id)
		condition["area.id"] = areaIdSlice
	}

	sql = sql.Where(condition)
	err = sql.Find(&resourcesJoin)
	if err != nil {
		return
	}
	resources = make([]models.Resource, 0)
	for _, one := range resourcesJoin {
		one.Resource.Area = one.Area
		resources = append(resources, one.Resource)
	}
	return
}

func (s *Resource) List() (resources []models.Resource, err error) {
	resources = make([]models.Resource, 0)
	err = s.Engine.Find(&resources)
	return
}

func (s *Resource) GroupByArea(resources []models.Resource) (areas []models.Area) {
	areasMap := map[int]models.Area{}
	for _, one := range resources {
		if _, ok := areasMap[one.Area.Id]; !ok {
			one.Area.Resources = make([]models.Resource, 0)
			areasMap[one.Area.Id] = one.Area
		}
		tmp, _ := areasMap[one.Area.Id]
		tmp.Resources = append(tmp.Resources, one)
		areasMap[one.Area.Id] = tmp
	}

	areas = make([]models.Area, 0)
	for _, one := range areasMap {
		areas = append(areas, one)
	}
	return areas
}
