package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Share struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"unique not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Code                string `xorm:"unique(name) not null default '' comment('代码') VARCHAR(64)" structs:"code,omitempty"`
	IndexType           int    `xorm:"not null default 0 comment('所属指数1上证2深证3创业') SMALLINT(6)" structs:"index_type,omitempty"`
}

var (
	ShareIndexTypeShanghai = 1
	ShareIndexTypeShenzhen = 2
	ShareIndexTypeChuangye = 3
)

func NewShare() *Share {
	ret := Share{}
	return &ret
}

func (m *Share) TableName() string {
	return "share"
}

func (m *Share) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	session.Asc("code")
	return
}

func (m *Share) SlicePtr() interface{} {
	ret := make([]Share, 0)
	return &ret
}

func (m *Share) Transfer(slicePtr interface{}) *[]Share {
	ret := slicePtr.(*[]Share)
	return ret
}

func (m *Share) WithDeleted() bool {
	return true
}
