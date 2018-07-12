package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type QuestResource struct {
	es.ModelWithId `xorm:"extends"`
	QuestId        int    `xorm:"unique not null default 0 comment('所属quest') INT(11)" structs:"quest_id,omitempty"`
	ResourceId     int    `xorm:"unique(quest_id) not null default 0 comment('要求资源id') INT(11)" structs:"resource_id,omitempty"`
	Number         int    `xorm:"not null default 0 comment('数量') INT(11)" structs:"number,omitempty"`
	Desc           string `xorm:"not null default '' comment('描述') VARCHAR(255)" structs:"desc,omitempty"`
	Status         int    `xorm:"not null default 1 comment('当前状态:1未匹配2已匹配') INT(11)" structs:"status,omitempty"`
}

var (
	QuestResourceStatusUnmatch = 1
	QuestResourceStatusMatched = 2
)

func (m *QuestResource) TableName() string {
	return "quest_resource"
}

func (m *QuestResource) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
