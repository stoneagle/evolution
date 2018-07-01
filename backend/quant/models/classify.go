package models

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Classify struct {
	GeneralWithDeleted `xorm:"extends"`
	AssetType          `xorm:"extends"`
	Source             `xorm:"extends"`
	Name               string `xorm:"varchar(128) notnull comment('名称')" form:"Name" json:"Name" structs:"name,omitempty"`
	Tag                string `xorm:"varchar(128) notnull comment('标签')" form:"Tag" json:"Tag" structs:"tag,omitempty"`
}

func (m *Classify) TableName() string {
	return "classify"
}

func (m *Classify) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	condition = builder.Eq{}
	general := structs.Map(m.GeneralWithDeleted)
	for key, value := range general {
		condition[keyPrefix+key] = value
	}

	assetType := structs.Map(m.AssetType)
	for key, value := range assetType {
		condition[keyPrefix+key] = value
	}

	source := structs.Map(m.Source)
	for key, value := range source {
		condition[keyPrefix+key] = value
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
