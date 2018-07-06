package php

import (
	"time"
)

type PlanProject struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	PlanId    int       `xorm:"not null default 0 comment('所属计划') INT(11)"`
	ProjectId int       `xorm:"not null default 0 comment('目标项目') INT(11)"`
	Hours     int       `xorm:"not null default 0 comment('投资时间') INT(11)"`
	UserId    int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	Ctime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
