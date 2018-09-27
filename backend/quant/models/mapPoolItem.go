package models

type MapPoolItem struct {
	GeneralWithoutId `xorm:"extends"`
	PoolId           int `xorm:"not null comment('投资池id')"`
	ItemId           int `xorm:"not null comment('标的id')"`
}

type MapPoolItemJoin struct {
	MapPoolItem
	Pool `xorm:"extends" json:"-"`
	Item `xorm:"extends" json:"-"`
}
