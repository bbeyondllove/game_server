package proto

import "time"

type C2SActivitStatus struct {
	Token string `json:"token"`
}

//打开宝箱返回
type S2CActivitStatus struct {
	S2CCommon
	Activity Activity `json:"activity"`
}

type Activity struct {
	TreasureBox  int `json:"treasureBox"`  // 欢乐宝箱
	RandList     int `json:"randList"`     // 排行榜
	ExchangeCdt  int `json:"exchangeCdt"`  // 兑换CDT
	ExchangeRole int `json:"exchangeRole"` // 兑换角色
}

// 活动状态信息
type ActivityItem struct {
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Status    int       `json:"status"` // 0 关闭 1 开启
	BeginTime time.Time `json:"begin_time"`
	EndTime   time.Time `json:"end_time"`
}

type ActivityStatusRsp struct {
	Type      int    `json:"type"`       // 活动类型
	Name      string `json:"name"`       // 活动名称
	Status    int    `json:"status"`     // 0 未开始 1 开启 2 结束
	StartTime string `json:"start_time"` // 开启时间
	EndTime   string `json:"end_time"`   // 结束时间
}
