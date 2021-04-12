package model

import "time"

type ActivityRedEnvelope struct {
	Id         int       `xorm:"int autoincr pk" json:"id"`
	UserId     string    `xorm:"int(20)" json:"user_id"`
	Number     float32   `xorm:"decimal(12,4)" json:"number"`
	IsReceive  int       `xorm:"int(1)" json:"is_receive"`
	CreateTime time.Time `xorm:"created" json:"create_time"`
	UpdateTime time.Time `xorm:"updated" json:"update_time"`
}
