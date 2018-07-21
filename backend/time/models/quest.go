package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Quest struct {
	es.ModelWithDeleted `xorm:"extends"`
	Name                string    `xorm:"not null default '' comment('内容') VARCHAR(255)" structs:"name,omitempty"`
	StartDate           time.Time `xorm:"comment('开始日期') DATETIME" structs:"start_date,omitempty"`
	EndDate             time.Time `xorm:"not null comment('结束日期') DATETIME" structs:"start_date,omitempty"`
	FounderId           int       `xorm:"not null default 0 comment('创建人id') INT(11)" structs:"founder_id,omitempty"`
	Members             int       `xorm:"not null default 0 comment('团队人数') INT(11)" structs:"members,omitempty"`
	Constraint          int       `xorm:"not null default 0 comment('限制:1重要紧急2重要不紧急3紧急不重要4不重要不紧急') INT(11)" structs:"constraint,omitempty"`
	Status              int       `xorm:"not null default 0 comment('当前状态:1招募中2执行中3已完成4未完成') INT(11)" structs:"status,omitempty"`
}

var (
	QuestStatusRecruit int = 1
	QuestStatusExec    int = 2
	QuestStatusFinish  int = 3
	QuestStatusFail    int = 4

	QuestConstraintImportantAndBusy int = 1
	QuestConstraintImportant        int = 2
	QuestConstraintBusy             int = 3
	QuestConstraintNormal           int = 4
)

func NewQuest() *Quest {
	ret := Quest{}
	return &ret
}

func (m *Quest) TableName() string {
	return "quest"
}

func (m *Quest) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	return
}

func (m *Quest) SlicePtr() interface{} {
	ret := make([]Quest, 0)
	return &ret
}

func (m *Quest) Transfer(slicePtr interface{}) *[]Quest {
	ret := slicePtr.(*[]Quest)
	return ret
}
