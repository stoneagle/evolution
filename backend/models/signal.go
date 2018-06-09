package models

type Signal struct {
	General       `xorm:"extends"`
	MapPoolItemId int     `xorm:"not null comment('关联id')"`
	Time          string  `xorm:"notnull comment('时间')" form:"Time" json:"Time"`
	Type          int     `xorm:"smallint(4) notnull comment('类别')" form:"Type" json:"Type"`
	Bid           float64 `xorm:"notnull comment('叫价')" form:"Bid" json:"Bid"`
}
