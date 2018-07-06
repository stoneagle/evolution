package models

import (
	"time"
)

type Task struct {
	GeneralWithDeleted `xorm:"extends"`
	Text               string    `xorm:"not null default '' comment('内容') VARCHAR(255)"`
	StartDate          time.Time `xorm:"not null comment('开始日期') DATETIME"`
	Progress           float32   `xorm:"not null default 0 comment('进度') FLOAT"`
	Parent             int       `xorm:"not null default 0 comment('父id') INT(11)"`
	Duration           int       `xorm:"not null comment('持续时间') INT(11)"`
	UserId             int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	EntityId           int       `xorm:"not null default 0 comment('相关实体') INT(11)"`
}
