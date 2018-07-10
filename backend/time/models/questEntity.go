package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type QuestEntity struct {
	es.ModelWithId `xorm:"extends"`
	QuestId        int    `xorm:"not null default 0 comment('所属quest') INT(11)" structs:"quest_id,omitempty"`
	EntityId       int    `xorm:"not null default 0 comment('要求资源id') INT(11)" structs:"entity_id,omitempty"`
	PhaseId        int    `xorm:"not null default 0 comment('所处阶段') INT(11)" structs:"phase_id,omitempty"`
	Number         int    `xorm:"not null default 0 comment('数量') INT(11)" structs:"number,omitempty"`
	Desc           string `xorm:"not null default '' comment('描述') VARCHAR(255)" structs:"desc,omitempty"`
	Status         int    `xorm:"not null default 1 comment('当前状态:1未匹配2已匹配') INT(11)" structs:"status,omitempty"`
}

var (
	QuestEntityStatusUnmatch = 1
	QuestEntityStatusMatched = 1
)

func (m *QuestEntity) TableName() string {
	return "quest_entity"
}

func (m *QuestEntity) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
