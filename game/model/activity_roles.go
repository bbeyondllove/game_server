package model

import "time"

type ActivityRoles struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	UserId     string    `xorm:"int(11)" json:"userId"`
	ItemId     int       `xorm:"int(11)" json:"itemId"`
	RoleId     int       `xorm:"int(11)" json:"roleId"`
	ExpireTime time.Time `xorm:"timestamp" json:"expireTime"`
	CreateTime time.Time `xorm:"created" json:"-"`
	UpdateTime time.Time `xorm:"updated" json:"-"`
}
