package models

import (
	"reflect"

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
	condition = builder.Eq{}

	areaCondition := m.Area.BuildCondition()
	if len(areaCondition) > 0 {
		for k, v := range areaCondition {
			condition[k] = v
		}
	}

	params := structs.Map(m)
	for key, value := range params {
		keyType := reflect.ValueOf(value).Kind()
		if keyType != reflect.Map && keyType != reflect.Slice && keyType != reflect.Struct {
			condition[keyPrefix+key] = value
		}
	}
	return condition
}
