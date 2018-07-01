package models

import "time"

type GeneralWithDeleted struct {
	Id        int       `xorm:"pk autoincr" structs:"omitempty"`
	CreatedAt time.Time `xorm:"created comment('创建时间')" structs:"omitempty"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" structs:"omitempty"`
	DeletedAt time.Time `xorm:"deleted comment('软删除时间')" structs:"omitempty"`
}

type General struct {
	Id        int       `xorm:"pk autoincr" structs:"omitempty"`
	CreatedAt time.Time `xorm:"created comment('创建时间')" structs:"omitempty"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" structs:"omitempty"`
}

type GeneralWithoutId struct {
	CreatedAt time.Time `xorm:"created comment('创建时间')" structs:"omitempty"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" structs:"omitempty"`
}
