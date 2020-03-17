package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type DailyLimitUp struct {
	es.ModelWithDeleted `xorm:"extends"`
	ShareId             int       `xorm:"unique not null default 0 comment('所属个股') INT(11)" structs:"share_id,omitempty"`
	Date                time.Time `xorm:"not null comment('日期') DATETIME" structs:"date,omitempty"`
	LimitTime           time.Time `xorm:"not null comment('涨停时间') DATETIME" structs:"limit_time,omitempty"`
	LimitCompetition    int       `xorm:"not null default 0 comment('竞价涨停') SMALLINT(4)" structs:"limit_competition,omitempty"`
	LimitContinuation   int       `xorm:"not null default 0 comment('连板次数') INT(11)" structs:"limit_continuation,omitempty"`
	HighTransaction     int       `xorm:"not null default 0 comment('高位成交-大于8%') INT(11)" structs:"high_transaction,omitempty"`
	Circulate           int       `xorm:"not null default 0 comment('实际流通') INT(11)" structs:"circulate,omitempty"`
	MainConceptId       int       `xorm:"not null default 0 comment('主要概念') INT(11)" structs:"main_concept_id,omitempty"`
	AdditionalConceptId int       `xorm:"not null default 0 comment('叠加概念') INT(11)" structs:"additional_concept_id,omitempty"`
	Share               *Share    `xorm:"-" structs:"-" json:"Share,omitempty"`
}

type DailyLimitUpJoin struct {
	DailyLimitUp `xorm:"extends" json:"-"`
	Share        `xorm:"extends" json:"-"`
}

var (
	DailyLimitUpStatusWait   = 1
	DailyLimitUpStatusFinish = 2
)

func NewDailyLimitUp() *DailyLimitUp {
	ret := DailyLimitUp{
		Share: NewShare(),
	}
	return &ret
}

func (m *DailyLimitUp) TableName() string {
	return "daily_limit_up"
}

func (m *DailyLimitUp) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	m.Share.BuildCondition(session)
	return
}

func (m *DailyLimitUp) SlicePtr() interface{} {
	ret := make([]DailyLimitUp, 0)
	return &ret
}

func (m *DailyLimitUp) Transfer(slicePtr interface{}) *[]DailyLimitUp {
	ret := slicePtr.(*[]DailyLimitUp)
	return ret
}

func (m *DailyLimitUp) Join() es.JoinGeneral {
	ret := DailyLimitUpJoin{}
	return &ret
}

func (j *DailyLimitUpJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	areaLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Share.TableName(),
		LeftTable:  j.Share.TableName(),
		LeftField:  "id",
		RightTable: j.DailyLimitUp.TableName(),
		RightField: "area_id",
	}
	links = append(links, areaLink)
	return links
}

func (j *DailyLimitUpJoin) SlicePtr() interface{} {
	ret := make([]DailyLimitUpJoin, 0)
	return &ret
}

func (j *DailyLimitUpJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.DailyLimitUp
	ret.Share = &join.Share
	return &ret
}

func (j *DailyLimitUpJoin) TransferCopy(modelPtr es.ModelGeneral) {
	dailyLimitupPtr := modelPtr.(*DailyLimitUp)
	(*dailyLimitupPtr) = (*j).DailyLimitUp
	(*dailyLimitupPtr).Share = &(*j).Share
	return
}

func (j *DailyLimitUpJoin) TransferCopySlice(slicePtr interface{}, targetPtr interface{}) {
	joinSlicePtr := slicePtr.(*[]DailyLimitUpJoin)
	joinSlice := *joinSlicePtr
	dailyLimitups := make([]DailyLimitUp, 0)
	for _, one := range joinSlice {
		dailyLimitupPtr := (&one).Transfer().(*DailyLimitUp)
		dailyLimitups = append(dailyLimitups, *dailyLimitupPtr)
	}
	dailyLimitupsPtr := targetPtr.(*[]DailyLimitUp)
	(*dailyLimitupsPtr) = dailyLimitups
	return
}

func (m *DailyLimitUp) WithDeleted() bool {
	return true
}
