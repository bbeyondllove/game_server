package proto

import "github.com/shopspring/decimal"

//获取圣诞糖果，圣诞树数据返回
type S2CGetSweetTree struct {
	S2CCommon
	UserId      string             `json:"userId"`      //用户ID
	UserDayCdt  decimal.Decimal    `json:userDayCdt"`   //当天当前用户总共兑换的cdt
	TotalDayCdt decimal.Decimal    `json:"totalDayCdt"` //当天所有用户总共兑换的cdt
	CurrentCdt  decimal.Decimal    `json:"currentCdt"`  //当前用户的cdt
	SweetTree   map[int][]ItemInfo `json:"sweetTree"`   //圣诞糖果，圣诞树信息
}

//获取碎片返回
type S2CGetPatch struct {
	S2CCommon
	UserId       string `json:"userId"`       //用户ID
	CurrentNum   int    `json:currentNum"`    //当前数量
	TradeNeedNum int    `json:"tradeNeedNum"` //兑换圣诞老人所需数量
}

//兑换CDT响应
type S2CTradeCdt struct {
	S2CCommon
	UserId    string  `json:"userId"`   //用户ID
	ChangeCdt float32 `json:"changeCdt` //用户兑换的cdt
}

//兑换圣诞老人解锁卡响应
type S2CTradeSantaClaus struct {
	S2CCommon
	UserId string   `json:"userId"` //用户ID
	Num    int      `json:"num"`    //得到的解锁卡数量
	Item   ItemInfo `json:"item"`
}

//通知圣诞老人已过期
type S2CSantaClausRoleExprie struct {
	S2CCommon
	UserId string `json:"userId"` //用户ID
	RoleId int    `json:"roleId"` //用户ID
}

//通知活动状态
type S2CDoubleYearStatus struct {
	S2CCommon
	StatusType int `json:"statusType"` //0:双旦;1:排行榜;
	Status     int `json:"status"`     //0：双蛋活动未开始 1:双蛋活动进行中 2:双蛋活动已结束
}

//通知得到圣诞老人解锁卡
type S2CBroadSantaCard struct {
	S2CCommon
	Content string `json:"content"`
}

//  获取活动状态
type C2SGetActiveStatus struct {
	Token      string `json:"token"`
	StatusType int    `json:"statusType"` //0:双旦;1:排行榜;
}

//发送红包
type S2CRedEnvelope struct {
	S2CCommon
	Cdt      float32 `json:"cdt"`      // 红包的CDT
	CdtalCdt float32 `json:"totalCdt"` // 用户总的CDT
}

//推 兑换红包信息，跑马灯
type S2CRedEnvelopeMsg struct {
	S2CCommon
	Cdt      float32 `json:"cdt"`      // cdt
	UserName string  `json:"userName"` // 用户昵称
	IsTotal  int     `json:"isTotal"`  //  0 单次   1 总的
}
