package proto

// 进入城市city
type C2SCityIcon struct {
	Token string `json:"token"`
}

// 欢乐宝箱有效广告统计
type C2SCEffectiveAdvertis struct {
	Token     string `json:"token"`
	RequestID string `json:"requestId"` // 唯一请求ID
}
