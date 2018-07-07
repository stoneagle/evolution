package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Phase struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	Desc                string `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Level               int    `xorm:"unique(name) not null default 0 comment('阶段级别') INT(11)" structs:"level,omitempty"`
	Threshold           int    `xorm:"not null default 0 comment('时间门槛') INT(11)" structs:"threshold,omitempty"`
	FieldId             int    `xorm:"unique(name) not null default 0 comment('隶属方向') INT(11)" structs:"field_id,omitempty"`
}

func (m *Phase) TableName() string {
	return "phase"
}

func (m *Phase) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
