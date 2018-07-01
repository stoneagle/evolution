package models

import (
	"time"
)

type GanttLinks struct {
	Id     int       `xorm:"not null pk autoincr INT(11)"`
	Source int       `xorm:"not null default 0 comment('来源') INT(11)"`
	Target int       `xorm:"not null default 0 comment('目标') INT(11)"`
	Type   string    `xorm:"not null default '0' comment('类别') VARCHAR(1)"`
	UserId int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	Del    int       `xorm:"not null default 0 comment('乱删除') TINYINT(1)"`
	Ctime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
