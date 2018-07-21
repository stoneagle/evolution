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
	Status    int
	Color     string
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
