package models

import "time"

type GeneralWithDeleted struct {
	Id        int       `xorm:"pk autoincr" form:"Id" json:"Id"`
	CreatedAt time.Time `xorm:"created comment('创建时间')" form:"CreatedAt" json:"CreatedAt"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" form:"UpdatedAt" json:"UpdatedAt"`
	DeletedAt time.Time `xorm:"deleted comment('软删除时间')" form:"DeletedAt" json:"DeletedAt"`
}

type General struct {
	Id        int       `xorm:"pk autoincr" form:"Id" json:"Id"`
	CreatedAt time.Time `xorm:"created comment('创建时间')" form:"CreatedAt" json:"CreatedAt"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" form:"UpdatedAt" json:"UpdatedAt"`
}
