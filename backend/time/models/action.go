package models

import (
	"time"
)

type Action struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	StartDate time.Time `xorm:"not null comment('开始时间') DATETIME"`
	EndDate   time.Time `xorm:"not null comment('结束时间') DATETIME"`
	Duration  int       `xorm:"not null default 0 comment('持续时间') SMALLINT(6)"`
	Text      string    `xorm:"not null default '' comment('内容') VARCHAR(255)"`
	UserId    int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	TaskId    int       `xorm:"not null default 0 comment('所属任务') INT(11)"`
	ExecTime  int       `xorm:"not null default 0 comment('实际时间') INT(6)"`
	PlanTime  int       `xorm:"not null comment('计划时间') SMALLINT(6)"`
	Status    int       `xorm:"not null default 0 comment('状态') SMALLINT(6)"`
	Desc      string    `xorm:"comment('描述') TEXT"`
	Ctime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
