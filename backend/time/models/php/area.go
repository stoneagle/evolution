package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type Area struct {
	Id      int       `xorm:"not null pk autoincr INT(11)"`
	Name    string    `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Parent  int       `xorm:"not null default 0 comment('父id') INT(11)" structs:"parent,omitempty"`
	Del     int       `xorm:"not null default 0 comment('软删除') TINYINT(1)" structs:"del,omitempty"`
	Level   int       `xorm:"not null default 0 comment('所属层级') SMALLINT(6)" structs:"level,omitempty"`
	FieldId int       `xorm:"not null default 0 comment('所属范畴id') SMALLINT(6)" structs:"field_id,omitempty"`
	Ctime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *Area) Transfer(src, des *xorm.Engine) {
	olds := make([]Area, 0)
	src.Find(&olds)
	news := make([]models.Area, 0)
	for _, one := range olds {
		if one.Del == 0 {
			tmp := models.Area{}
			tmp.Id = one.Id
			tmp.Name = one.Name
			tmp.Parent = one.Parent
			tmp.FieldId = one.FieldId
			tmp.CreatedAt = one.Ctime
			tmp.UpdatedAt = one.Utime
			news = append(news, tmp)
		}
	}
	affected, err := des.Insert(&news)
	if err != nil {
		fmt.Printf("area transfer error:%v\r\n", err.Error())
	} else {
		fmt.Printf("area transfer success:%v\r\n", affected)
	}
}
