package model

import "time"

type BuildType struct {
	Id           string    `xorm:"int(20) notnull pk" json:"id"`
	SmallType    string    `xorm:"varchar(20)" json:"small_type"`
	BuildingName string    `xorm:"varchar(255) notnull" json:"building_name"`
	ImageName    string    `xorm:"varchar(50)" json:"image_name"`
	CanSale      int       `xorm:"int(1) notnull" json:"can_sale"`
	CreateTime   time.Time `xorm:"timestamp notnull" json:"createTime"`
	UpdateTime   time.Time `xorm:"timestamp notnull" json:"updateTime"`
}
