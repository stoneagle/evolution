package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type TopFund struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"unique not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Type                int    `xorm:"not null default 0 comment('资金类别1国家2派系3营业厅4牛散') SMALLINT(6)" structs:"type,omitempty"`
}

var (
	TopFundTypeCountry = 1
	TopFundTypeSite    = 2
	TopFundTypeHall    = 3
	TopFundTypePerson  = 4
)

func NewTopFund() *TopFund {
	ret := TopFund{}
	return &ret
}

func (m *TopFund) TableName() string {
	return "top_fund"
}

func (m *TopFund) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	session.Asc("type")
	return
}

func (m *TopFund) SlicePtr() interface{} {
	ret := make([]TopFund, 0)
	return &ret
}

func (m *TopFund) Transfer(slicePtr interface{}) *[]TopFund {
	ret := slicePtr.(*[]TopFund)
	return ret
}

func (m *TopFund) WithDeleted() bool {
	return true
}
