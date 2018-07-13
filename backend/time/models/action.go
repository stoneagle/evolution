package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Action struct {
	es.ModelWithId `xorm:"extends"`
	TaskId         int       `xorm:"not null default 0 comment('所属task') INT(11)" structs:"task_id,omitempty"`
	StartDate      time.Time `xorm:"TIMESTAMP not null comment('开始时间') DATETIME" structs:"-"`
	EndDate        time.Time `xorm:"TIMESTAMP not null comment('结束时间') DATETIME" structs:"-"`
	Name           string    `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	UserId         int       `xorm:"not null default 0 comment('执行人') INT(11)" structs:"user_id,omitempty"`
	Time           int       `xorm:"not null default 0 comment('花费时间:单位分钟') INT(11)" structs:"time,omitempty"`

	TaskIds []int `xorm:"-" structs:"task_id,omitempty"`
	Task    Task  `xorm:"-" structs:"-"`
}

func (m *Action) TableName() string {
	return "action"
}

type ActionJoin struct {
	Action `xorm:"extends" json:"-"`
	Task   `xorm:"extends" json:"-"`
}

func (m *Action) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
