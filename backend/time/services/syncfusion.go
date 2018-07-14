package services

import (
	"errors"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

type Syncfusion struct {
	structs.Service
}

func NewSyncfusion(engine *xorm.Engine, cache *redis.Client) *Syncfusion {
	ret := Syncfusion{}
	ret.Engine = engine
	ret.Cache = cache
	return &ret
}

func (s *Syncfusion) ListKanban(userId int) (kanbans []models.SyncfusionKanban, err error) {
	questTeam := models.QuestTeam{}
	questTeam.UserId = userId
	questTeam.Quest.Status = models.QuestStatusExec
	questTeams, err := NewQuestTeam(s.Engine, s.Cache).ListWithCondition(&questTeam)
	if err != nil {
		return
	}

	questIds := make([]int, 0)
	for _, one := range questTeams {
		questIds = append(questIds, one.QuestId)
	}
	project := models.Project{}
	project.QuestIds = questIds
	projects, err := NewProject(s.Engine, s.Cache).ListWithCondition(&project)
	if err != nil {
		return
	}

	projectIds := make([]int, 0)
	areaIds := make([]int, 0)
	projectsMap := make(map[int]models.Project)
	for _, one := range projects {
		projectIds = append(projectIds, one.Id)
		projectsMap[one.Id] = one
		areaIds = append(areaIds, one.AreaId)
	}
	task := models.Task{}
	task.ProjectIds = projectIds
	tasks, err := NewTask(s.Engine, s.Cache).ListWithCondition(&task)
	if err != nil {
		return
	}

	area := models.Area{}
	area.Ids = areaIds
	areas, err := NewArea(s.Engine, s.Cache).ListWithCondition(&area)
	if err != nil {
		return
	}
	areasMap := make(map[int]models.Area)
	for _, one := range areas {
		areasMap[one.Id] = one
	}

	fields, err := NewField(s.Engine, s.Cache).List()
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
		area, ok := areasMap[one.Resource.AreaId]
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
