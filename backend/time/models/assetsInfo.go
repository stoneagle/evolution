package models

import (
	"time"
)

type AssetsInfo struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	ObjId      int       `xorm:"not null default 0 comment('领域对象id') SMALLINT(6)"`
	TradeNum   int       `xorm:"not null default 0 comment('交易次数') INT(11)"`
	HeadCount  int       `xorm:"not null default 0 comment('人力资源数量') INT(11)"`
	TimeSpan   int       `xorm:"not null default 0 comment('日期间隔/天') INT(11)"`
	Value      int       `xorm:"not null default 0 comment('当前价值') INT(11)"`
	IncomeFlow int       `xorm:"not null default 0 comment('现金流/元') INT(11)"`
	Position   string    `xorm:"not null default '' comment('展示位置') VARCHAR(255)"`
	Ctime      time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime      time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}
