package models

import "time"

type SyncfusionGantt struct {
	Id        int
	Name      string
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
