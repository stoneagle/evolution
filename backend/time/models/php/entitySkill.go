package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

"github.com/go-xorm/xorm"


ype EntitySkill struct {
Id          int       `xorm:"not null pk autoincr INT(11)"`
Name        string    `xorm:"not null default '' comment('技能名称') VARCHAR(255)"`
Description string    `xorm:"not null comment('技能描述') TEXT"`
	RankDesc    string    `xorm:"comment('等级描述，用逗号分隔') TEXT"`
	ImgUrl      string    `xorm:"comment('技能图片') TEXT"`
	MaxPoints   int       `xorm:"not null comment('最大级别') SMALLINT(6)"`
	TypeId      int       `xorm:"not null default 0 comment('所处层级') INT(11)"`
	Ctime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
Utime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`


unc (c *EntitySkill) Transfer(src, des *xorm.Engine) {
oldLinks := make([]AreaSkillLink, 0)
src.Find(&oldLinks)
news := make([]models.Resource, 0)
for _, one := range oldLinks {
	tmp := models.Resource{}
	skill := new(EntitySkill)
	_, err := src.Id(one.SkillId).Get(skill)
	if err != nil {
		fmt.Printf("skillId:%v, not exist\r\n", one.SkillId)
		}
		tmp.Name = skill.Name
		tmp.Desc = skill.Description
		tmp.Year = 0
		tmp.AreaId = one.AreaId
		tmp.CreatedAt = one.Ctime
		tmp.UpdatedAt = one.Utime
		news = append(news, tmp)
}
affected, err := des.Insert(&news)
if err != nil {
	fmt.Printf("entity skill transfer error:%v\r\n", err.Error())
} else {
	fmt.Printf("entity skill transfer success:%v\r\n", affected)
}


