package model

import "time"

//奖励表
type Awards struct {
	Id        string  `xorm:"int(20) notnull pk" json:"id"`
	AwardType int     `xorm:"int(10) notnull" json:"award_type"`
	AwardId   int     `xorm:"int(10) notnull" json:"award_id"`
	ItemId    int     `xorm:"int(10) notnull" json:"item_id"`
	ItemNum   float32 `xorm:"Decimal notnull" json:"item_num"`

	CreateTime time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"update_time"`
}
