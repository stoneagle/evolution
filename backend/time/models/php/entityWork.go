package php

import (
	"evolution/backend/time/models"
	"fmt"
	"time"

"github.com/go-xorm/xorm"


ype EntityWork struct {
Id        int       `xorm:"not null pk autoincr INT(11)"`
Name      string    `xorm:"not null default '' comment('名称') VARCHAR(255)"`
Desc      string    `xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Year      int       `xorm:"not null default 0 comment('年代') INT(11)"`
	CountryId int       `xorm:"not null default 0 comment('所属国家') INT(11)"`
	AreaId    int       `xorm:"not null default 0 comment('所属领域') INT(11)"`
	Ctime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
Utime     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`


unc (c *EntityWork) Transfer(src, des *xorm.Engine) {
olds := make([]EntityWork, 0)
src.Find(&olds)
news := make([]models.Resource, 0)
for _, one := range olds {
	tmp := models.Resource{}
	tmp.Name = one.Name
	tmp.Desc = one.Desc
	tmp.Year = one.Year
		tmp.AreaId = one.AreaId
		tmp.CreatedAt = one.Ctime
		tmp.UpdatedAt = one.Utime
		news = append(news, tmp)
}
affected, err := des.Insert(&news)
if err != nil {
	fmt.Printf("entity work transfer error:%v\r\n", err.Error())
} else {
	fmt.Printf("entity work transfer success:%v\r\n", affected)
}


