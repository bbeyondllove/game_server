package model

import "time"

//任务表
type Tasks struct {
	Id          int    `xorm:"int(10) notnull pk" json:"id"`
	TaskType    int    `xorm:"int(10) notnull " json:"task_type"`
	Title       string `xorm:"varchar(256) notnull " json:"title"`
	Desc        string `xorm:"varchar(512) notnull " json:"desc"`
	TaskKey     string `xorm:"varchar(256) notnull " json:"task_key"`
	TaskValue   int    `xorm:"int(10) notnull " json:"task_value"`
	MessageId   string `xorm:"varchar(256) notnull " json:"message_id"`
	EventId     int    `xorm:"int(10) notnull " json:"event_id"`
	FrontTaskId int    `xorm:"int(10) notnull " json:"front_task_id"`
	Status      int    `xorm:"tinyint(1)"  json:"status"`

	CreateTime time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"update_time"`
}

//任务奖励表
type TaskAward struct {
	TaskId   int    `xorm:"int(10) notnull pk" json:"task_id"`
	AwardId  string `xorm:"varchar(256) notnull " json:"award_id"`
	AwardNum string `xorm:"varchar(256) notnull " json:"award_num"`

	CreateTime time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"update_time"`
}
