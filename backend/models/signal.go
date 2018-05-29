package models

type Signal struct {
	General       `xorm:"extends"`
	MapPoolItemId int     `xorm:"not null comment('关联id')"`
	Time          string  `xorm:"notnull comment('时间')" form:"time" json:"time"`
	Type          int     `xorm:"smallint(4) notnull comment('类别')" form:"type" json:"type"`
	Bid           float64 `xorm:"notnull comment('叫价')" form:"bid" json:"bid"`
}
