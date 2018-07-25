package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type EntityAsset struct {
	Id     int       `xorm:"not null pk autoincr INT(11)"`
	Name   string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	Desc   string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Year   int       `xorm:"not null default 0 comment('年份') INT(11)"`
	AreaId int       `xorm:"not null default 0 comment('隶属领域') INT(11)"`
	Status int       `xorm:"not null default 0 comment('状态，公共、商业、私人') INT(11)"`
	Ctime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *EntityAsset) Transfer(src, des *xorm.Engine) {
	session := des.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Printf("asset session error:%v\r\n", err.Error())
		return
	}

	olds := make([]EntityAsset, 0)
	src.Find(&olds)
	insertNum := 0
	for _, one := range olds {
		resource := models.Resource{}
		resource.Name = one.Name
		resource.Desc = one.Desc
		resource.Year = one.Year
		resource.CreatedAt = one.Ctime
		resource.UpdatedAt = one.Utime
		_, err = session.Insert(&resource)
		if err != nil {
			session.Rollback()
			fmt.Printf("asset resource insert error:%v\r\n", err.Error())
			return
		}
		mapAreaResource := models.MapAreaResource{}
		mapAreaResource.AreaId = one.AreaId
		mapAreaResource.ResourceId = resource.Id
		_, err = session.Insert(&mapAreaResource)
		if err != nil {
			session.Rollback()
			fmt.Printf("asset resource map insert error:%v\r\n", err.Error())
			return
		}
		insertNum++
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("asset session commit error:%v\r\n", err)
	}
	fmt.Printf("asset transfer success:%v\r\n", insertNum)
}
