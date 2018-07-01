package models

import (
	"time"
)

type EntitySkill struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	Name        string    `xorm:"not null default '' comment('技能名称') VARCHAR(255)"`
	Description string    `xorm:"not null comment('技能描述') TEXT"`
	RankDesc    string    `xorm:"comment('等级描述，用逗号分隔') TEXT"`
	ImgUrl      string    `xorm:"comment('技能图片') TEXT"`
	MaxPoints   int       `xorm:"not null comment('最大级别') SMALLINT(6)"`
	TypeId      int       `xorm:"not null default 0 comment('所处层级') INT(11)"`
	Ctime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
