package services

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"

	"evolution/backend/common/logger"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"fmt"
	"time"
)

type Syncfusion struct {
	Base
	structs.Service
}

func NewSyncfusion(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Syncfusion {
	ret := Syncfusion{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Syncfusion) ListKanban(userId int) (kanbans []models.SyncfusionKanban, err error) {
	questTeam := models.QuestTeam{}
	questTeam.UserId = userId
	questTeam.Quest.Status = models.QuestStatusExec
	questTeams, err := s.QuestTeamSvc.ListWithCondition(&questTeam)
	if err != nil {
		return
	}

	questIds := make([]int, 0)
	for _, one := range questTeams {
		questIds = append(questIds, one.QuestId)
	}
	project := models.Project{}
	project.QuestTarget.QuestIds = questIds
	projects, err := s.ProjectSvc.ListWithCondition(&project)
	if err != nil {
		return
	}

	projectIds := make([]int, 0)
	areaIds := make([]int, 0)
	projectsMap := make(map[int]models.Project)
	for _, one := range projects {
		projectIds = append(projectIds, one.Id)
		projectsMap[one.Id] = one
		areaIds = append(areaIds, one.QuestTarget.AreaId)
	}
	task := models.Task{}
	task.ProjectIds = projectIds
	tasks, err := s.TaskSvc.ListWithCondition(&task)
	if err != nil {
		return
	}

	area := models.Area{}
	area.Ids = areaIds
	areas, err := s.AreaSvc.ListWithCondition(&area)
	if err != nil {
		return
	}
	areasMap := make(map[int]models.Area)
	for _, one := range areas {
		areasMap[one.Id] = one
	}

	fields, err := s.FieldSvc.List()
	if err != nil {
		return
	}
	fieldsMap := make(map[int]models.Field)
	for _, one := range fields {
		fieldsMap[one.Id] = one
	}

	kanbans = make([]models.SyncfusionKanban, 0)
	for _, one := range tasks {
		project, ok := projectsMap[one.ProjectId]
		if !ok {
			err = errors.New(fmt.Sprintf("task %s not find its project", one.Name))
			return
		}
		area, ok := areasMap[one.Resource.Area.Id]
		if !ok {
			err = errors.New(fmt.Sprintf("task %s not find its area", one.Name))
			return
		}
		field, ok := fieldsMap[area.FieldId]
		if !ok {
			err = errors.New(fmt.Sprintf("task %s not find its field", one.Name))
			return
		}
		kanban := models.SyncfusionKanban{}
		kanban.Id = one.Id
		kanban.Name = one.Name
		kanban.Desc = one.Desc
		kanban.Status = one.Status
		kanban.StatusName, ok = models.TaskStatusNameMap[one.Status]
		if !ok {
			err = errors.New(fmt.Sprintf("task %s not find its status name", one.Name))
			return
		}
		kanban.Tags = field.Name + "," + area.Name + "," + one.Resource.Name
		kanban.ResourceId = one.Resource.Id
		kanban.ResourceName = one.Resource.Name
		kanban.ProjectId = project.Id
		kanban.ProjectName = project.Name
		kanban.FieldId = field.Id
		kanban.FieldName = field.Name
		kanbans = append(kanbans, kanban)
	}
	return kanbans, err
}

func (s *Syncfusion) ListSchedule(userId int, startDate, endDate time.Time) (schedules []models.SyncfusionSchedule, err error) {
	action := models.Action{}
	action.UserId = userId
	action.StartDate = startDate
	action.EndDate = endDate
	actions, err := s.ActionSvc.ListWithCondition(&action)
	if err != nil {
		return
	}
	schedules = make([]models.SyncfusionSchedule, 0)
	hour, _ := time.ParseDuration("-4h")
	for _, one := range actions {
		schedule := models.SyncfusionSchedule{}
		schedule.Id = one.Id
		schedule.Name = one.Name
		schedule.StartDate = one.StartDate.Add(hour)
		schedule.EndDate = one.EndDate.Add(hour)
		schedule.AllDay = false
		schedule.Recurrence = false
		schedules = append(schedules, schedule)
	}
	return
}

func (s *Syncfusion) ListTreeGrid(fieldId, parentId int) (treeGrids []models.SyncfusionTreeGrid, err error) {
	area := models.Area{}
	area.FieldId = fieldId
	if parentId == 0 {
		area.Type = models.AreaTypeRoot
	} else {
		area.Parent = parentId
	}

	areas, err := s.AreaSvc.ListWithCondition(&area)
	if err != nil {
		return
	}

	treeGrids = make([]models.SyncfusionTreeGrid, 0)
	for _, one := range areas {
		treeGrid := models.SyncfusionTreeGrid{}
		treeGrid.Id = one.Id
		treeGrid.Name = one.Name
		if parentId == 0 {
			treeGrid.Parent = nil
		} else {
			treeGrid.Parent = one.Parent
		}

		if one.Type == models.AreaTypeLeaf {
			treeGrid.IsParent = false
			treeGrid.IsExpanded = true
		} else {
			treeGrid.IsParent = true
			treeGrid.IsExpanded = false
			children := make([]models.SyncfusionTreeGrid, 0)
			children = append(children, models.SyncfusionTreeGrid{})
			treeGrid.Children = children
		}
		treeGrids = append(treeGrids, treeGrid)
	}
	return
}

func (s *Syncfusion) ListGantt(userId int) (gantts []models.SyncfusionGantt, err error) {
	questTeam := models.QuestTeam{}
	questTeam.UserId = userId
	questTeam.Quest.Status = models.QuestStatusExec
	questTeams, err := s.QuestTeamSvc.ListWithCondition(&questTeam)
	if err != nil {
		return
	}

	questIds := make([]int, 0)
	quests := make([]models.Quest, 0)
	for _, one := range questTeams {
		questIds = append(questIds, one.QuestId)
		quests = append(quests, one.Quest)
	}

	project := models.Project{}
	project.QuestTarget.QuestIds = questIds
	projects, err := s.ProjectSvc.ListWithCondition(&project)
	if err != nil {
		return
	}

	projectIds := make([]int, 0)
	for _, one := range projects {
		projectIds = append(projectIds, one.Id)
	}
	task := models.Task{}
	task.ProjectIds = projectIds
	tasks, err := s.TaskSvc.ListWithCondition(&task)
	if err != nil {
		return
	}

	tasksMap := s.buildTaskMap(tasks)
	projectsMap := s.buildProjectMap(projects, tasksMap)
	gantts = s.buildQuestSlice(quests, projectsMap)
	return
}

func (s *Syncfusion) buildTaskMap(tasks []models.Task) map[int][]models.SyncfusionGantt {
	result := make(map[int][]models.SyncfusionGantt)
	for _, one := range tasks {
		if _, ok := result[one.ProjectId]; !ok {
			result[one.ProjectId] = make([]models.SyncfusionGantt, 0)
		}
		gantt := models.SyncfusionGantt{}
		gantt.Id = one.Id
		gantt.Parent = one.ProjectId
		gantt.Relate = one.Resource.Name
		gantt.Name = one.Name
		gantt.StartDate = one.StartDate
		gantt.EndDate = one.EndDate
		gantt.Progress = 0
		gantt.Duration = 0
		gantt.Expanded = false
		child := make([]models.SyncfusionGantt, 0)
		gantt.Children = child
		result[one.ProjectId] = append(result[one.ProjectId], gantt)
	}
	return result
}

func (s *Syncfusion) buildProjectMap(projects []models.Project, tasksMap map[int][]models.SyncfusionGantt) map[int][]models.SyncfusionGantt {
	result := make(map[int][]models.SyncfusionGantt)
	for _, one := range projects {
		if _, ok := result[one.QuestTarget.QuestId]; !ok {
			result[one.QuestTarget.QuestId] = make([]models.SyncfusionGantt, 0)
		}
		gantt := models.SyncfusionGantt{}
		gantt.Id = one.Id
		gantt.Parent = one.QuestTarget.QuestId
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
			child := make([]models.SyncfusionGantt, 0)
			gantt.Children = child
		}
		result[one.QuestTarget.QuestId] = append(result[one.QuestTarget.QuestId], gantt)
	}
	return result
}

func (s *Syncfusion) buildQuestSlice(quests []models.Quest, projectsMap map[int][]models.SyncfusionGantt) []models.SyncfusionGantt {
	result := make([]models.SyncfusionGantt, 0)
	for _, one := range quests {
		gantt := models.SyncfusionGantt{}
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
