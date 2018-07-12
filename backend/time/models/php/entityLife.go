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

	olds := make([]EntityLife, 0)
	src.Find(&olds)
	news := make([]models.Resource, 0)
	for _, one := range olds {
		tmp := models.Resource{}
		tmp.Name = one.Name
		tmp.Desc = one.Desc
		tmp.Year = 0
		tmp.AreaId = area.Id
		tmp.CreatedAt = one.Ctime
		tmp.UpdatedAt = one.Utime
		news = append(news, tmp)
	}
	affected, err := des.Insert(&news)
	if err != nil {
		fmt.Printf("entity life transfer error:%v\r\n", err.Error())
	} else {
		fmt.Printf("entity life transfer success:%v\r\n", affected)
	}
}
