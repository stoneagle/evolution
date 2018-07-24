package services

import (
	"errors"
	"math"
	"strconv"

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
	questsPtr, err := s.Pack.QuestTeamSvc.GetQuestsByUser(userId, models.QuestStatusExec)
	if err != nil {
		return
	}
	questIds := make([]int, 0)
	for _, one := range *questsPtr {
		questIds = append(questIds, one.Id)
	}

	fieldsMap, err := s.Pack.FieldSvc.Map()
	if err != nil {
		return
	}

	kanbans = make([]models.SyncfusionKanban, 0)
	task := models.NewTask()
	task.Project.Status = models.ProjectStatusWait
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
	fieldsMap, err := s.Pack.FieldSvc.Map()
	if err != nil {
		return
	}

	hour, _ := time.ParseDuration("-4h")
	for _, one := range *actionsPtr {
		schedule := models.SyncfusionSchedule{}
		schedule.Id = one.Id
		schedule.Name = one.Name
		fieldIdStr := strconv.Itoa(one.Area.FieldId)
		schedule.FieldId = fieldIdStr
		field, ok := fieldsMap[one.Area.FieldId]
		if !ok {
			err = errors.New("field not exist")
			return
		}
		schedule.Field = field
		schedule.StartDate = one.StartDate.Add(hour)
		schedule.EndDate = one.EndDate.Add(hour)
		schedule.AllDay = false
		schedule.Recurrence = false
		schedule.Area = *one.Area
		schedule.Resource = *one.Resource
		schedule.Task = *one.Task
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

func (s *Syncfusion) ListGantt(userId int, ganttLevel, ganttStatus string) (gantts []models.SyncfusionGantt, err error) {
	questsPtr, err := s.Pack.QuestTeamSvc.GetQuestsByUser(userId, models.QuestStatusExec)
	if err != nil {
		return
	}
	questIds := make([]int, 0)
	for _, one := range *questsPtr {
		questIds = append(questIds, one.Id)
	}

	project := models.NewProject()
	switch ganttStatus {
	case models.SyncfusionGanttStatusFinish:
		project.Status = models.ProjectStatusFinish
	case models.SyncfusionGanttStatusWait:
		project.Status = models.ProjectStatusWait
	default:
		err = errors.New(fmt.Sprintf("gantt status not match:%v", ganttStatus))
		return
	}
	project.QuestTarget.QuestIds = questIds
	projectsGeneralPtr := project.SlicePtr()
	err = s.Pack.ProjectSvc.List(project, projectsGeneralPtr)
	if err != nil {
		return
	}
	projectsPtr := project.Transfer(projectsGeneralPtr)

	task := models.NewTask()
	switch ganttStatus {
	case models.SyncfusionGanttStatusFinish:
		project.Status = models.ProjectStatusFinish
	case models.SyncfusionGanttStatusWait:
		project.Status = models.ProjectStatusWait
	default:
		err = errors.New(fmt.Sprintf("gantt status not match:%v", ganttStatus))
		return
	}
	task.QuestTarget.QuestIds = questIds
	tasksGeneralPtr := task.SlicePtr()
	err = s.Pack.TaskSvc.List(task, tasksGeneralPtr)
	if err != nil {
		return
	}
	tasksPtr := task.Transfer(tasksGeneralPtr)

	fieldsMap, err := s.Pack.FieldSvc.Map()
	if err != nil {
		return
	}

	switch ganttLevel {
	case models.SyncfusionGanttLevelQuest:
		projectsMap := map[int][]models.SyncfusionGantt{}
		gantts, err = s.buildQuestSlice(questsPtr, projectsMap)
		if err != nil {
			return gantts, err
		}
	case models.SyncfusionGanttLevelProject:
		tasksMap := map[int][]models.SyncfusionGantt{}
		projectsMap, err := s.buildProjectMap(projectsPtr, tasksMap, fieldsMap)
		if err != nil {
			return gantts, err
		}
		gantts, err = s.buildQuestSlice(questsPtr, projectsMap)
		if err != nil {
			return gantts, err
		}
	case models.SyncfusionGanttLevelTask:
		tasksMap, err := s.buildTaskMap(tasksPtr)
		if err != nil {
			return gantts, err
		}
		projectsMap, err := s.buildProjectMap(projectsPtr, tasksMap, fieldsMap)
		if err != nil {
			return gantts, err
		}
		gantts, err = s.buildQuestSlice(questsPtr, projectsMap)
		if err != nil {
			return gantts, err
		}
	default:
		err = errors.New(fmt.Sprintf("gantt level not match: %v", ganttLevel))
		return
	}
	return
}

func (s *Syncfusion) buildTaskMap(tasksPtr *[]models.Task) (result map[int][]models.SyncfusionGantt, err error) {
	result = make(map[int][]models.SyncfusionGantt)
	for _, one := range *tasksPtr {
		if _, ok := result[one.ProjectId]; !ok {
			result[one.ProjectId] = make([]models.SyncfusionGantt, 0)
		}
		err = one.Hook()
		if err != nil {
			return
		}
		err = one.Project.Hook()
		if err != nil {
			return
		}
		gantt := models.SyncfusionGantt{}
		gantt.Id = one.UuidNumber
		gantt.Parent = one.Project.UuidNumber
		gantt.Relate = one.Resource.Name
		gantt.RelateId = one.Id
		gantt.Name = one.Name
		gantt.StartDate = one.StartDate
		gantt.EndDate = one.EndDate
		if one.Status == models.TaskStatusDone {
			gantt.Progress = 100
		} else {
			gantt.Progress = 0
		}
		diffHours := gantt.EndDate.Sub(gantt.StartDate).Hours()
		gantt.Duration = int(math.Ceil(diffHours / 24))
		gantt.Status = one.Status
		gantt.Expanded = false
		child := make([]models.SyncfusionGantt, 0)
		gantt.Children = child
		result[one.ProjectId] = append(result[one.ProjectId], gantt)
	}
	return
}

func (s *Syncfusion) buildProjectMap(projectsPtr *[]models.Project, tasksMap map[int][]models.SyncfusionGantt, fieldsMap map[int]models.Field) (result map[int][]models.SyncfusionGantt, err error) {
	result = make(map[int][]models.SyncfusionGantt)
	for _, one := range *projectsPtr {
		if _, ok := result[one.QuestTarget.QuestId]; !ok {
			result[one.QuestTarget.QuestId] = make([]models.SyncfusionGantt, 0)
		}
		err = one.Hook()
		if err != nil {
			return
		}
		gantt := models.SyncfusionGantt{}
		field, ok := fieldsMap[one.Area.FieldId]
		if ok {
			gantt.Color = field.Color
		}
		gantt.Id = one.UuidNumber
		// still not set Parent, need set in Quest levetl
		gantt.Relate = one.Area.Name
		gantt.RelateId = one.Id
		gantt.Name = one.Name
		gantt.StartDate = one.StartDate
		gantt.Status = one.Status
		gantt.Expanded = false
		gantt.Progress = 0
		if children, ok := tasksMap[one.Id]; ok {
			gantt.Children = children
			var maxEndDate time.Time
			for _, task := range children {
				if task.EndDate.After(maxEndDate) {
					maxEndDate = task.EndDate
				}
			}
			gantt.EndDate = maxEndDate
			diffHours := gantt.EndDate.Sub(gantt.StartDate).Hours()
			gantt.Duration = int(math.Ceil(diffHours / 24))
		} else {
			child := make([]models.SyncfusionGantt, 0)
			gantt.Children = child
			gantt.EndDate = one.StartDate
			gantt.Duration = 0
		}
		result[one.QuestTarget.QuestId] = append(result[one.QuestTarget.QuestId], gantt)
	}
	return
}

func (s *Syncfusion) buildQuestSlice(questsPtr *[]models.Quest, projectsMap map[int][]models.SyncfusionGantt) (result []models.SyncfusionGantt, err error) {
	result = make([]models.SyncfusionGantt, 0)
	for _, one := range *questsPtr {
		err = one.Hook()
		if err != nil {
			return
		}
		gantt := models.SyncfusionGantt{}
		gantt.Id = one.UuidNumber
		gantt.Parent = 0
		gantt.RelateId = one.Id
		gantt.Name = one.Name
		gantt.Progress = 0
		gantt.Expanded = false
		if children, ok := projectsMap[one.Id]; ok {
			minStartDate := time.Now()
			var maxEndDate time.Time
			for k, project := range children {
				if project.StartDate.Before(minStartDate) {
					minStartDate = project.StartDate
				}
				if project.EndDate.After(maxEndDate) {
					maxEndDate = project.EndDate
				}
				children[k].Parent = one.UuidNumber
			}
			gantt.Children = children
			gantt.StartDate = minStartDate
			gantt.EndDate = maxEndDate
			diffHours := gantt.EndDate.Sub(gantt.StartDate).Hours()
			gantt.Duration = int(math.Ceil(diffHours / 24))
		} else {
			child := make([]models.SyncfusionGantt, 0)
			gantt.Children = child
			gantt.StartDate = one.StartDate
			gantt.EndDate = one.EndDate
			gantt.Duration = 0
		}
		result = append(result, gantt)
	}
	return
}
