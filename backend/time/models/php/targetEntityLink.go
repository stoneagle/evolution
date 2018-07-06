package php

import (
	"time"
)

type TargetEntityLink struct {
	Id       int       `xorm:"not null pk autoincr INT(11)"`
	TargetId int       `xorm:"not null default 0 comment('领域对象ID') INT(11)"`
	EntityId int       `xorm:"not null default 0 comment('相关实体ID') INT(11)"`
	Ctime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
