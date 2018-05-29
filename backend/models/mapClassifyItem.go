package models

type MapClassifyItem struct {
	General    `xorm:"extends"`
	ClassifyId int `xorm:"not null comment('分类id')"`
	ItemId     int `xorm:"not null comment('投资商品id')"`
}
