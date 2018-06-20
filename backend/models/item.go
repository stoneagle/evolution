package models

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Item struct {
	GeneralWithDeleted `xorm:"extends"`
	Code               string     `xorm:"varchar(16) notnull unique comment('编号')" form:"Code" json:"Code" structs:"code,omitempty"`
	Name               string     `xorm:"varchar(32) notnull comment('名称')" form:"Name" json:"Name" structs:"name,omitempty"`
	Status             int        `xorm:"smallint(4) notnull comment('状态')" form:"Status" json:"Status" structs:"status,omitempty"`
	Classify           []Classify `xorm:"-"`
}

func (m *Item) TableName() string {
	return "item"
}

func (m *Item) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	condition = builder.Eq{}
	general := structs.Map(m.GeneralWithDeleted)
	for key, value := range general {
		condition[keyPrefix+key] = value
	}

	if len(m.Classify) > 0 {
		subCondition := m.Classify[0].BuildCondition()
		for key, value := range subCondition {
			condition[key] = value
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
