package model

import "time"

//已推送给玩家公告
type UserNotice struct {
	Id         int       `xorm:"int(10) notnull pk" json:"id"`
	NoticeId   int       `xorm:"int(20) notnull comment('公告ID')" json:"noticeId"`
	UserId     string    `xorm:"int(20) notnull comment('用户ID')" json:"userId"`
	CreateTime time.Time `xorm:"timestamp notnull comment('创建时间')" json:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull comment('更新时间')" json:"updateTime"`
}
