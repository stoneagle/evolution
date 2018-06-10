package models

type Item struct {
	GeneralWithDeleted `xorm:"extends"`
	Code               string `xorm:"varchar(16) notnull unique comment('编号')" form:"Code" json:"Code"`
	Name               string `xorm:"varchar(32) notnull comment('名称')" form:"Name" json:"Name"`
	Status             int    `xorm:"smallint(4) notnull comment('状态')" form:"Status" json:"Status"`
}
