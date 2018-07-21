package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Field struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	Desc                string `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Color               string `xorm:"not null default '' comment('颜色') VARCHAR(255)"`
}

func NewField() *Field {
	ret := Field{}
	return &ret
}

func (m *Field) TableName() string {
	return "field"
}

func (m *Field) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	return
}

func (m *Field) SlicePtr() interface{} {
	ret := make([]Field, 0)
	return &ret
}

func (m *Field) Transfer(slicePtr interface{}) *[]Field {
	ret := slicePtr.(*[]Field)
	return ret
}
