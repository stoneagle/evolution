package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Concept struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"unique not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
}

func NewConcept() *Concept {
	ret := Concept{}
	return &ret
}

func (m *Concept) TableName() string {
	return "concept"
}

func (m *Concept) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	session.Asc("code")
	return
}

func (m *Concept) SlicePtr() interface{} {
	ret := make([]Concept, 0)
	return &ret
}

func (m *Concept) Transfer(slicePtr interface{}) *[]Concept {
	ret := slicePtr.(*[]Concept)
	return ret
}

func (m *Concept) WithDeleted() bool {
	return true
}
