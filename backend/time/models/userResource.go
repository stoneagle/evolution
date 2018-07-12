package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type UserResource struct {
	es.ModelWithId `xorm:"extends"`
	UserId         int `xorm:"unique not null default 0 comment('隶属用户') INT(11)" structs:"user_id,omitempty"`
	ResourceId     int `xorm:"unique(user_id) not null default 0 comment('隶属实体') INT(11)" structs:"resource_id,omitempty"`
	Time           int `xorm:"not null default 0 comment('总投入时间') INT(11)" structs:"time,omitempty"`
}

func (m *UserResource) TableName() string {
	return "user_resource"
}

func (m *UserResource) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
