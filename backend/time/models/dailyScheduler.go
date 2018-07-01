package models

import (
	"time"
)

type DailyScheduler struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	StartDate time.Time `xorm:"not null comment('开始时间') DATETIME"`
	EndDate   time.Time `xorm:"not null comment('结束时间') DATETIME"`
	UserId    int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	DailyId   int       `xorm:"not null default 0 comment('所属模板') INT(11)"`
	Ctime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
