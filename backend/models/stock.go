package models

type Stock struct {
	GeneralWithDeleted `xorm:"extends"`
	Code               string `xorm:"varchar(16) notnull unique comment('编号')" form:"code" json:"code"`
	Name               string `xorm:"varchar(32) notnull comment('名称')" form:"name" json:"name"`
	Status             int    `xorm:"smallint(4) notnull comment('状态')" form:"status" json:"status"`
	Source             int    `xorm:"smallint(4) notnull comment('来源')" form:"source" json:"source"`
}
