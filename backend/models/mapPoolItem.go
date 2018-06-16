package models

type MapPoolItem struct {
	General `xorm:"extends"`
	PoolId  int `xorm:"not null comment('投资池id')"`
	ItemId  int `xorm:"not null comment('标的id')"`
}
