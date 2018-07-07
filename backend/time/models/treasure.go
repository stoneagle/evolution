package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Treasure struct {
	es.ModelWithId `xorm:"extends"`
	UserId         int `xorm:"unique not null default 0 comment('隶属用户') INT(11)" structs:"user_id,omitempty"`
	EntityId       int `xorm:"unique(user_id) not null default 0 comment('隶属实体') INT(11)" structs:"entity_id,omitempty"`
	PhaseId        int `xorm:"not null default 0 comment('所处阶段') INT(11)" structs:"phase_id,omitempty"`
	Status         int `xorm:"not null default 0 comment('状态:1已收录2开拓中') INT(11)" structs:"status,omitempty"`
	SumTime        int `xorm:"not null default 0 comment('总投入时间') INT(11)" structs:"sum_time,omitempty"`
}

var (
	TreasureStatusCollect = 1
	TreasureStatusExec    = 2
)

func (m *Treasure) TableName() string {
	return "treasure"
}

func (m *Treasure) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
