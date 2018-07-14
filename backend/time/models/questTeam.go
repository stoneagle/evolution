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
	EndDate        time.Time `xorm:"not null comment('开始日期') DATETIME" structs:"end_date,omitempty"`
	UserId         int       `xorm:"unique(quest_id) not null default 0 comment('成员id') INT(11)" structs:"user_id,omitempty"`
	Quest          Quest     `xorm:"-"`
}

func (m *QuestTeam) TableName() string {
	return "quest_team"
}

type QuestTeamJoin struct {
	QuestTeam `xorm:"extends" json:"-"`
	Quest     `xorm:"extends" json:"-"`
}

func (m *QuestTeam) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)

	questCondition := m.Quest.BuildCondition()
	if len(questCondition) > 0 {
		for k, v := range questCondition {
			condition[k] = v
		}
	}

	return condition
}
