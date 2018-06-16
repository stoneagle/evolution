package models

type Classify struct {
	GeneralWithDeleted `xorm:"extends"`
	AssetType          `xorm:"extends"`
	Source             `xorm:"extends"`
	Name               string `xorm:"varchar(128) notnull comment('名称')" form:"Name" json:"Name"`
	Tag                string `xorm:"varchar(128) notnull comment('标签')" form:"Tag" json:"Tag"`
}
