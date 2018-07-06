package models

import (
	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Phase struct {
	GeneralWithDeleted `xorm:"extends"`
	Name               string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	Desc               string `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Level              int    `xorm:"unique(name) not null default 0 comment('阶段级别') INT(11)" structs:"level,omitempty"`
	FieldId            int    `xorm:"unique(name) not null default 0 comment('隶属方向') INT(11)" structs:"field_id,omitempty"`
}

func (m *Phase) TableName() string {
	return "phase"
}

func (m *Phase) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.General.BuildCondition(params, keyPrefix)
	return condition
}
