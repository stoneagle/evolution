package models

type Pool struct {
	GeneralWithDeleted `xorm:"extends"`
	AssetType          `xorm:"extends"`
	Strategy           string `xorm:"varchar(128) notnull comment('策略')" form:"Strategy" json:"Strategy"`
	Name               string `xorm:"varchar(128) notnull comment('名称')" form:"Name" json:"Name"`
	Status             int    `xorm:"smallint(4) notnull comment('状态')" form:"Status" json:"Status"`
	Item               []Item `xorm:"-"`
}
