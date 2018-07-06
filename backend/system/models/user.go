package models

import (
	es "evolution/backend/common/structs"

	"github.com/fatih/structs"
	"github.com/go-xorm/builder"
)

type User struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string `xorm:"not null unique default '' comment('名称') VARCHAR(255)"`
	Email               string `xorm:"not null default '' comment('邮箱') VARCHAR(255)"`
	Password            string `xorm:"not null default '' comment('密码') VARCHAR(255)"`
}

func (m *User) TableName() string {
	return "user"
}

func (m *User) BuildCondition() (condition builder.Eq) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition = m.Model.BuildCondition(params, keyPrefix)
	return condition
}
