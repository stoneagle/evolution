package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type EntityQuest struct {
	Id     int       `xorm:"not null pk autoincr INT(11)"`
	Name   string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	Desc   string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	AreaId int       `xorm:"not null default 0 comment('隶属领域') INT(11)"`
	Ctime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *EntityQuest) Transfer(src, des *xorm.Engine) {
	session := des.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Printf("quest session error:%v\r\n", err.Error())
		return
	}

	olds := make([]EntityQuest, 0)
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
			fmt.Printf("quest resource insert error:%v\r\n", err.Error())
			return
		}
		mapAreaResource := models.MapAreaResource{}
		mapAreaResource.AreaId = one.AreaId
		mapAreaResource.ResourceId = resource.Id
		_, err = session.Insert(&mapAreaResource)
		if err != nil {
			session.Rollback()
			fmt.Printf("quest resource map insert error:%v\r\n", err.Error())
			return
		}
		insertNum++
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("quest session commit error:%v\r\n", err)
	}
	fmt.Printf("quest transfer success:%v\r\n", insertNum)
}
