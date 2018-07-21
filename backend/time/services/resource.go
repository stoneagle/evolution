package services

import (
	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/common/utils"
	"evolution/backend/time/models"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Resource struct {
	Pack ServicePackage
	structs.Service
}

func NewResource(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Resource {
	ret := Resource{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Resource) Create(modelPtr structs.ModelGeneral) (err error) {
	resourcePtr := modelPtr.(*models.Resource)
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	_, err = session.Insert(&resourcePtr)
	if err != nil {
		session.Rollback()
		return
	}
	if resourcePtr.Area.Id != 0 {
		mapAreaResource := models.NewMapAreaResource()
		mapAreaResource.AreaId = resourcePtr.Area.Id
		mapAreaResource.ResourceId = resourcePtr.Id
		_, err = session.Insert(mapAreaResource)
		if err != nil {
			session.Rollback()
			return
		}
	}
	if len(resourcePtr.Area.Ids) > 0 {
		for _, areaId := range resourcePtr.Area.Ids {
			mapAreaResource := models.NewMapAreaResource()
			mapAreaResource.AreaId = areaId
			mapAreaResource.ResourceId = resourcePtr.Id
			_, err = session.Insert(mapAreaResource)
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

func (s *Resource) Update(id int, modelPtr structs.ModelGeneral) (err error) {
	resourcePtr := modelPtr.(*models.Resource)
	session := s.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	_, err = session.Id(id).Update(resourcePtr)
	if err != nil {
		session.Rollback()
		return
	}
	if len(resourcePtr.Area.Ids) > 0 {
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
		for _, newAreaId := range resourcePtr.Area.Ids {
			newAreaIdsSlice = append(newAreaIdsSlice, newAreaId)
		}
		needAddSlice := utils.SliceDiff(newAreaIdsSlice, oldAreaIdsSlice)
		needDeleteSlice := utils.SliceDiff(oldAreaIdsSlice, newAreaIdsSlice)

		for _, areaId := range needAddSlice {
			mapAreaResource := models.NewMapAreaResource()
			mapAreaResource.AreaId = areaId.(int)
			mapAreaResource.ResourceId = id
			_, err = session.Insert(&mapAreaResource)
			if err != nil {
				session.Rollback()
				return
			}
		}

		for _, areaId := range needDeleteSlice {
			mapAreaResource := models.NewMapAreaResource()
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

func (s *Resource) Delete(id int, modelPtr structs.ModelGeneral) (err error) {
	resourcePtr := modelPtr.(*models.Resource)
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
		_, err = session.Unscoped().Delete(resourcePtr)
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
		_, err = session.Id(id).Delete(resourcePtr)
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

func (s *Resource) GroupByArea(resourcesPtr *[]models.Resource) (areas []models.Area) {
	areasMap := map[int]*models.Area{}
	for _, one := range *resourcesPtr {
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
		areas = append(areas, *one)
	}
	return areas
}
