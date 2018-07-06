package php

import (
	"time"
)

type Target struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	Name       string    `xorm:"not null default '' comment('名称') VARCHAR(256)"`
	Desc       string    `xorm:"not null default '' comment('描述') VARCHAR(256)"`
	UserId     int       `xorm:"not null default 0 comment('用户ID') INT(11)"`
	PriorityId int       `xorm:"not null default 0 comment('优先级') SMALLINT(6)"`
	FieldId    int       `xorm:"not null default 0 comment('所属范畴') SMALLINT(6)"`
	Ctime      time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime      time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
