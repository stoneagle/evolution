package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Country struct {
	es.ModelWithId `xorm:"extends"`
	Name           string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	EnName         string `xorm:"not null default '' comment('英文名称') VARCHAR(255)"`
}

func (m *Country) TableName() string {
	return "country"
}

func (m *Country) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
