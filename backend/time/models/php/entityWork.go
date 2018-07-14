package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type EntityWork struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	Name      string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	Desc      string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Year      int       `xorm:"not null default 0 comment('年代') INT(11)"`
	CountryId int       `xorm:"not null default 0 comment('所属国家') INT(11)"`
	AreaId    int       `xorm:"not null default 0 comment('所属领域') INT(11)"`
	Ctime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *EntityWork) Transfer(src, des *xorm.Engine) {
	session := des.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Printf("work session error:%v\r\n", err.Error())
		return
	}

	olds := make([]EntityWork, 0)
	src.Find(&olds)
	insertNum := 0
	areaMap := map[string][]int{}
	resourceMap := map[string]models.Resource{}
	for _, one := range olds {
		resource := models.Resource{}
		resource.Name = one.Name
		resource.Desc = one.Desc
		resource.Year = one.Year
		resource.CreatedAt = one.Ctime
		resource.UpdatedAt = one.Utime

		areaSlice, ok := areaMap[resource.Name]
		if !ok {
			areaMap[resource.Name] = make([]int, 0)
		}
		areaSlice = append(areaSlice, one.AreaId)
		areaMap[resource.Name] = areaSlice

		if _, ok := resourceMap[resource.Name]; !ok {
			resourceMap[resource.Name] = resource
		}
	}
	for resourceName, areaSlice := range areaMap {
		resource, ok := resourceMap[resourceName]
		if !ok {
			session.Rollback()
			fmt.Printf("work resource can not find in map\r\n")
			return
		}
		_, err = session.Insert(&resource)
		if err != nil {
			session.Rollback()
			fmt.Printf("work resource insert error:%v\r\n", err.Error())
			return
		}
		for _, areaId := range areaSlice {
			mapAreaResource := models.MapAreaResource{}
			mapAreaResource.AreaId = areaId
			mapAreaResource.ResourceId = resource.Id
			_, err = session.Insert(&mapAreaResource)
			if err != nil {
				session.Rollback()
				fmt.Printf("work resource map insert error:%v\r\n", err.Error())
				return
			}
		}
		insertNum++
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("resource session commit error:%v\r\n", err)
	}
	fmt.Printf("resource transfer success:%v\r\n", insertNum)
}
