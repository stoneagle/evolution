package models

type Pool struct {
	GeneralWithDeleted `xorm:"extends"`
	Name               string `xorm:"varchar(128) notnull comment('名称')" form:"name" json:"name"`
	Strategy           string `xorm:"varchar(128) notnull comment('策略')" form:"strategy" json:"strategy"`
	Status             int    `xorm:"smallint(4) notnull comment('状态')" form:"status" json:"status"`
	Type               int    `xorm:"smallint(4) notnull comment('类别')" form:"type" json:"type"`
}
