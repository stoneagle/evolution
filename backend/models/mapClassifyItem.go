package models

type MapClassifyItem struct {
	General    `xorm:"extends"`
	ClassifyId int `xorm:"not null comment('分类id')"`
	ItemId     int `xorm:"not null comment('标的id')"`
}

type MapClassifyItemJoin struct {
	MapClassifyItem
	Classify `xorm:"extends" json:"-"`
	Item     `xorm:"extends" json:"-"`
}
