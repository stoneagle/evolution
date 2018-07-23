package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/jinzhu/now"
	uuid "github.com/satori/go.uuid"
)

type Target struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	Name       string    `xorm:"not null default '' comment('名称') VARCHAR(256)"`
	Desc       string    `xorm:"not null default '' comment('描述') VARCHAR(256)"`
	UserId     int       `xorm:"not null default 0 comment('用户ID') INT(11)"`
	PriorityId int       `xorm:"not null default 0 comment('优先级') SMALLINT(6)"`
	FieldId    int       `xorm:"not null default 0 comment('所属范畴') SMALLINT(6)"`
	Ctime      time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime      time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

var (
	TargetPriorityImportantAndBusy int         = 1
	TargetPriorityImportant        int         = 2
	TargetPriorityBusy             int         = 3
	TargetPriorityNormal           int         = 4
	TargetPriorityMap              map[int]int = map[int]int{
		TargetPriorityImportantAndBusy: models.QuestConstraintImportantAndBusy,
		TargetPriorityImportant:        models.QuestConstraintImportant,
		TargetPriorityBusy:             models.QuestConstraintBusy,
		TargetPriorityNormal:           models.QuestConstraintNormal,
	}
)

func (c *Target) Transfer(src, des *xorm.Engine, userId int) {
	session := des.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Printf("target session error:%v\r\n", err.Error())
		return
	}

	oldTargets := make([]Target, 0)
	src.Find(&oldTargets)
	insertNum := 0
	for _, target := range oldTargets {
		quest := models.NewQuest()
		quest.Name = target.Name
		quest.FounderId = userId
		quest.Members = 1
		status, ok := TargetPriorityMap[target.PriorityId]
		if !ok {
			session.Rollback()
			fmt.Printf("target priority map not exist:%v\r\n", target.PriorityId)
			return
		}
		quest.Constraint = status
		quest.Status = models.QuestStatusExec
		quest.StartDate = target.Ctime
		quest.EndDate = now.New(time.Now()).EndOfYear()
		quest.CreatedAt = target.Ctime
		quest.UpdatedAt = target.Utime
		u := uuid.NewV4()
		quest.Uuid = u.String()
		_, err = session.Insert(quest)
		if err != nil {
			session.Rollback()
			fmt.Printf("quest insert error:%v\r\n", err.Error())
			return
		}
		questTeam := models.QuestTeam{}
		questTeam.StartDate = quest.StartDate
		questTeam.EndDate = quest.EndDate
		questTeam.UserId = userId
		questTeam.QuestId = quest.Id
		_, err = session.Insert(&questTeam)
		if err != nil {
			session.Rollback()
			fmt.Printf("quest team insert error:%v\r\n", err.Error())
			return
		}

		oldTargetEntities := make([]TargetEntityLink, 0)
		src.Where("target_id = ?", target.Id).Find(&oldTargetEntities)
		targetAreaSlice := map[int]bool{}
		for _, targetEntityLink := range oldTargetEntities {
			entityId := targetEntityLink.EntityId
			newResourceJoin, err := getResourceJoin(src, des, entityId, target.FieldId)
			if err != nil {
				fmt.Printf("new resource join get error %v\r\n", err.Error())
				session.Rollback()
				return
			}
			// fmt.Printf("quest_id:%v, area_id:%v\r\n", quest.Id, newResourceJoin.Area.Id)
			// fmt.Printf("area_name:%v\r\n", newResourceJoin.Resource.Name)
			_, ok := targetAreaSlice[newResourceJoin.Area.Id]
			if !ok {
				targetAreaSlice[newResourceJoin.Area.Id] = true
			}
		}
		for areaId, _ := range targetAreaSlice {
			questTarget := models.QuestTarget{}
			questTarget.QuestId = quest.Id
			questTarget.AreaId = areaId
			questTarget.Status = models.QuestTargetStatusWait
			_, err = session.Insert(&questTarget)
			if err != nil {
				fmt.Printf("quest target insert error %v\r\n", err.Error())
				session.Rollback()
				return
			}
		}
		insertNum++
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("target session commit error:%v\r\n", err)
	}
	fmt.Printf("target transfer success:%v\r\n", insertNum)
}
