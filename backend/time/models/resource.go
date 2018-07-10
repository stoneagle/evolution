package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Resource struct {
	es.ModelWithId `xorm:"extends"`
	UserId         int `xorm:"unique not null default 0 comment('隶属用户') INT(11)" structs:"user_id,omitempty"`
	EntityId       int `xorm:"unique(user_id) not null default 0 comment('隶属实体') INT(11)" structs:"entity_id,omitempty"`
	PhaseId        int `xorm:"not null default 0 comment('所处阶段') INT(11)" structs:"phase_id,omitempty"`
	Status         int `xorm:"not null default 0 comment('状态:1已收录2开拓中') INT(11)" structs:"status,omitempty"`
	SumTime        int `xorm:"not null default 0 comment('总投入时间') INT(11)" structs:"sum_time,omitempty"`
}

var (
	ResourceStatusCollect = 1
	ResourceStatusExec    = 2
)

func (m *Resource) TableName() string {
	return "resource"
}

func (m *Resource) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
