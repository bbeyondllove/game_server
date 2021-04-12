package model

import "time"

//用户任务表
type UserTask struct {
	Id         string    `xorm:"int(20) notnull pk" json:"id"`
	UserId     string    `xorm:"int(20) notnull " json:"user_id"`
	TaskId     int       `xorm:"int(10) notnull " json:"task_id"`
	TaskType   int       `xorm:"int(10) notnull " json:"task_type"`
	SigninType int       `xorm:"int(10) notnull " json:"signin_type"`
	AwardInfo  string    `xorm:"varchar(512) notnull " json:"award_info"`
	Status     int       `xorm:"int(10) notnull " json:"status"`
	CreateTime      time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime      time.Time `xorm:"timestamp notnull" json:"update_time"`
}
