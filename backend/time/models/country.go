package models

import (
	es "evolution/backend/common/structs"
)

type Country struct {
	es.ModelWithId `xorm:"extends"`
	Name           string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	EnName         string `xorm:"not null default '' comment('英文名称') VARCHAR(255)"`
}

func (m *Country) TableName() string {
	return "country"
}
