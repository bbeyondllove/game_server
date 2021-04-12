package model

import (
	"time"
)

//世界地图坐标表
type WorldMap struct {
	Id              string `xorm:"int(20) notnull pk" json:"id"`
	PositionX       int    `xorm:"int(6) notnull" json:"position_x"`
	PositionY       int    `xorm:"int(6) notnull" json:"position_y"`
	SmallType       string `xorm:"varchar(20)" json:"small_type"`
	CanSale         int    `xorm:"int(1) notnull" json:"can_sale"`
	IsSale          int    `xorm:"int(1) notnull" json:"is_sale"`
	BuildingName    string `xorm:"varchar(50)" json:"building_name"` //建筑名
	ShopName        string `xorm:"varchar(50)" json:"shop_name"`     //商家名称
	Desc            string `xorm:"varchar(500)" json:"desc"`
	H5Url           string `xorm:"varchar(80)" json:"h5_url"`
	WebUrl          string `xorm:"varchar(80)" json:"web_url"`
	PassportAviable string `xorm:"varchar(200)" json:"passport_aviable"`
	ImageUrl        string `xorm:"varchar(80)" json:"image_url"`

	CreateTime time.Time `xorm:"timestamp notnull" json:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"updateTime"`
}
