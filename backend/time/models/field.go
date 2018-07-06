package models

type Field struct {
	GeneralWithDeleted `xorm:"extends"`
	Name               string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	Desc               string `xorm:"not null default '' comment('描述') VARCHAR(255)"`
}

func (m *Field) TableName() string {
	return "field"
}
