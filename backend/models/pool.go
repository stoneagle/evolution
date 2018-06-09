package models

type Pool struct {
	GeneralWithDeleted `xorm:"extends"`
	Name               string `xorm:"varchar(128) notnull comment('名称')" form:"Name" binding:"required" json:"Name"`
	Strategy           string `xorm:"varchar(128) notnull comment('策略')" form:"Strategy" binding:"required" json:"Strategy"`
	Status             int    `xorm:"smallint(4) notnull comment('状态')" form:"Status" json:"Status"`
	Type               int    `xorm:"smallint(4) notnull comment('类别')" form:"Type" binding:"required" json:"Type"`
}
