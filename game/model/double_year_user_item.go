package model

import "time"

//玩家背包表
type DoubleYearUserItem struct {
	Id       int    `xorm:"int(10) notnull pk" json:"id"`
	UserId   string `xorm:"int(20) notnull" json:"user_id"`
	ItemType int    `xorm:"int(10) notnull" json:"item_type"`
	ItemId   int    `xorm:"int(10) notnull" json:"item_id"`
	ItemNum  int    `xorm:"int(10) notnull" json:"item_num"`
	IsTrade  int    `xorm:"int(10) notnull" json:"is_trade"`

	CreateTime time.Time `xorm:"timestamp notnull" json:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"updateTime"`
}
