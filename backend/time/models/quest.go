package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Quest struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string    `xorm:"not null default '' comment('内容') VARCHAR(255)" structs:"name,omitempty"`
	StartDate           time.Time `xorm:"not null comment('开始日期') DATETIME" structs:"start_date,omitempty"`
	ResourceId          int       `xorm:"not null default 0 comment('所属财富') INT(11)" structs:"treasure_id,omitempty"`
	Constraint          int       `xorm:"not null default 0 comment('限制:1重要紧急2重要不紧急3紧急不重要4不重要不紧急') INT(11)" structs:"constraint,omitempty"`
	Status              int       `xorm:"not null default 0 comment('当前状态:1执行中2已完成') INT(11)" structs:"status,omitempty"`
}

func (m *Quest) TableName() string {
	return "quest"
}

func (m *Quest) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
