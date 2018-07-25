package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type Country struct {
	Id     int       `xorm:"not null pk autoincr INT(11)"`
	Name   string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	EnName string    `xorm:"not null default '' comment('英文名称') VARCHAR(255)"`
	Ctime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *Country) Transfer(src, des *xorm.Engine) {
	olds := make([]Country, 0)
	src.Find(&olds)
	news := make([]models.Country, 0)
	for _, one := range olds {
		tmp := models.Country{}
		tmp.Id = one.Id
		tmp.Name = one.Name
		tmp.EnName = one.EnName
		tmp.CreatedAt = one.Ctime
		tmp.UpdatedAt = one.Utime
		news = append(news, tmp)
	}
	affected, err := des.Insert(&news)
	if err != nil {
		fmt.Printf("country transfer error:%v\r\n", err.Error())
	} else {
		fmt.Printf("country transfer success:%v\r\n", affected)
	}
}
