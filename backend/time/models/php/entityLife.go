package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type EntityLife struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	Name  string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	Desc  string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Ctime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *EntityLife) Transfer(src, des *xorm.Engine) {
	// 插入一条日常，作为life的归属节点
	area := models.Area{}
	area.FieldId = 6
	area.Parent = 0
	area.Type = models.AreaTypeRoot
	area.Name = "生活杂务"
	_, err := des.Insert(&area)
	if err != nil {
		fmt.Printf("life area insert error:%v\r\n", err.Error())
		return
	}

	session := des.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		fmt.Printf("life session error:%v\r\n", err.Error())
		return
	}

	olds := make([]EntityLife, 0)
	src.Find(&olds)
	insertNum := 0
	for _, one := range olds {
		resource := models.Resource{}
		resource.Name = one.Name
		resource.Desc = one.Desc
		resource.Year = 0
		resource.CreatedAt = one.Ctime
		resource.UpdatedAt = one.Utime
		_, err = session.Insert(&resource)
		if err != nil {
			session.Rollback()
			fmt.Printf("life resource insert error:%v\r\n", err.Error())
			return
		}
		mapAreaResource := models.MapAreaResource{}
		mapAreaResource.AreaId = area.Id
		mapAreaResource.ResourceId = resource.Id
		_, err = session.Insert(&mapAreaResource)
		if err != nil {
			session.Rollback()
			fmt.Printf("life resource map insert error:%v\r\n", err.Error())
			return
		}
		insertNum++
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("life session commit error:%v\r\n", err)
	}
	fmt.Printf("life transfer success:%v\r\n", insertNum)
}
