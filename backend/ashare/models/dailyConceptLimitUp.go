package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type DailyConceptLimitUp struct {
	es.ModelWithDeleted  `xorm:"extends"`
	ConceptId            int       `xorm:"unique not null default 0 comment('所属概念') INT(11)" structs:"concept_id,omitempty"`
	Date                 time.Time `xorm:"not null comment('日期') DATETIME" structs:"date,omitempty"`
	LimitUpCount         int       `xorm:"not null default 0 comment('总涨停') SMALLINT(6)" structs:"limit_up_sum,omitempty"`
	LimitCompetitionSum  int       `xorm:"not null default 0 comment('总竞价涨停') SMALLINT(6)" structs:"limit_competition_sum,omitempty"`
	LimitContinuationSum int       `xorm:"not null default 0 comment('总连板') SMALLINT(6)" structs:"limit_continuation_sum,omitempty"`
	LimitContinuationMax int       `xorm:"not null default 1 comment('最高连板') SMALLINT(6)" structs:"limit_continuation_max,omitempty"`
	CirculateSum         int       `xorm:"not null default 0 comment('实际流通') INT(11)" structs:"circulate_sum,omitempty"`
	TransactionSum       int       `xorm:"not null default 0 comment('总成交') INT(11)" structs:"transaction_sum,omitempty"`
	HighTransactionSum   int       `xorm:"not null default 0 comment('总高位成交-大于8%') INT(11)" structs:"high_transaction_sum,omitempty"`
}

func NewDailyConceptLimitUp() *DailyConceptLimitUp {
	ret := DailyConceptLimitUp{}
	return &ret
}

func (m *DailyConceptLimitUp) TableName() string {
	return "daily_concept_limit_up"
}

func (m *DailyConceptLimitUp) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	session.Asc("date")
	return
}

func (m *DailyConceptLimitUp) SlicePtr() interface{} {
	ret := make([]DailyConceptLimitUp, 0)
	return &ret
}

func (m *DailyConceptLimitUp) Transfer(slicePtr interface{}) *[]DailyConceptLimitUp {
	ret := slicePtr.(*[]DailyConceptLimitUp)
	return ret
}

func (m *DailyConceptLimitUp) WithDeleted() bool {
	return true
}
