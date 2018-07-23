package models

import (
	"encoding/binary"
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
	uuid "github.com/satori/go.uuid"
)

type Project struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string    `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	QuestTargetId       int       `xorm:"not null default 0 comment('所属target') INT(11)" structs:"quest_target_id,omitempty"`
	StartDate           time.Time `xorm:"not null comment('开始日期') DATETIME"`
	Status              int       `xorm:"not null default 1 comment('当前状态:1未完成2已完成') INT(11)" structs:"status,omitempty"`
	Uuid                string    `xorm:"not null default '' comment('唯一ID') VARCHAR(255)" structs:"uuid,omitempty"`
	UuidNumber          uint32    `xorm:"-" structs:"-"`

	Quest       *Quest       `xorm:"-" structs:"-" json:"Quest,omitempty"`
	QuestTarget *QuestTarget `xorm:"-" structs:"-" json:"QuestTarget,omitempty"`
	Area        *Area        `xorm:"-" structs:"-" json:"Area,omitempty"`
}

var (
	ProjectStatusWait   = 1
	ProjectStatusFinish = 2
)

type ProjectJoin struct {
	Project     Project `xorm:"extends" json:"-"`
	QuestTarget `xorm:"extends" json:"-"`
	Area        `xorm:"extends" json:"-"`
	Quest       `xorm:"extends" json:"-"`
}

func NewProject() *Project {
	ret := Project{
		QuestTarget: NewQuestTarget(),
		Area:        NewArea(),
		Quest:       NewQuest(),
	}
	return &ret
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

func (m *Project) Transfer(slicePtr interface{}) *[]Project {
	ret := slicePtr.(*[]Project)
	return ret
}

func (m *Project) Hook() (err error) {
	if len(m.Uuid) == 0 {
		u := uuid.NewV4()
		m.Uuid = u.String()
		m.UuidNumber = binary.BigEndian.Uint32(u[0:4])
	} else {
		u, err := uuid.FromString(m.Uuid)
		if err != nil {
			return err
		}
		m.UuidNumber = binary.BigEndian.Uint32(u[0:4])
	}
	return nil
}

func (m *Project) Join() es.JoinGeneral {
	ret := ProjectJoin{}
	return &ret
}

func (m *Project) WithDeleted() bool {
	return true
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
	questLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Quest.TableName(),
		LeftTable:  j.Quest.TableName(),
		LeftField:  "id",
		RightTable: j.QuestTarget.TableName(),
		RightField: "quest_id",
	}
	links = append(links, questTargetLink)
	links = append(links, areaLink)
	links = append(links, questLink)
	return links
}

func (j *ProjectJoin) SlicePtr() interface{} {
	ret := make([]ProjectJoin, 0)
	return &ret
}

func (j *ProjectJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.Project
	ret.Area = &join.Area
	ret.Quest = &join.Quest
	ret.QuestTarget = &join.QuestTarget
	return &ret
}

func (j *ProjectJoin) TransferCopy(modelPtr es.ModelGeneral) {
	projectPtr := modelPtr.(*Project)
	(*projectPtr) = (*j).Project
	(*projectPtr).Area = &(*j).Area
	(*projectPtr).Quest = &(*j).Quest
	(*projectPtr).QuestTarget = &(*j).QuestTarget
	return
}

func (j *ProjectJoin) TransferCopySlice(slicePtr interface{}, targetPtr interface{}) {
	joinSlicePtr := slicePtr.(*[]ProjectJoin)
	joinSlice := *joinSlicePtr
	projects := make([]Project, 0)
	for _, one := range joinSlice {
		projectPtr := (&one).Transfer().(*Project)
		projects = append(projects, *projectPtr)
	}
	projectsPtr := targetPtr.(*[]Project)
	(*projectsPtr) = projects
	return
}
