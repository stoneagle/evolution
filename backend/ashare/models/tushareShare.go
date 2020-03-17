package models

type TushareShare struct {
	Symbol string `xorm:"not null default '' comment('代码') TEXT" structs:"symbol,omitempty"`
	Name   string `xorm:"not null default '' comment('名称') TEXT" structs:"code,omitempty"`
	Market string `xorm:"not null default '' comment('板块') TEXT" structs:"market,omitempty"`
}

func NewTushareShare() *TushareShare {
	ret := TushareShare{}
	return &ret
}

func (m *TushareShare) TableName() string {
	return "tushare_share"
}
