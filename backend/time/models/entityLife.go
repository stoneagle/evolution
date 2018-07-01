package models

import (
	"time"
)

type EntityLife struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`
	Name  string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	Desc  string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Ctime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
