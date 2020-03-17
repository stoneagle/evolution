package models

type TushareLimitDaily struct {
	Symbol string `xorm:"not null default '' comment('代码') TEXT" structs:"symbol,omitempty"`
	Name   string `xorm:"not null default '' comment('名称') TEXT" structs:"code,omitempty"`
	Market string `xorm:"not null default '' comment('板块') TEXT" structs:"market,omitempty"`
}

func NewTushareLimitDaily() *TushareLimitDaily {
	ret := TushareLimitDaily{}
	return &ret
}

func (m *TushareLimitDaily) TableName() string {
	return "tushare_limit_daily"
}
