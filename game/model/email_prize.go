package model

import "time"

type EmailPrize struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	EmailId    int       `xorm:"int(11)" json:"emailId"`
	PrizeType  int       `xorm:"int(11)" json:"prizeType"`
	PrizeId    int       `xorm:"int(11)" json:"prizeId"`
	PrizeName  string    `xorm:"varchar(64)" json:"przieName"`
	PrizeNum   int       `xorm:"int(11)" json:"prizeNum"`
	PrizeImg   string    `xorm:"varchar(255)" json:"prizeImg"`
	IsReceive  int       `xorm:"int(1)" json:"isReceive"`
	Extend     string    `xorm:"varchar(255)" json:"-"`
	CreateTime time.Time `xorm:"created" json:"-"`
	UpdateTime time.Time `xorm:"updated" json:"-"`
}
