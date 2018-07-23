package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Phase struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"not null unique default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Desc                string `xorm:"not null default '' comment('描述') VARCHAR(255)" structs:"desc,omitempty"`
	Level               int    `xorm:"unique(name) not null default 0 comment('阶段级别') INT(11)" structs:"level,omitempty"`
	Threshold           int    `xorm:"not null default 0 comment('时间门槛') INT(11)" structs:"threshold,omitempty"`
	FieldId             int    `xorm:"unique(name) not null default 0 comment('隶属方向') INT(11)" structs:"field_id,omitempty"`
}

func NewPhase() *Phase {
	ret := Phase{}
	return &ret
}

func (m *Phase) TableName() string {
	return "phase"
}

func (m *Phase) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	session.Asc("level")
	return
}

func (m *Phase) SlicePtr() interface{} {
	ret := make([]Phase, 0)
	return &ret
}

func (m *Phase) Transfer(slicePtr interface{}) *[]Phase {
	ret := slicePtr.(*[]Phase)
	return ret
}

func (m *Phase) WithDeleted() bool {
	return true
}
