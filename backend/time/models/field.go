package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Field struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	Desc                string `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Color               string `xorm:"not null default '' comment('颜色') VARCHAR(255)"`
}

func (m *Field) TableName() string {
	return "field"
}

func (m *Field) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
