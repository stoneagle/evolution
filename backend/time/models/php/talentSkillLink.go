package php

import (
	"time"
)

type TalentSkillLink struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	SkillId  int       `xorm:"not null default 0 comment('所处层级') INT(11)"`
	TalentId string    `xorm:"not null default '0' comment('相关天赋，用逗号分隔') VARCHAR(255)"`
	Ctime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
