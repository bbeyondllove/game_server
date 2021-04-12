package model

import "time"

// 用户宝箱领取记录
type UserTreasureBoxRecord struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	UserId     string    `xorm:"bigint(20) comment('用户ID')" json:"user_id"`
	Cdt        float32   `xorm:"float comment('cdt数量')" json:"cdt"`
	WatchTime  int       `xorm:"int(11) comment('观看时长')" json:"watch_time"`
	OpenTime   time.Time `xorm:"timestamp comment('打开时间点')" json:"open_time"`
	CreateTime time.Time `xorm:"created notnull comment('创建时间')" json:"create_time"`
	UpdateTime time.Time `xorm:"updated notnull comment('更新时间')" json:"update_time"`
}
