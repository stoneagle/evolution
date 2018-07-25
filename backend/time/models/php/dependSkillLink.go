package php

import (
	"time"
)

type DependSkillLink struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	SkillId  int       `xorm:"not null default 0 comment('技能id') INT(11)"`
	DependId string    `xorm:"not null default '0' comment('前置技能') VARCHAR(255)"`
	Ctime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
