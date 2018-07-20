package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Project struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string    `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	QuestTargetId       int       `xorm:"not null default 0 comment('所属target') INT(11)" structs:"quest_target_id,omitempty"`
	StartDate           time.Time `xorm:"not null comment('开始日期') DATETIME"`
	Duration            int       `xorm:"not null default 0 comment('持续时间') INT(11)" structs:"duration,omitempty"`

	QuestTarget QuestTarget `xorm:"-" structs:"-"`
	Area        Area        `xorm:"-" structs:"-"`
	Quest       Quest       `xorm:"-" structs:"-"`
}

type ProjectJoin struct {
	Project     `xorm:"extends" json:"-"`
	QuestTarget `xorm:"extends" json:"-"`
	Area        `xorm:"extends" json:"-"`
}

func (m *Project) TableName() string {
	return "project"
}

func (m *Project) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	m.QuestTarget.BuildCondition(session)
	return
}

func (m *Project) SlicePtr() interface{} {
	ret := make([]Project, 0)
	return &ret
}

func (m *Project) Join() es.JoinGeneral {
	ret := ProjectJoin{}
	return &ret
}

func (j *ProjectJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	questTargetLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.QuestTarget.TableName(),
		LeftTable:  j.QuestTarget.TableName(),
		LeftField:  "id",
		RightTable: j.Project.TableName(),
		RightField: "quest_target_id",
	}
	areaLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Area.TableName(),
		LeftTable:  j.Area.TableName(),
		LeftField:  "id",
		RightTable: j.QuestTarget.TableName(),
		RightField: "area_id",
	}
	links = append(links, questTargetLink)
	links = append(links, areaLink)
	return links
}

func (j *ProjectJoin) SlicePtr() interface{} {
	ret := make([]ProjectJoin, 0)
	return &ret
}

func (j *ProjectJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.Project
	ret.Area = join.Area
	ret.QuestTarget = join.QuestTarget
	return &ret
}

func (j *ProjectJoin) TransferSlicePtr(slicePtr interface{}) interface{} {
	joinSlicePtr := slicePtr.(*[]ProjectJoin)
	joinSlice := *joinSlicePtr
	projects := make([]Project, 0)
	for _, one := range joinSlice {
		projectPtr := (&one).Transfer().(*Project)
		projects = append(projects, *projectPtr)
	}
	return &projects
}
