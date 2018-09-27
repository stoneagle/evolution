package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type EntityCircle struct {
	Id     int       `xorm:"not null pk autoincr INT(11)"`
	Name   string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	Desc   string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	AreaId int       `xorm:"not null default 0 comment('隶属领域') INT(11)"`
	Ctime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *EntityCircle) Transfer(src, des *xorm.Engine) {
	session := des.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Printf("circle session error:%v\r\n", err.Error())
		return
	}

	olds := make([]EntityCircle, 0)
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
			fmt.Printf("circle resource insert error:%v\r\n", err.Error())
			return
		}
		mapAreaResource := models.MapAreaResource{}
		mapAreaResource.AreaId = one.AreaId
		mapAreaResource.ResourceId = resource.Id
		_, err = session.Insert(&mapAreaResource)
		if err != nil {
			session.Rollback()
			fmt.Printf("circle resource map insert error:%v\r\n", err.Error())
			return
		}
		insertNum++
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("circle session commit error:%v\r\n", err)
	}
	fmt.Printf("circle transfer success:%v\r\n", insertNum)
}
