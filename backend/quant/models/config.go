package models

import "evolution/backend/quant/rpc/engine"

type AssetType struct {
	Asset       engine.AssetType `xorm:"smallint(4) notnull comment('所属资产')" form:"Resource" json:"Asset" structs:"asset,omitempty"`
	AssetString string           `xorm:"-" form:"AssetString" json:"AssetString" structs:"-"`
	Type        string           `xorm:"varchar(128) notnull comment('类别')" form:"Type" json:"Type" structs:"type,omitempty"`
}

type Source struct {
	Main string `xorm:"varchar(128) notnull comment('主来源')" form:"Main" json:"Main" structs:"main,omitempty"`
	Sub  string `xorm:"varchar(128) comment('子类别')" form:"Sub" json:"Sub" structs:"sub,omitempty"`
}
