package php

import (
	"time"
)

type UserSkillLink struct {
	Id      int       `xorm:"not null pk autoincr INT(11)"`
	SkillId int       `xorm:"not null default 0 comment('技能ID') INT(11)"`
	UserId  int       `xorm:"not null default 0 comment('用户ID') INT(11)"`
	Level   int       `xorm:"not null default 0 comment('技能级别') SMALLINT(6)"`
	Ctime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
