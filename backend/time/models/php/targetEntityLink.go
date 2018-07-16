package php

import (
	"errors"
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type TargetEntityLink struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	TargetId int       `xorm:"not null default 0 comment('领域对象ID') INT(11)"`
	EntityId int       `xorm:"not null default 0 comment('相关实体ID') INT(11)"`
	Ctime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

type TargetEntityLinkJoin struct {
	TargetEntityLink `xorm:"extends"`
	Target           `xorm:"extends"`
}

func (c *TargetEntityLink) Transfer(src, des *xorm.Engine, userId int) {
	olds := make([]TargetEntityLinkJoin, 0)
	err := src.Table("target_entity_link").Join("LEFT", "target", "target.id = target_entity_link.target_id").Find(&olds)
	if err != nil {
		fmt.Printf("target join get error:%v\r\n", err.Error())
		return
	}

	fieldPhaseMap := map[int]int{}
	insertNum := 0

	for _, one := range olds {
		field := one.Target.FieldId
		entityId := one.TargetEntityLink.EntityId

		_, ok := fieldPhaseMap[field]
		if !ok {
			phases := make([]models.Phase, 0)
			des.Where("field_id = ?", field).Asc("level").Find(&phases)
			if len(phases) == 0 {
				fmt.Printf("phase not exist")
				return
			}
			fieldPhaseMap[field] = phases[0].Id
		}

		// 找到新版本resource的Id
		newResourceJoin, err := getResourceJoin(src, des, entityId, field)
		if err != nil {
			fmt.Printf("get new resource error:%v\r\n", err)
			return
		}

		userResource := models.UserResource{}
		userResource.UserId = userId
		userResource.ResourceId = newResourceJoin.Resource.Id
		userResource.CreatedAt = one.TargetEntityLink.Ctime
		userResource.UpdatedAt = one.TargetEntityLink.Utime

		has, err := des.Where("user_id = ?", userId).And("resource_id = ?", newResourceJoin.Resource.Id).Get(new(models.UserResource))
		if err != nil {
			fmt.Printf("target transfer has check error:%v\r\n", err.Error())
			continue
		}
		if !has {
			_, err = des.Insert(&userResource)
			if err != nil {
				fmt.Printf("target transfer insert error:%v\r\n", err.Error())
				continue
			}
		}
		insertNum++
	}
	fmt.Printf("target and entity transfer success:%v\r\n", insertNum)
}

func getResourceJoin(src, des *xorm.Engine, entityId int, field int) (newResourceJoin models.ResourceJoin, err error) {
	var name string
	switch field {
	case 1:
		oldEntity := EntitySkill{}
		_, err = src.Id(entityId).Get(&oldEntity)
		if err != nil {
			fmt.Printf("entity skill get error:%v\r\n", err.Error())
			return
		}
		name = oldEntity.Name
	case 2:
		oldEntity := EntityAsset{}
		_, err = src.Id(entityId).Get(&oldEntity)
		if err != nil {
			fmt.Printf("entity asset get error:%v\r\n", err.Error())
			return
		}
		name = oldEntity.Name
	case 3:
		oldEntity := EntityWork{}
		_, err = src.Id(entityId).Get(&oldEntity)
		if err != nil {
			fmt.Printf("entity work get error:%v\r\n", err.Error())
			return
		}
		name = oldEntity.Name
	case 4:
		oldEntity := EntityCircle{}
		_, err = src.Id(entityId).Get(&oldEntity)
		if err != nil {
			fmt.Printf("entity circle get error:%v\r\n", err.Error())
			return
		}
		name = oldEntity.Name
	case 5:
		oldEntity := EntityQuest{}
		_, err = src.Id(entityId).Get(&oldEntity)
		if err != nil {
			fmt.Printf("entity quest get error:%v\r\n", err.Error())
			return
		}
		name = oldEntity.Name
	case 6:
		oldEntity := EntityLife{}
		_, err = src.Id(entityId).Get(&oldEntity)
		if err != nil {
			fmt.Printf("entity life get error:%v\r\n", err.Error())
			return
		}
		name = oldEntity.Name
	}

	_, err = des.Table("resource").Join("LEFT", "map_area_resource", "map_area_resource.resource_id = resource.id").Join("LEFT", "area", "area.id = map_area_resource.area_id").Where("area.field_id = ?", field).And("resource.name = ?", name).Get(&newResourceJoin)
	if err != nil {
		return
	}
	if newResourceJoin.Resource.Name == "" {
		fmt.Printf("%v:%v:%v\r\n", field, name, newResourceJoin.Resource.Name)
		errors.New(fmt.Sprintf("new resource name not exist:%v\r\n", newResourceJoin))
		return
	}
	return
}
