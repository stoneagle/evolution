package models

import (
	"time"
)

type SyncfusionGantt struct {
	Id        int
	Name      string
	Relate    string
	Parent    int
	Progress  int
	Duration  int
	Expanded  bool
	StartDate time.Time
	EndDate   time.Time
	Children  []SyncfusionGantt
}

type SyncfusionTreeGrid struct {
	Id         int
	Name       string
	Parent     interface{}
	IsParent   bool
	IsExpanded bool
	Children   []SyncfusionTreeGrid `json:"Children,omitempty"`
}

type SyncfusionSchedule struct {
	Id         int
	Name       string
	StartDate  time.Time
	EndDate    time.Time
	AllDay     bool
	Recurrence bool
}

type SyncfusionKanban struct {
	Id           int
	Name         string
	Desc         string
	Tags         string
	Status       int
	StatusName   string
	ResourceId   int
	ResourceName string
	ProjectId    int
	ProjectName  string
	FieldId      int
	FieldName    string
}

var (
	SyncfusionScheduleViewWeek     string = "week"
	SyncfusionScheduleViewDay      string = "day"
	SyncfusionScheduleViewMonth    string = "month"
	SyncfusionScheduleViewAgenda   string = "agenda"
	SyncfusionScheduleViewWorkWeek string = "workweek"
)

func (m *SyncfusionGantt) BuildTaskMap(tasks []Task) map[int][]SyncfusionGantt {
	result := make(map[int][]SyncfusionGantt)
	for _, one := range tasks {
		if _, ok := result[one.ProjectId]; !ok {
			result[one.ProjectId] = make([]SyncfusionGantt, 0)
		}
		gantt := SyncfusionGantt{}
		gantt.Id = one.Id
		gantt.Parent = one.ProjectId
		gantt.Relate = one.Resource.Name
		gantt.Name = one.Name
		gantt.StartDate = one.StartDate
		gantt.EndDate = one.EndDate
		gantt.Progress = 0
		gantt.Duration = 0
		gantt.Expanded = false
		child := make([]SyncfusionGantt, 0)
		gantt.Children = child
		result[one.ProjectId] = append(result[one.ProjectId], gantt)
	}
	return result
}

func (m *SyncfusionGantt) BuildProjectMap(projects []Project, tasksMap map[int][]SyncfusionGantt) map[int][]SyncfusionGantt {
	result := make(map[int][]SyncfusionGantt)
	for _, one := range projects {
		if _, ok := result[one.QuestId]; !ok {
			result[one.QuestId] = make([]SyncfusionGantt, 0)
		}
		gantt := SyncfusionGantt{}
		gantt.Id = one.Id
		gantt.Parent = one.QuestId
		gantt.Relate = one.Area.Name
		gantt.Name = one.Name
		gantt.StartDate = one.StartDate
		gantt.EndDate = one.StartDate
		gantt.Progress = 0
		gantt.Duration = 0
		gantt.Expanded = false
		if children, ok := tasksMap[one.Id]; ok {
			gantt.Children = children
		} else {
			child := make([]SyncfusionGantt, 0)
			gantt.Children = child
		}
		result[one.QuestId] = append(result[one.QuestId], gantt)
	}
	return result
}

func (m *SyncfusionGantt) BuildQuestSlice(quests []Quest, projectsMap map[int][]SyncfusionGantt) []SyncfusionGantt {
	result := make([]SyncfusionGantt, 0)
	for _, one := range quests {
		gantt := SyncfusionGantt{}
		gantt.Id = one.Id
		gantt.Parent = 0
		gantt.Name = one.Name
		gantt.StartDate = one.StartDate
		gantt.EndDate = one.EndDate
		gantt.Progress = 0
		gantt.Duration = 0
		gantt.Expanded = false
		if children, ok := projectsMap[one.Id]; ok {
			gantt.Children = children
		}
		result = append(result, gantt)
	}

	return result
}

func (m *SyncfusionSchedule) BuildActionSlice(actions []Action) []SyncfusionSchedule {
	result := make([]SyncfusionSchedule, 0)
	hour, _ := time.ParseDuration("-4h")
	for _, one := range actions {
		schedule := SyncfusionSchedule{}
		schedule.Id = one.Id
		schedule.Name = one.Name
		schedule.StartDate = one.StartDate.Add(hour)
		schedule.EndDate = one.EndDate.Add(hour)
		schedule.AllDay = false
		schedule.Recurrence = false
		result = append(result, schedule)
	}
	return result
}
