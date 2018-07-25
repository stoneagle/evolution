package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	uuid "github.com/satori/go.uuid"
)

type Project struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	Text      string    `xorm:"not null default '' comment('内容') VARCHAR(255)"`
	StartDate time.Time `xorm:"not null comment('开始日期') DATETIME"`
	Duration  int       `xorm:"not null default 0 comment('持续时间') INT(11)"`
	Progress  float32   `xorm:"not null default 0 comment('进度') FLOAT"`
	UserId    int       `xorm:"not null default 0 comment('所属用户') INT(11)"`
	TargetId  int       `xorm:"not null default 0 comment('所属目标') INT(11)"`
	Del       int       `xorm:"not null default 0 comment('软删除') TINYINT(1)"`
	Ctime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

type ProjectJoin struct {
	Project `xorm:"extends"`
	Target  `xorm:"extends"`
}

func (c *Project) Transfer(src, des *xorm.Engine, userId int) {
	session := des.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Printf("project to action session error:%v\r\n", err.Error())
		return
	}

	oldProjectJoins := make([]ProjectJoin, 0)
	err = src.Table("project").Join("LEFT", "target", "target.id = project.target_id").Find(&oldProjectJoins)
	if err != nil {
		fmt.Printf("project join get error:%v\r\n", err.Error())
		return
	}

	projectNum := 0
	taskNum := 0
	actionNum := 0
	for _, oldProject := range oldProjectJoins {
		newProject := models.NewProject()
		quest := models.NewQuest()
		_, err = session.Where("name = ?", oldProject.Target.Name).Get(quest)
		if err != nil {
			session.Rollback()
			fmt.Printf("new quest get error:%v\r\n", err.Error())
			return
		}
		oldTasks := make([]Task, 0)
		newAreaTasksMap := map[int][]models.Task{}
		src.Where("parent = ?", oldProject.Project.Id).Find(&oldTasks)
		// get old task to merge area
		for _, oldTask := range oldTasks {
			newResourceJoin, err := getResourceJoin(src, des, oldTask.EntityId, oldProject.Target.FieldId)
			if err != nil {
				fmt.Printf("new resource join get error %v\r\n", err.Error())
				session.Rollback()
				return
			}
			newTask := models.NewTask()
			newTask.ResourceId = newResourceJoin.Resource.Id
			newTask.UserId = userId
			newTask.StartDate = oldTask.StartDate
			newTask.Name = oldTask.Text
			if oldTask.Progress == 0 {
				newTask.Status = models.TaskStatusProgress
			} else {
				newTask.Status = models.TaskStatusDone
			}
			if oldTask.Del != 0 {
				newTask.DeletedAt = oldTask.Utime
			}
			newTask.CreatedAt = oldTask.Ctime
			newTask.UpdatedAt = oldTask.Utime
			tasksSlice, ok := newAreaTasksMap[newResourceJoin.Area.Id]
			if !ok {
				newAreaTasksMap[newResourceJoin.Area.Id] = make([]models.Task, 0)
			}
			tasksSlice = append(tasksSlice, *newTask)
			newAreaTasksMap[newResourceJoin.Area.Id] = tasksSlice
		}
		newProject.QuestTarget.QuestId = quest.Id
		newProject.StartDate = oldProject.Project.StartDate
		if oldProject.Project.Del != 0 {
			newProject.DeletedAt = oldProject.Project.Utime
		}
		newProject.CreatedAt = oldProject.Project.Ctime
		newProject.UpdatedAt = oldProject.Project.Utime
		newProject.Status = models.ProjectStatusWait

		// base on merge area to insert project and task
		for areaId, tasksSlice := range newAreaTasksMap {
			newProject.Id = 0
			newProject.QuestTarget.AreaId = areaId
			// get relate quest target
			relateQuestTarget := models.NewQuestTarget()
			has, err := des.Where("quest_id = ?", newProject.QuestTarget.QuestId).And("area_id = ?", newProject.QuestTarget.AreaId).Get(relateQuestTarget)
			if !has {
				fmt.Printf("quest target not exist quest_id:%v, area_id:%v\r\n", newProject.QuestTarget.QuestId, newProject.QuestTarget.AreaId)
				session.Rollback()
				return
			}
			if err != nil {
				fmt.Printf("quest target get error %+v\r\n", err.Error())
				session.Rollback()
				return
			}
			newProject.QuestTargetId = relateQuestTarget.Id

			newArea := models.NewArea()
			_, err = session.Id(areaId).Get(newArea)
			if err != nil {
				fmt.Printf("new area get error %v\r\n", err.Error())
				session.Rollback()
				return
			}
			newProject.Name = oldProject.Project.Text + ":" + newArea.Name
			u := uuid.NewV4()
			newProject.Uuid = u.String()
			_, err = session.Insert(newProject)
			if err != nil {
				fmt.Printf("new project insert error %v\r\n", err.Error())
				session.Rollback()
				return
			}
			for _, newTask := range tasksSlice {
				newTask.ProjectId = newProject.Id
				u := uuid.NewV4()
				newTask.Uuid = u.String()
				_, err = session.Insert(&newTask)
				if err != nil {
					fmt.Printf("new task insert error %v\r\n", err.Error())
					session.Rollback()
					return
				}
				taskNum++

				// get old task and its relate action
				needOldTask := Task{}
				has, err := src.Where("text = ?", newTask.Name).And("start_date = ?", newTask.StartDate).Get(&needOldTask)
				if err != nil || !has {
					fmt.Printf("old relate task get error %v\r\n", err.Error())
					session.Rollback()
					return
				}
				oldActions := make([]Action, 0)
				src.Where("task_id = ?", needOldTask.Id).Find(&oldActions)
				for _, oldAction := range oldActions {
					newAction := models.NewAction()
					newAction.TaskId = newTask.Id
					newAction.UserId = userId
					newAction.Name = oldAction.Text
					newAction.Time = oldAction.ExecTime / 60
					newAction.StartDate = oldAction.StartDate
					newAction.EndDate = oldAction.EndDate
					newAction.CreatedAt = oldAction.Ctime
					newAction.UpdatedAt = oldAction.Utime
					_, err = session.Insert(newAction)
					if err != nil {
						fmt.Printf("new action insert error %v\r\n", err.Error())
						session.Rollback()
						return
					}
					actionNum++
				}
			}
			projectNum++
		}
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("project to action session commit error:%v\r\n", err)
	}
	fmt.Printf("project transfer success:%v\r\n", projectNum)
	fmt.Printf("task transfer success:%v\r\n", taskNum)
	fmt.Printf("action transfer success:%v\r\n", actionNum)
}
