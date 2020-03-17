package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type TopFundTransaction struct {
	es.ModelWithDeleted `xorm:"extends"`
	ShareId             int       `xorm:"unique not null default 0 comment('所属个股') INT(11)" structs:"share_id,omitempty"`
	TopFundId           int       `xorm:"unique not null default 0 comment('所属资金') INT(11)" structs:"top_fund_id,omitempty"`
	Date                time.Time `xorm:"not null comment('日期') DATETIME" structs:"date,omitempty"`
	Type                int       `xorm:"not null default 1 comment('行为1买入2卖出') SMALLINT(4)" structs:"type,omitempty"`
	Transaction         int       `xorm:"not null default 0 comment('成交金额') INT(11)" structs:"transaction,omitempty"`
}

var (
	TopFundTransactionTypeBuy  = 1
	TopFundTransactionTypeSell = 2
)

func NewTopFundTransaction() *TopFundTransaction {
	ret := TopFundTransaction{}
	return &ret
}

func (m *TopFundTransaction) TableName() string {
	return "top_fund_transaction"
}

func (m *TopFundTransaction) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	session.Asc("date")
	return
}

func (m *TopFundTransaction) SlicePtr() interface{} {
	ret := make([]TopFundTransaction, 0)
	return &ret
}

func (m *TopFundTransaction) Transfer(slicePtr interface{}) *[]TopFundTransaction {
	ret := slicePtr.(*[]TopFundTransaction)
	return ret
}

func (m *TopFundTransaction) WithDeleted() bool {
	return true
}
