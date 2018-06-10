package models

import (
	"quant/backend/rpc/engine"
)

type Classify struct {
	General  `xorm:"extends"`
	Resource engine.QuantType `xorm:"smallint(4) notnull comment('所属资源')" form:"Resource" json:"Resource"`
	Type     string           `xorm:"varchar(128) notnull comment('类别')" form:"Type" json:"Type"`
	SubType  string           `xorm:"varchar(128) comment('子类别')" form:"SubType" json:"SubType"`
	Name     string           `xorm:"varchar(128) notnull comment('名称')" form:"Name" json:"Name"`
	Tag      string           `xorm:"varchar(128) notnull comment('标签')" form:"Tag" json:"Tag"`
}
