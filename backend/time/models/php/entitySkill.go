package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type EntitySkill struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	Name        string    `xorm:"not null default '' comment('技能名称') VARCHAR(255)"`
	Description string    `xorm:"not null comment('技能描述') TEXT"`
	RankDesc    string    `xorm:"comment('等级描述，用逗号分隔') TEXT"`
	ImgUrl      string    `xorm:"comment('技能图片') TEXT"`
	MaxPoints   int       `xorm:"not null comment('最大级别') SMALLINT(6)"`
	TypeId      int       `xorm:"not null default 0 comment('所处层级') INT(11)"`
	Ctime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	Utime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
}

func (c *EntitySkill) Transfer(src, des *xorm.Engine) {
	session := des.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		fmt.Printf("skill session error:%v\r\n", err.Error())
		return
	}

	oldLinks := make([]AreaSkillLink, 0)
	src.Find(&oldLinks)
	insertNum := 0
	areaMap := map[string][]int{}
	resourceMap := map[string]models.Resource{}
	for _, one := range oldLinks {
		resource := models.Resource{}
		skill := new(EntitySkill)
		_, err := src.Id(one.SkillId).Get(skill)
		if err != nil {
			fmt.Printf("skillId:%v, not exist\r\n", one.SkillId)
		}
		resource.Name = skill.Name
		resource.Desc = skill.Description
		resource.Year = 0
		resource.CreatedAt = one.Ctime
		resource.UpdatedAt = one.Utime

		areaSlice, ok := areaMap[resource.Name]
		if !ok {
			areaMap[resource.Name] = make([]int, 0)
		}
		areaSlice = append(areaSlice, one.AreaId)
		areaMap[resource.Name] = areaSlice

		if _, ok := resourceMap[resource.Name]; !ok {
			resourceMap[resource.Name] = resource
		}

	}
	for resourceName, areaSlice := range areaMap {
		resource, ok := resourceMap[resourceName]
		if !ok {
			session.Rollback()
			fmt.Printf("work resource can not find in map\r\n")
			return
		}
		_, err = session.Insert(&resource)
		if err != nil {
			session.Rollback()
			fmt.Printf("skill resource insert error:%v\r\n", err.Error())
			return
		}
		for _, areaId := range areaSlice {
			mapAreaResource := models.MapAreaResource{}
			mapAreaResource.AreaId = areaId
			mapAreaResource.ResourceId = resource.Id
			_, err = session.Insert(&mapAreaResource)
			if err != nil {
				session.Rollback()
				fmt.Printf("skill resource map insert error:%v\r\n", err.Error())
				return
			}
		}
		insertNum++
	}
	err = session.Commit()
	if err != nil {
		fmt.Printf("skill session commit error:%v\r\n", err)
	}
	fmt.Printf("skill transfer success:%v\r\n", insertNum)
}
