package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Resource struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"unique not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Desc                string `xorm:"not null default '' comment('描述') VARCHAR(255)" structs:"desc,omitempty"`
	Year                int    `xorm:"not null default 0 comment('年份') INT(11)" structs:"year,omitempty"`

	Area    Area `xorm:"-"`
	WithSub bool `xorm:"-" structs:"-"`
}

type ResourceJoin struct {
	MapAreaResource `xorm:"extends" json:"-"`
	Resource        `xorm:"extends" json:"-"`
	Area            `xorm:"extends" json:"-"`
}

func (m *Resource) TableName() string {
	return "resource"
}

func (m *Resource) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	m.Area.BuildCondition(session)
	return
}

func (m *Resource) SlicePtr() interface{} {
	ret := make([]Resource, 0)
	return &ret
}

func (m *Resource) Join() es.JoinGeneral {
	ret := ResourceJoin{}
	return &ret
}

func (j *ResourceJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	mapAreaResourceLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.MapAreaResource.TableName(),
		LeftTable:  j.MapAreaResource.TableName(),
		LeftField:  "resource_id",
		RightTable: j.Resource.TableName(),
		RightField: "id",
	}
	areaLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Area.TableName(),
		LeftTable:  j.Area.TableName(),
		LeftField:  "id",
		RightTable: j.MapAreaResource.TableName(),
		RightField: "area_id",
	}
	links = append(links, mapAreaResourceLink)
	links = append(links, areaLink)
	return links
}

func (j *ResourceJoin) SlicePtr() interface{} {
	ret := make([]ResourceJoin, 0)
	return &ret
}

func (j *ResourceJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.Resource
	ret.Area = join.Area
	return &ret
}

func (j *ResourceJoin) TransferSlicePtr(slicePtr interface{}) interface{} {
	joinSlicePtr := slicePtr.(*[]ResourceJoin)
	joinSlice := *joinSlicePtr
	resources := make([]Resource, 0)
	for _, one := range joinSlice {
		resourcePtr := (&one).Transfer().(*Resource)
		resources = append(resources, *resourcePtr)
	}
	return &resources
}
