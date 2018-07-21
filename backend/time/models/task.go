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

	Resource       *Resource    `xorm:"-" structs:"-" json:"Resource,omitempty"`
	Project        *Project     `xorm:"-" structs:"-" json:"Project,omitempty"`
	QuestTarget    *QuestTarget `xorm:"-" structs:"-" json:"QuestTarget,omitempty"`
	Area           *Area        `xorm:"-" structs:"-" json:"Area,omitempty"`
	StartDateReset bool         `xorm:"-" structs:"-"`
	EndDateReset   bool         `xorm:"-" structs:"-"`
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
	Task        `xorm:"extends" json:"-"`
	Resource    `xorm:"extends" json:"-"`
	Project     `xorm:"extends" json:"-"`
	QuestTarget `xorm:"extends" json:"-"`
	Area        `xorm:"extends" json:"-"`
}

func NewTask() *Task {
	ret := Task{
		Resource:    NewResource(),
		Area:        NewArea(),
		Project:     NewProject(),
		QuestTarget: NewQuestTarget(),
	}
	return &ret
}

func (m *Task) TableName() string {
	return "task"
}

func (m *Task) BuildCondition(session *xorm.Session) {
	keyPrefix := m.TableName() + "."
	params := structs.Map(m)
	condition := m.Model.BuildCondition(params, keyPrefix)
	session.Where(condition)
	m.QuestTarget.BuildCondition(session)
	m.Project.BuildCondition(session)
	m.Area.BuildCondition(session)
	m.Resource.BuildCondition(session)
	return
}

func (m *Task) SlicePtr() interface{} {
	ret := make([]Task, 0)
	return &ret
}

func (m *Task) Transfer(slicePtr interface{}) *[]Task {
	ret := slicePtr.(*[]Task)
	return ret
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
	projectLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Project.TableName(),
		LeftTable:  j.Project.TableName(),
		LeftField:  "id",
		RightTable: j.Task.TableName(),
		RightField: "project_id",
	}
	questTargetLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.QuestTarget.TableName(),
		LeftTable:  j.QuestTarget.TableName(),
		LeftField:  "id",
		RightTable: j.Project.TableName(),
		RightField: "quest_target_id",
	}
	areaLink := es.JoinLinks{
		Type:       es.InnerJoin,
		Table:      j.Area.TableName(),
		LeftTable:  j.Area.TableName(),
		LeftField:  "id",
		RightTable: j.QuestTarget.TableName(),
		RightField: "area_id",
	}
	links = append(links, resourceLink)
	links = append(links, projectLink)
	links = append(links, questTargetLink)
	links = append(links, areaLink)
	return links
}

func (j *TaskJoin) SlicePtr() interface{} {
	ret := make([]TaskJoin, 0)
	return &ret
}

func (j *TaskJoin) Transfer() es.ModelGeneral {
	join := *j
	ret := join.Task
	ret.Resource = &join.Resource
	ret.Project = &join.Project
	ret.QuestTarget = &join.QuestTarget
	ret.Area = &join.Area
	return &ret
}

func (j *TaskJoin) TransferCopy(modelPtr es.ModelGeneral) {
	taskPtr := modelPtr.(*Task)
	(*taskPtr) = (*j).Task
	(*taskPtr).Resource = &(*j).Resource
	(*taskPtr).Project = &(*j).Project
	(*taskPtr).QuestTarget = &(*j).QuestTarget
	(*taskPtr).Area = &(*j).Area
	return
}

func (j *TaskJoin) TransferCopySlice(slicePtr interface{}, targetPtr interface{}) {
	joinSlicePtr := slicePtr.(*[]TaskJoin)
	joinSlice := *joinSlicePtr
	tasks := make([]Task, 0)
	for _, one := range joinSlice {
		taskPtr := (&one).Transfer().(*Task)
		tasks = append(tasks, *taskPtr)
	}
	tasksPtr := targetPtr.(*[]Task)
	(*tasksPtr) = tasks
	return
}
