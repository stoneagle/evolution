package models

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type Area struct {
	GeneralWithId `xorm:"extends"`
	Name          string `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Parent        int    `xorm:"not null default 0 comment('父id') INT(11)" structs:"parent,omitempty"`
	Del           int    `xorm:"not null default 0 comment('软删除') TINYINT(1)" structs:"del,omitempty"`
	Level         int    `xorm:"not null default 0 comment('所属层级') SMALLINT(6)" structs:"level,omitempty"`
	FieldId       int    `xorm:"not null default 0 comment('所属范畴id') SMALLINT(6)" structs:"field_id,omitempty"`
}

type AreaTree struct {
	Value    string     `json:"value"`
	Children []AreaNode `json:"children"`
}

type AreaNode struct {
	Id       int        `json:"id"`
	Value    string     `json:"value"`
	Children []AreaNode `json:"children"`
}

const (
	AreaFieldLife = iota
	AreaFieldSkill
	AreaFieldAsset
	AreaFieldWork
	AreaFieldCircle
	AreaFieldQuest
)

var (
	AreaFieldMap map[int]string = map[int]string{
		AreaFieldLife:   "日常",
		AreaFieldSkill:  "知识",
		AreaFieldAsset:  "财富",
		AreaFieldWork:   "文化",
		AreaFieldCircle: "社交",
		AreaFieldQuest:  "挑战",
	}
)

func (m *Area) TableName() string {
	return "area"
}

func (m *Area) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	condition = builder.Eq{}

	params := structs.Map(m)
	for key, value := range params {
		keyType := reflect.ValueOf(value).Kind()
		if keyType != reflect.Map && keyType != reflect.Slice && keyType != reflect.Struct {
			condition[keyPrefix+key] = value
		}
	}
	return condition
}
