package model

import "time"

//宝箱CDT分配规则
type TreasureBoxCdtConfig struct {
	Id           int     `xorm:"int(10) notnull autoincr pk" json:"id"`
	Probability  int64   `xorm:"int(20) notnull comment('概率')" json:"probability"`
	RewardItem   string  `xorm:"varchar(54) notnull comment('奖励道具名')" json:"rewardItem"`
	RewardNumber float32 `xorm:"decimal(12,2) notnull comment('奖励数量')" json:"rewardNumber"`

	CreateTime time.Time `xorm:"timestamp notnull comment('创建时间')" json:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull comment('更新时间')" json:"updateTime"`
}
