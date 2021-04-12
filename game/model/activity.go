package model

import "time"

type ActivityConfig struct {
	Id           int       `xorm:"int autoincr pk" json:"id"`
	ActivityType int       `xorm:"int(20) comment('活动类型,10:签到,7000:广告,8000:春节活动')" json:"activity_type"`
	Name         string    `xorm:"varchar(256) comment('活动名称')" json:"name"`
	StartTime    time.Time `xorm:"DATETIME notnull comment('开始时间')" json:"start_time"`
	FinishTime   time.Time `xorm:"DATETIME notnull comment('结束时间')" json:"finish_time"`
	CreateTime   time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime   time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}
