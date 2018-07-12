package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Resource struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"unique not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Desc                string `xorm:"not null default '' comment('描述') VARCHAR(255)" structs:"desc,omitempty"`
	Year                int    `xorm:"not null default 0 comment('年份') INT(11)" structs:"year,omitempty"`
	AreaId              int    `xorm:"unique(name) not null default 0 comment('隶属领域') INT(11)" structs:"area_id,omitempty"`
	Area                Area   `xorm:"-" structs:"-"`
	WithSub             bool   `xorm:"-" structs:"-"`
}

type ResourceJoin struct {
	Resource `xorm:"extends" json:"-"`
	Area     `xorm:"extends" json:"-"`
}

type ResourceSyncfusion struct {
	Id       int    `json:"Id"`
	ParentId int    `json:"ParentId"`
	IsParent bool   `json:"IsParent"`
	Name     string `json:"Name"`
	Parent   string `json:"Parent"`
}

func (m *Resource) TableName() string {
	return "resource"
}

func (m *Resource) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)

	areaCondition := m.Area.BuildCondition()
	if len(areaCondition) > 0 {
		for k, v := range areaCondition {
			condition[k] = v
		}
	}

	return condition
}
