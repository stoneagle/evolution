package models

import (
	"time"
)

type Area struct {
	Id      int       `xorm:"not null pk autoincr INT(11)"`
	Name    string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
	Parent  int       `xorm:"not null default 0 comment('父id') INT(11)"`
	Del     int       `xorm:"not null default 0 comment('软删除') TINYINT(1)"`
	Level   int       `xorm:"not null default 0 comment('所属层级') SMALLINT(6)"`
	FieldId int       `xorm:"not null default 0 comment('所属范畴id') SMALLINT(6)"`
	Ctime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

type AreaTree struct {
	Value    string     `json:"value"`
	Children []AreaNode `json:"children"`
}

type AreaNode struct {
	Id       int        `json:"id"`
	Value    string     `json:"value"`
	Children []AreaNode `json:"children"`
}

const (
	AreaFieldSkill  int = 1
	AreaFieldAsset  int = 2
	AreaFieldWork   int = 3
	AreaFieldCircle int = 4
	AreaFieldQuest  int = 5
)

var (
	AreaFiledMap map[int]string = map[int]string{
		AreaFieldSkill:  "知识",
		AreaFieldAsset:  "财富",
		AreaFieldWork:   "文化",
		AreaFieldCircle: "社交",
		AreaFieldQuest:  "挑战",
	}
)
