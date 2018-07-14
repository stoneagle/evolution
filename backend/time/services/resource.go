package services

import (
	"evolution/backend/common/structs"
	"evolution/backend/common/utils"
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
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	_, err = session.Insert(&model)
	if err != nil {
		session.Rollback()
		return
	}
	if model.Area.Id != 0 {
		mapAreaResource := models.MapAreaResource{}
		mapAreaResource.AreaId = model.Area.Id
		mapAreaResource.ResourceId = model.Id
		_, err = session.Insert(&mapAreaResource)
		if err != nil {
			session.Rollback()
			return
		}
	}
	if len(model.Area.Ids) > 0 {
		for _, areaId := range model.Area.Ids {
			mapAreaResource := models.MapAreaResource{}
			mapAreaResource.AreaId = areaId
			mapAreaResource.ResourceId = model.Id
			_, err = session.Insert(&mapAreaResource)
			if err != nil {
				session.Rollback()
				return
			}
		}
	}
	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return
}

func (s *Resource) Update(id int, model models.Resource) (err error) {
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	_, err = session.Id(id).Update(&model)
	if err != nil {
		session.Rollback()
		return
	}
	if len(model.Area.Ids) > 0 {
		oldMap := make([]models.MapAreaResource, 0)
		err = session.Where("map_area_resource.resource_id = ?", id).Find(&oldMap)
		if err != nil {
			session.Rollback()
			return
		}

		oldAreaIdsSlice := make([]interface{}, 0)
		newAreaIdsSlice := make([]interface{}, 0)
		for _, old := range oldMap {
			oldAreaIdsSlice = append(oldAreaIdsSlice, old.AreaId)
		}
		for _, newAreaId := range model.Area.Ids {
			newAreaIdsSlice = append(newAreaIdsSlice, newAreaId)
		}
		needAddSlice := utils.SliceDiff(newAreaIdsSlice, oldAreaIdsSlice)
		needDeleteSlice := utils.SliceDiff(oldAreaIdsSlice, newAreaIdsSlice)

		for _, areaId := range needAddSlice {
			mapAreaResource := models.MapAreaResource{}
			mapAreaResource.AreaId = areaId.(int)
			mapAreaResource.ResourceId = id
			_, err = session.Insert(&mapAreaResource)
			if err != nil {
				session.Rollback()
				return
			}
		}

		for _, areaId := range needDeleteSlice {
			mapAreaResource := models.MapAreaResource{}
			mapAreaResource.AreaId = areaId.(int)
			mapAreaResource.ResourceId = id
			_, err = session.Delete(&mapAreaResource)
			if err != nil {
				session.Rollback()
				return
			}
		}
	}
	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
	}
	return
}

func (s *Resource) Delete(id int, model models.Resource) (err error) {
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	relateTasks := make([]models.Task, 0)
	err = session.Where("resource_id = ?", id).Find(&relateTasks)
	if err != nil {
		session.Rollback()
		return
	}

	if len(relateTasks) == 0 {
		_, err = session.Unscoped().Delete(&model)
		if err != nil {
			session.Rollback()
			return
		}
		mapAreaResource := models.MapAreaResource{}
		mapAreaResource.ResourceId = id
		_, err = session.Delete(&mapAreaResource)
		if err != nil {
			session.Rollback()
			return
		}
	} else {
		_, err = session.Id(id).Delete(&model)
		if err != nil {
			session.Rollback()
			return
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return
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
	sql = sql.And("resource.deleted_at is NULL")
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

func (s *Resource) ListAreas(id int) (areas []models.Area, err error) {
	resourcesJoin := make([]models.ResourceJoin, 0)
	sql := s.Engine.Unscoped().Table("resource").Join("INNER", "map_area_resource", "map_area_resource.resource_id = resource.id").Join("INNER", "area", "area.id = map_area_resource.area_id")

	sql = sql.Where("map_area_resource.resource_id = ?", id)
	err = sql.Find(&resourcesJoin)
	if err != nil {
		return
	}
	areas = make([]models.Area, 0)
	for _, one := range resourcesJoin {
		areas = append(areas, one.Area)
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
