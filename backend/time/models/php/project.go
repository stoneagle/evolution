package php

import (
	"time"
)

type Project struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	Text      string    `xorm:"not null default '' comment('内容') VARCHAR(255)"`
	StartDate time.Time `xorm:"not null comment('开始日期') DATETIME"`
	Duration  int       `xorm:"not null default 0 comment('持续时间') INT(11)"`
	Progress  float32   `xorm:"not null default 0 comment('进度') FLOAT"`
	UserId    int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	TargetId  int       `xorm:"not null default 0 comment('所属目标') INT(11)"`
	Del       int       `xorm:"not null default 0 comment('软删除') TINYINT(1)"`
	Ctime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
