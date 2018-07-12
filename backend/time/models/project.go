package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Project struct {
	es.ModelWithDeleted `xorm:"extends"`
	QuestId             int       `xorm:"not null default 0 comment('所属quest') INT(11)" structs:"quest_id,omitempty"`
	AreaId              int       `xorm:"not null default 0 comment('所属area') INT(11)" structs:"area_id,omitempty"`
	Name                string    `xorm:"not null unique default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	StartDate           time.Time `xorm:"not null comment('开始日期') DATETIME"`
	Duration            int       `xorm:"not null default 0 comment('持续时间') INT(11)" structs:"duration,omitempty"`
	QuestIds            []int     `xorm:"-" structs:"quest_id,omitempty"`
}

func (m *Project) TableName() string {
	return "project"
}

func (m *Project) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
