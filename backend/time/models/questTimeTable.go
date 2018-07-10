package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type QuestTimeTable struct {
	es.ModelWithId `xorm:"extends"`
	QuestId        int       `xorm:"not null default 0 comment('所属quest') INT(11)" structs:"quest_id,omitempty"`
	StartTime      time.Time `xorm:"not null comment('开始时间') DATETIME" structs:"start_time,omitempty"`
	EndTime        time.Time `xorm:"not null comment('结束时间') DATETIME" structs:"end_time,omitempty"`
	Type           int       `xorm:"not null default 1 comment('类别:1工作日2节假日') INT(11)" structs:"status,omitempty"`
}

var (
	QuestTimeTableTypeWorkday = 1
	QuestTimeTableTypeHoliday = 2
)

func (m *QuestTimeTable) TableName() string {
	return "quest_time_table"
}

func (m *QuestTimeTable) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
