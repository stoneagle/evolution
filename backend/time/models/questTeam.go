package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type QuestTeam struct {
	es.ModelWithId `xorm:"extends"`
	QuestId        int       `xorm:"unique not null default 0 comment('所属quest') INT(11)" structs:"quest_id,omitempty"`
	StartDate      time.Time `xorm:"not null comment('开始日期') DATETIME" structs:"start_date,omitempty"`
	EndDate        time.Time `xorm:"not null comment('开始日期') DATETIME" structs:"start_date,omitempty"`
	UserId         int       `xorm:"unique(quest_id) not null default 0 comment('成员id') INT(11)" structs:"founder_id,omitempty"`
}

func (m *QuestTeam) TableName() string {
	return "quest_team"
}

func (m *QuestTeam) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
