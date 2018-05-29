package models

import "time"

type GeneralWithDeleted struct {
	Id        int       `xorm:"pk autoincr" form:"id" json:"id"`
	CreatedAt time.Time `xorm:"created comment('创建时间')" form:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" form:"updatedAt" json:"updatedAt"`
	DeletedAt time.Time `xorm:"deleted comment('软删除时间')" form:"deletedAt" json:"deletedAt"`
}

type General struct {
	Id        int       `xorm:"pk autoincr" form:"id" json:"id"`
	CreatedAt time.Time `xorm:"created comment('创建时间')" form:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updated comment('修改时间')" form:"updatedAt" json:"updatedAt"`
}
