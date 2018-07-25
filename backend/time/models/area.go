package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Area struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string     `xorm:"unique not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Parent              int        `xorm:"unique(name) not null default 0 comment('父id') INT(11)" structs:"parent,omitempty"`
	FieldId             int        `xorm:"not null default 0 comment('所属范畴id') SMALLINT(6)" structs:"field_id,omitempty"`
	Type                int        `xorm:"not null default 0 comment('类别1根2节点3叶') SMALLINT(6)" structs:"type,omitempty"`
	Resources           []Resource `xorm:"-" structs:"-" json:"Resources,omitempty"`
}

var (
	AreaTypeRoot = 1
	AreaTypeNode = 2
	AreaTypeLeaf = 3
)

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

func NewArea() *Area {
	ret := Area{}
	return &ret
}

func (m *Area) TableName() string {
	return "area"
}

func (m *Area) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	session.Asc("parent")
	return
}

func (m *Area) SlicePtr() interface{} {
	ret := make([]Area, 0)
	return &ret
}

func (m *Area) Transfer(slicePtr interface{}) *[]Area {
	ret := slicePtr.(*[]Area)
	return ret
}

func (m *Area) WithDeleted() bool {
	return true
}
