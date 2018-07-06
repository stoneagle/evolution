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
	olds := make([]EntityLife, 0)
	src.Find(&olds)
	news := make([]models.Entity, 0)
	for _, one := range olds {
		tmp := models.Entity{}
		tmp.Name = one.Name
		tmp.Desc = one.Desc
		tmp.Year = 0
		tmp.AreaId = 0
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
