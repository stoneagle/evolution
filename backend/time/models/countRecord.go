package models

import (
	"time"
)

type CountRecord struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	ActionId int       `xorm:"not null default 0 comment('所属行动') INT(11)"`
	Status   int       `xorm:"not null default 0 comment('状态') SMALLINT(6)"`
	UserId   int       `xorm:"not null comment('所属用户') INT(11)"`
	InitTime int       `xorm:"not null default 0 comment('起始时间') INT(11)"`
	Ctime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
