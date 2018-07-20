package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Country struct {
	es.ModelWithId `xorm:"extends"`
	Name           string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	EnName         string `xorm:"not null default '' comment('英文名称') VARCHAR(255)"`
}

func (m *Country) TableName() string {
	return "country"
}

func (m *Country) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	return
}

func (m *Country) SlicePtr() interface{} {
	ret := make([]Country, 0)
	return &ret
}
