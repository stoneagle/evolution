package php

import (
	"time"
)

type Plan struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	FromDate string    `xorm:"not null default '' comment('开始日期') VARCHAR(255)"`
	ToDate   string    `xorm:"not null default '' comment('结束日期') VARCHAR(255)"`
	DailyId  int       `xorm:"not null default 0 comment('作息模板') INT(11)"`
	UserId   int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	Ctime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
