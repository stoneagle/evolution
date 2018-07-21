package models

import (
	es "evolution/backend/common/structs"
	"fmt"
	"time"

	"github.com/fatih/structs"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

type Action struct {
	es.ModelWithId `xorm:"extends"`
	TaskId         int       `xorm:"not null default 0 comment('所属task') INT(11)" structs:"task_id,omitempty"`
	StartDate      time.Time `xorm:"TIMESTAMP not null comment('开始时间') DATETIME" structs:"-"`
	EndDate        time.Time `xorm:"TIMESTAMP not null comment('结束时间') DATETIME" structs:"-"`
	Name           string    `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	UserId         int       `xorm:"not null default 0 comment('执行人') INT(11)" structs:"user_id,omitempty"`
	Time           int       `xorm:"not null default 0 comment('花费时间:单位分钟') INT(11)" structs:"time,omitempty"`

	TaskIds []int `xorm:"-" structs:"task_id,omitempty" json:"TaskIds,omitempty"`
	Task    *Task `xorm:"-" structs:"-" json:"Task,omitempty"`
}

type ActionJoin struct {
	Action `xorm:"extends" json:"-"`
	Task   `xorm:"extends" json:"-"`
}

func NewAction() *Action {
	ret := Action{
		Task: NewTask(),
	}
	return &ret
}

func (m *Action) TableName() string {
	return "action"
}

func (m *Action) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)

	emptyTime := time.Time{}
	if m.StartDate != emptyTime {
		session.And(builder.Gte{fmt.Sprintf("%s.start_date", m.TableName()): m.StartDate})
	}
	if m.EndDate != emptyTime {
		session.And(builder.Lte{fmt.Sprintf("%s.start_date", m.TableName()): m.EndDate})
	}

	return
}

func (m *Action) SlicePtr() interface{} {
	ret := make([]Action, 0)
	return &ret
}

func (m *Action) Transfer(slicePtr interface{}) *[]Action {
	ret := slicePtr.(*[]Action)
	return ret
}

func (m *Action) Join() es.JoinGeneral {
	ret := ActionJoin{}
	return &ret
}

func (j *ActionJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	actionLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Task.TableName(),
		LeftTable:  j.Task.TableName(),
		LeftField:  "id",
		RightTable: j.Action.TableName(),
		RightField: "task_id",
	}
	links = append(links, actionLink)
	return links
}

func (j *ActionJoin) SlicePtr() interface{} {
	ret := make([]ActionJoin, 0)
	return &ret
}

func (j *ActionJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.Action
	ret.Task = &join.Task
	return &ret
}

func (j *ActionJoin) TransferCopy(modelPtr es.ModelGeneral) {
	actionPtr := modelPtr.(*Action)
	(*actionPtr) = (*j).Action
	(*actionPtr).Task = &(*j).Task
	return
}

func (j *ActionJoin) TransferCopySlice(slicePtr interface{}, targetPtr interface{}) {
	joinSlicePtr := slicePtr.(*[]ActionJoin)
	joinSlice := *joinSlicePtr
	actions := make([]Action, 0)
	for _, one := range joinSlice {
		actionPtr := (&one).Transfer().(*Action)
		actions = append(actions, *actionPtr)
	}
	actionsPtr := targetPtr.(*[]Action)
	(*actionsPtr) = actions
	return
}
