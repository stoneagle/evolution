package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type MapAreaResource struct {
	es.Model   `xorm:"extends"`
	AreaId     int `xorm:"unique not null default 0 comment('隶属领域') INT(11)" structs:"area_id,omitempty"`
	ResourceId int `xorm:"unique(area_id) not null default 0 comment('隶属资源') INT(11)" structs:"resource_id,omitempty"`
}

func NewMapAreaResource() *MapAreaResource {
	ret := MapAreaResource{}
	return &ret
}

func (m *MapAreaResource) TableName() string {
	return "map_area_resource"
}

func (m *MapAreaResource) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	return
}

func (m *MapAreaResource) SlicePtr() interface{} {
	ret := make([]MapAreaResource, 0)
	return &ret
}
