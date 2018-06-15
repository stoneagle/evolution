package models

import (
	"quant/backend/rpc/engine"
)

type Classify struct {
	GeneralWithDeleted `xorm:"extends"`
	Asset              engine.AssetType `xorm:"smallint(4) notnull comment('所属资产')" form:"Resource" json:"Asset"`
	AssetString        string           `xorm:"-" form:"AssetString" json:"AssetString"`
	Type               string           `xorm:"varchar(128) notnull comment('类别')" form:"Type" json:"Type"`
	Source             string           `xorm:"varchar(128) notnull comment('来源')" form:"Source" json:"Source"`
	Sub                string           `xorm:"varchar(128) comment('子类别')" form:"Sub" json:"Sub"`
	Name               string           `xorm:"varchar(128) notnull comment('名称')" form:"Name" json:"Name"`
	Tag                string           `xorm:"varchar(128) notnull comment('标签')" form:"Tag" json:"Tag"`
}
