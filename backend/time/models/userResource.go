package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type UserResource struct {
	es.ModelWithId `xorm:"extends"`
	UserId         int      `xorm:"unique not null default 0 comment('隶属用户') INT(11)" structs:"user_id,omitempty"`
	ResourceId     int      `xorm:"unique(user_id) not null default 0 comment('隶属实体') INT(11)" structs:"resource_id,omitempty"`
	Time           int      `xorm:"not null default 0 comment('总投入时间') INT(11)" structs:"time,omitempty"`
	Resource       Resource `xorm:"-"`
}

type UserResourceJoin struct {
	UserResource    `xorm:"extends" json:"-"`
	Resource        `xorm:"extends" json:"-"`
	MapAreaResource `xorm:"extends" json:"-"`
	Area            `xorm:"extends" json:"-"`
}

var (
	UserResourceStatusRelax = 1
	UserResourceStatusExec  = 2
)

func (m *UserResource) TableName() string {
	return "user_resource"
}

func (m *UserResource) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	m.Resource.BuildCondition(session)
	return
}

func (m *UserResource) SlicePtr() interface{} {
	ret := make([]UserResource, 0)
	return &ret
}

func (m *UserResource) Join() es.JoinGeneral {
	ret := UserResourceJoin{}
	return &ret
}

func (j *UserResourceJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	resourceLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Resource.TableName(),
		LeftTable:  j.Resource.TableName(),
		LeftField:  "id",
		RightTable: j.UserResource.TableName(),
		RightField: "resource_id",
	}
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
	links = append(links, resourceLink)
	links = append(links, mapAreaResourceLink)
	links = append(links, areaLink)
	return links
}

func (j *UserResourceJoin) SlicePtr() interface{} {
	ret := make([]UserResourceJoin, 0)
	return &ret
}

func (j *UserResourceJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.UserResource
	ret.Resource = join.Resource
	ret.Resource.Area = join.Area
	return &ret
}

func (j *UserResourceJoin) TransferSlicePtr(slicePtr interface{}) interface{} {
	joinSlicePtr := slicePtr.(*[]UserResourceJoin)
	joinSlice := *joinSlicePtr
	userResources := make([]UserResource, 0)
	for _, one := range joinSlice {
		userResourcePtr := (&one).Transfer().(*UserResource)
		userResources = append(userResources, *userResourcePtr)
	}
	return &userResources
}
