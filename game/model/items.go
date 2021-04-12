package model

import "time"

//道具表
type Items struct {
	Id        int     `xorm:"int(10) notnull pk" json:"id"`
	ItemType  int     `xorm:"int(10) notnull" json:"item_type"`
	ItemName  string  `xorm:"varchar(255) notnull" json:"item_name"`
	IsBind    int     `xorm:"tinyint(1)"  json:"is_bind"`
	Quality   int     `xorm:"int(10)" json:"quality"`
	IsPile    int     `xorm:"tinyint(1)" json:"is_pile"`
	GetFrom   string  `xorm:"varchar(255)"   json:"get_from"`
	UseJump   string  `xorm:"varchar(255)"    json:"use_jump"`
	Price     float32 `xorm:"Decimal notnull" json:"Price"`
	Recommend int     `xorm:"tinyint(1) notnull " json:"recommend"`
	Desc      string  `xorm:"varchar(1024)" json:"desc"`
	Attr1     string  `xorm:"varchar(100)" json:"attr1"`
	ImgUrl    string  `xorm:"varchar(255)" json:"img_url"`
	IsGift    int     `xorm:"tinyint(1)"  json:"is_gift"`

	CreateTime time.Time `xorm:"timestamp notnull" json:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"updateTime"`
}
