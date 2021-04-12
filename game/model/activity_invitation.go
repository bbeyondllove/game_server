package model

import "time"

type ActivityInvitation struct {
	Id           int64     `xorm:"int(20) autoincr pk" json:"-"`
	UserId       string    `xorm:"int(20)" json:"userId"`
	InviterId    string    `xorm:"int(20)" json:"inviterId"`
	IsCreateRole int       `xorm:"int(1)" json:"is_create_role"`
	CreateTime   time.Time `xorm:"created" json:"createTime"`
	UpdateTime   time.Time `xorm:"updated" json:"-"`
}

type ActivityUserInvitList struct {
	UserId     string    `json:"userId"`
	NickName   string    `json:"nickName"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}
