package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Project struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string    `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	QuestTargetId       int       `xorm:"not null default 0 comment('所属target') INT(11)" structs:"quest_target_id,omitempty"`
	StartDate           time.Time `xorm:"not null comment('开始日期') DATETIME"`
	Duration            int       `xorm:"not null default 0 comment('持续时间') INT(11)" structs:"duration,omitempty"`

	QuestTarget QuestTarget `xorm:"-" structs:"-"`
	Area        Area        `xorm:"-" structs:"-"`
	Quest       Quest       `xorm:"-" structs:"-"`
}

func (m *Project) TableName() string {
	return "project"
}

type ProjectJoin struct {
	Project     `xorm:"extends" json:"-"`
	QuestTarget `xorm:"extends" json:"-"`
	Area        `xorm:"extends" json:"-"`
}

func (m *Project) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)

	questTargetCondition := m.QuestTarget.BuildCondition()
	if len(questTargetCondition) > 0 {
		for k, v := range questTargetCondition {
			condition[k] = v
		}
	}

	return condition
}
