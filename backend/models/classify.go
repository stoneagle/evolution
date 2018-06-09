package models

type Classify struct {
	General `xorm:"extends"`
	Name    string `xorm:"varchar(32) notnull comment('名称')" form:"Name" json:"Name"`
	Type    int    `xorm:"smallint(4) notnull comment('类别')" form:"Type" json:"Type"`
}
