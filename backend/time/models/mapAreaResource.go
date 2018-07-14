package models

type MapAreaResource struct {
	AreaId     int `xorm:"unique not null default 0 comment('隶属领域') INT(11)" structs:"area_id,omitempty"`
	ResourceId int `xorm:"unique(area_id) not null default 0 comment('隶属资源') INT(11)" structs:"resource_id,omitempty"`
}
