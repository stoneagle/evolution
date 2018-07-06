package models

import (
	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Entity struct {
	GeneralWithDeleted `xorm:"extends"`
	Name               string `xorm:"unique not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Desc               string `xorm:"not null default '' comment('描述') VARCHAR(255)" structs:"desc,omitempty"`
	Year               int    `xorm:"not null default 0 comment('年份') INT(11)" structs:"year,omitempty"`
	AreaId             int    `xorm:"unique(name) not null default 0 comment('隶属领域') INT(11)" structs:"area_id,omitempty"`
	Area               Area   `xorm:"-"`
}

type EntityJoin struct {
	Entity `xorm:"extends" json:"-"`
	Area   `xorm:"extends" json:"-"`
}

func (m *Entity) TableName() string {
	return "entity"
}

func (m *Entity) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.General.BuildCondition(params, keyPrefix)

	areaCondition := m.Area.BuildCondition()
	if len(areaCondition) > 0 {
		for k, v := range areaCondition {
			condition[k] = v
		}
	}
	return condition
}
