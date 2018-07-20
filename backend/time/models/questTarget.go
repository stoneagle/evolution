package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type QuestTarget struct {
	es.ModelWithId `xorm:"extends"`
	QuestId        int    `xorm:"unique not null default 0 comment('所属quest') INT(11)" structs:"quest_id,omitempty"`
	AreaId         int    `xorm:"unique(quest_id) not null default 0 comment('目标资源id') INT(11)" structs:"area_id,omitempty"`
	Desc           string `xorm:"not null default '' comment('描述') VARCHAR(255)" structs:"desc,omitempty"`
	Status         int    `xorm:"not null default 1 comment('当前状态:1未完成2已完成') INT(11)" structs:"status,omitempty"`

	QuestIds []int `xorm:"-" structs:"quest_id,omitempty"`
	Area     Area  `xorm:"-" structs:"-"`
}

type QuestTargetJoin struct {
	QuestTarget `xorm:"extends" json:"-"`
	Area        `xorm:"extends" json:"-"`
}

var (
	QuestTargetStatusWait   = 1
	QuestTargetStatusFinish = 2
)

func (m *QuestTarget) TableName() string {
	return "quest_target"
}

func (m *QuestTarget) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	m.Area.BuildCondition(session)
	return
}

func (m *QuestTarget) SlicePtr() interface{} {
	ret := make([]QuestTarget, 0)
	return &ret
}

func (m *QuestTarget) Join() es.JoinGeneral {
	ret := QuestTargetJoin{}
	return &ret
}

func (j *QuestTargetJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	areaLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Area.TableName(),
		LeftTable:  j.Area.TableName(),
		LeftField:  "id",
		RightTable: j.QuestTarget.TableName(),
		RightField: "area_id",
	}
	links = append(links, areaLink)
	return links
}

func (j *QuestTargetJoin) SlicePtr() interface{} {
	ret := make([]QuestTargetJoin, 0)
	return &ret
}

func (j *QuestTargetJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.QuestTarget
	ret.Area = join.Area
	return &ret
}

func (j *QuestTargetJoin) TransferCopy(modelPtr es.ModelGeneral) {
	questTargetPtr := modelPtr.(*QuestTarget)
	(*questTargetPtr) = (*j).QuestTarget
	(*questTargetPtr).Area = (*j).Area
	return
}

func (j *QuestTargetJoin) TransferSlicePtr(slicePtr interface{}) interface{} {
	joinSlicePtr := slicePtr.(*[]QuestTargetJoin)
	joinSlice := *joinSlicePtr
	questTargets := make([]QuestTarget, 0)
	for _, one := range joinSlice {
		questTargetPtr := (&one).Transfer().(*QuestTarget)
		questTargets = append(questTargets, *questTargetPtr)
	}
	return &questTargets
}
