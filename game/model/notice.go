package model

import "time"

type Notice struct {
	Id            int64     `xorm:"int(20) autoincr pk" json:"id"`
	NoticeType    int       `xorm:"int(2) comment('类型,0:版本更新公告,1:普通公告,2:强制更新公告,3:更新前公告')" json:"notice_type"`
	NoticeTitle   string    `xorm:"varchar(128) comment('标题')" json:"notice_title"`
	NoticeContent string    `xorm:"text comment('类容')" json:"notice_content"`
	NoticeUrl     string    `xorm:"varchar(255) comment('公告的URL')" json:"notice_url"`
	Version       string    `xorm:"varchar(16) comment('版本')" json:"version"`
	Remark        string    `xorm:"varchar(255) NOT NULL comment('说明')" json:"remark"`
	Operator      int64     `xorm:"bigint(20) NOT NULL comment('操作者')" json:"operator"`
	IsNoticed     int       `xorm:"int(2) NOT NULL comment('是否已已通知,只有类型3公告才需要处理,0:未通知 1:已通知')" json:"is_noticed"`
	NoticeTime    string    `xorm:"varchar(255) NOT NULL comment('通知时间')" json:"notice_time"`
	CreateTime    time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime    time.Time `xorm:"timestamp notnull" json:"update_time"`
}

type TimeNotice struct {
	Id           int64     `xorm:"int(20) autoincr pk" json:"id"`
	NoticeId     int64     `xorm:"int(20) notnull comment('通知ID')" json:"notice_id"`
	Content      string    `xorm:"varchar(128) notnull comment('通知内容')" json:"content"`
	IntervalTime int       `xorm:"int(20) notnull comment('间隔时间,单位:秒')" json:"interval_time"`
	StartTime    time.Time `xorm:"-" json:"-" desc:"开始时间"`
	CreateTime   time.Time `xorm:"timestamp notnull comment('创建时间')" json:"-"`
	UpdateTime   time.Time `xorm:"timestamp notnull comment('更新时间')" json:"-"`
}
