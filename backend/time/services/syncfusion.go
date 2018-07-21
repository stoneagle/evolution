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
	Pack ServicePackage
	structs.Service
}

func NewSyncfusion(engine *xorm.Engine, cache *redis.Client, log *logger.Logger) *Syncfusion {
	ret := Syncfusion{}
	ret.Init(engine, cache, log)
	return &ret
}

func (s *Syncfusion) ListKanban(userId int) (kanbans []models.SyncfusionKanban, err error) {
	questTeam := models.NewQuestTeam()
	questTeam.UserId = userId
	questTeam.Quest.Status = models.QuestStatusExec
	questTeamsGeneralPtr := questTeam.SlicePtr()
	err = s.Pack.QuestTeamSvc.List(questTeam, questTeamsGeneralPtr)
	if err != nil {
		return
	}
	questTeamsPtr := questTeam.Transfer(questTeamsGeneralPtr)
	questIds := make([]int, 0)
	for _, one := range *questTeamsPtr {
		questIds = append(questIds, one.QuestId)
	}

	field := models.NewField()
	fieldsGeneralPtr := field.SlicePtr()
	err = s.Pack.FieldSvc.List(field, fieldsGeneralPtr)
	if err != nil {
		return
	}
	fieldsPtr := field.Transfer(fieldsGeneralPtr)
	fieldsMap := make(map[int]models.Field)
	for _, one := range *fieldsPtr {
		fieldsMap[one.Id] = one
	}

	kanbans = make([]models.SyncfusionKanban, 0)
	task := models.NewTask()
	task.QuestTarget.Status = models.QuestTargetStatusWait
	task.QuestTarget.QuestIds = questIds
	tasksGeneralPtr := task.SlicePtr()
	err = s.Pack.TaskSvc.List(task, tasksGeneralPtr)
	if err != nil {
		return
	}
	tasksPtr := task.Transfer(tasksGeneralPtr)
	for _, one := range *tasksPtr {
		field, ok := fieldsMap[one.Area.FieldId]
		if !ok {
			err = errors.New(fmt.Sprintf("task %s not find its field", one.Name))
			s.Logger.Log(logger.WarnLevel, one, err)
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
		kanban.Tags = field.Name + "," + one.Area.Name + "," + one.Resource.Name
		kanban.ResourceId = one.Resource.Id
		kanban.ResourceName = one.Resource.Name
		kanban.ProjectId = one.Project.Id
		kanban.ProjectName = one.Project.Name
		kanban.FieldId = field.Id
		kanban.FieldName = field.Name
		kanbans = append(kanbans, kanban)
	}

	return kanbans, err
}

func (s *Syncfusion) ListSchedule(userId int, startDate, endDate time.Time) (schedules []models.SyncfusionSchedule, err error) {
	action := models.NewAction()
	action.UserId = userId
	action.StartDate = startDate
	action.EndDate = endDate
	actionsGeneralPtr := action.SlicePtr()
	err = s.Pack.ActionSvc.List(action, actionsGeneralPtr)
	if err != nil {
		return
	}
	actionsPtr := action.Transfer(actionsGeneralPtr)

	schedules = make([]models.SyncfusionSchedule, 0)
	hour, _ := time.ParseDuration("-4h")
	for _, one := range *actionsPtr {
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
	area := models.NewArea()
	area.FieldId = fieldId
	if parentId == 0 {
		area.Type = models.AreaTypeRoot
	} else {
		area.Parent = parentId
	}

	areasGeneralPtr := area.SlicePtr()
	err = s.Pack.AreaSvc.List(area, areasGeneralPtr)
	if err != nil {
		return
	}
	areasPtr := area.Transfer(areasGeneralPtr)

	treeGrids = make([]models.SyncfusionTreeGrid, 0)
	for _, one := range *areasPtr {
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
	questTeam := models.NewQuestTeam()
	questTeam.UserId = userId
	questTeam.Quest.Status = models.QuestStatusExec
	questTeamsGeneralPtr := questTeam.SlicePtr()
	err = s.Pack.QuestTeamSvc.List(questTeam, questTeamsGeneralPtr)
	if err != nil {
		return
	}
	questTeamsPtr := questTeam.Transfer(questTeamsGeneralPtr)

	questIds := make([]int, 0)
	quests := make([]models.Quest, 0)
	for _, one := range *questTeamsPtr {
		questIds = append(questIds, one.QuestId)
		quests = append(quests, *one.Quest)
	}

	task := models.NewTask()
	task.QuestTarget.Ids = questIds
	tasksGeneralPtr := task.SlicePtr()
	err = s.Pack.TaskSvc.List(task, tasksGeneralPtr)
	if err != nil {
		return
	}
	tasksPtr := task.Transfer(tasksGeneralPtr)

	projectFilterMap := map[int]bool{}
	projects := make([]models.Project, 0)
	for k, _ := range *tasksPtr {
		task := (*tasksPtr)[k]
		project := models.NewProject()
		*(project) = *(task.Project)
		project.QuestTarget = models.NewQuestTarget()
		project.Area = models.NewArea()
		*(project.QuestTarget) = *(task.QuestTarget)
		*(project.Area) = *(task.Area)
		_, ok := projectFilterMap[project.Id]
		if ok {
			continue
		}
		projectFilterMap[project.Id] = true
		projects = append(projects, *project)
	}

	tasksMap := s.buildTaskMap(tasksPtr)
	projectsMap := s.buildProjectMap(&projects, tasksMap)
	gantts = s.buildQuestSlice(quests, projectsMap)
	return
}

func (s *Syncfusion) buildTaskMap(tasksPtr *[]models.Task) map[int][]models.SyncfusionGantt {
	result := make(map[int][]models.SyncfusionGantt)
	for _, one := range *tasksPtr {
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

func (s *Syncfusion) buildProjectMap(projectsPtr *[]models.Project, tasksMap map[int][]models.SyncfusionGantt) map[int][]models.SyncfusionGantt {
	result := make(map[int][]models.SyncfusionGantt)
	for _, one := range *projectsPtr {
		if _, ok := result[one.QuestTarget.QuestId]; !ok {
			result[one.QuestTarget.QuestId] = make([]models.SyncfusionGantt, 0)
		}
		gantt := models.SyncfusionGantt{}
		gantt.Id = one.Id
		gantt.Parent = one.QuestTarget.QuestId
		gantt.Relate = one.Area.Name
		gantt.Name = one.Name
		gantt.StartDate = one.StartDate
		gantt.Progress = 0
		gantt.Duration = 0
		gantt.Expanded = false
		if children, ok := tasksMap[one.Id]; ok {
			gantt.Children = children
			var maxEndDate time.Time
			for _, task := range children {
				if task.EndDate.After(maxEndDate) {
					maxEndDate = task.EndDate
				}
			}
			gantt.EndDate = maxEndDate
		} else {
			child := make([]models.SyncfusionGantt, 0)
			gantt.EndDate = one.StartDate
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
