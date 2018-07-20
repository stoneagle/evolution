package models

import (
	es "evolution/backend/common/structs"
	"time"

	"github.com/fatih/structs"
	"github.com/go-xorm/xorm"
)

type Task struct {
	es.ModelWithDeleted `xorm:"extends"`
	ProjectId           int       `xorm:"not null default 0 comment('所属项目') INT(11)" structs:"project_id,omitempty"`
	ResourceId          int       `xorm:"not null default 0 comment('所属资源') INT(11)" structs:"resource_id,omitempty"`
	UserId              int       `xorm:"not null default 0 comment('负责人') INT(11)" structs:"user_id,omitempty"`
	Name                string    `xorm:"not null default '' comment('名称') VARCHAR(255)" structs:"name,omitempty"`
	Desc                string    `xorm:"comment('描述') TEXT" structs:"desc,omitempty"`
	StartDate           time.Time `xorm:"comment('开始日期') DATETIME"`
	EndDate             time.Time `xorm:"comment('结束日期') DATETIME"`
	Duration            int       `xorm:"not null comment('持续时间') INT(11)" structs:"duration,omitempty"`
	Status              int       `xorm:"not null default 0 comment('当前状态:1未分配2准备做3执行中4已完成') INT(11)" structs:"status,omitempty"`

	ProjectIds     []int    `xorm:"-" structs:"project_id,omitempty"`
	Resource       Resource `xorm:"-" structs:"-"`
	StartDateReset bool     `xorm:"-" structs:"-"`
	EndDateReset   bool     `xorm:"-" structs:"-"`
}

var (
	TaskStatusBacklog  int            = 1
	TaskStatusTodo     int            = 2
	TaskStatusProgress int            = 3
	TaskStatusDone     int            = 4
	TaskStatusNameMap  map[int]string = map[int]string{
		TaskStatusBacklog:  "backlog",
		TaskStatusTodo:     "todo",
		TaskStatusProgress: "progress",
		TaskStatusDone:     "done",
	}
)

type TaskJoin struct {
	Task            `xorm:"extends" json:"-"`
	Resource        `xorm:"extends" json:"-"`
	MapAreaResource `xorm:"extends" json:"-"`
}

func (m *Task) TableName() string {
	return "task"
}

func (m *Task) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	return
}

func (m *Task) SlicePtr() interface{} {
	ret := make([]Task, 0)
	return &ret
}

func (m *Task) Join() es.JoinGeneral {
	ret := TaskJoin{}
	return &ret
}

func (j *TaskJoin) Links() []es.JoinLinks {
	links := make([]es.JoinLinks, 0)
	resourceLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Resource.TableName(),
		LeftTable:  j.Resource.TableName(),
		LeftField:  "id",
		RightTable: j.Task.TableName(),
		RightField: "resource_id",
	}
	mapAreaResourceLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.MapAreaResource.TableName(),
		LeftTable:  j.MapAreaResource.TableName(),
		LeftField:  "resource_id",
		RightTable: j.Resource.TableName(),
		RightField: "id",
	}
	links = append(links, resourceLink)
	links = append(links, mapAreaResourceLink)
	return links
}

func (j *TaskJoin) SlicePtr() interface{} {
	ret := make([]TaskJoin, 0)
	return &ret
}

func (j *TaskJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.Task
	ret.Resource = join.Resource
	ret.Resource.Area.Id = join.MapAreaResource.AreaId
	return &ret
}

func (j *TaskJoin) TransferCopy(modelPtr es.ModelGeneral) {
	taskPtr := modelPtr.(*Task)
	(*taskPtr) = (*j).Task
	(*taskPtr).Resource = (*j).Resource
	(*taskPtr).Resource.Area.Id = (*j).MapAreaResource.AreaId
	return
}

func (j *TaskJoin) TransferSlicePtr(slicePtr interface{}) interface{} {
	joinSlicePtr := slicePtr.(*[]TaskJoin)
	joinSlice := *joinSlicePtr
	tasks := make([]Task, 0)
	for _, one := range joinSlice {
		taskPtr := (&one).Transfer().(*Task)
		tasks = append(tasks, *taskPtr)
	}
	return &tasks
}
