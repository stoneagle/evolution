package models

import (
	"time"
)

type AreaSkillLink struct {
	Id      int       `xorm:"not null pk autoincr INT(11)"`
	SkillId int       `xorm:"not null default 0 comment('技能id') INT(11)"`
	AreaId  int       `xorm:"not null default 0 comment('领域id') INT(11)"`
	Ctime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
