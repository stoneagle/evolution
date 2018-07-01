package models

import "time"

type General struct {
	Ctime time.Time `xorm:"created comment('创建时间')" structs:"omitempty"`
	Utime time.Time `xorm:"updated comment('修改时间')" structs:"omitempty"`
}

type GeneralWithId struct {
	Id      int `xorm:"pk autoincr" structs:"omitempty"`
	General `xorm:"extends"`
}

type GeneralWithDeleted struct {
	GeneralWithId `xorm:"extends"`
	DeletedAt     time.Time `xorm:"deleted comment('软删除时间')" structs:"omitempty"`
}
