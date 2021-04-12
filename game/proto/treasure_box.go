package proto

import "github.com/shopspring/decimal"

//打开宝箱
type C2SOpenStreasureBox struct {
	Token        string `json:"token"`
	ActivityType int    `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件 2.宝箱
	LocationID   int    `json:"locationId"`   //位置ID
	X            int    `json:"x"`            //玩家当前x坐标
	Y            int    `json:"y"`            //玩家当前y坐标
}

//完成宝箱
type C2SFinishStreasureBox struct {
	Token        string `json:"token"`
	FinishStatus int    `json:"finishStatus"` //完成状态 0 未完成， 1 已完成
	ActivityType int    `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件 2.宝箱
	LocationID   int    `json:"locationId"`   //位置ID
	X            int    `json:"x"`            //玩家当前x坐标
	Y            int    `json:"y"`            //玩家当前y坐标
}

//打开宝箱返回
type S2COpenStreasureBox struct {
	S2CCommon
	DayNum        int64 `json:"dayNum"`        //每天的数量
	ResidueDegree int64 `json:"residueDegree"` // 剩余次数
}

//完成宝箱返回
type S2CFinishStreasureBox struct {
	S2CCommon
	LocationID   int         `json:"locationId"`   //位置ID
	ActivityType int         `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件
	X            int         `json:"x"`            //玩家当前x坐标
	Y            int         `json:"y"`            //玩家当前y坐标
	UserId       string      `json:"userId"`
	BoxAward     interface{} `json:"boxAward,omitempty"` //奖励
}

type BoxAward struct {
	AwardId  int     `json:"awardId"`
	ItemId   int     `json:"itemId"`
	ItemNum  float32 `json:"itemNum"`
	ImgUrl   string  `json:"imgUrl"`
	Desc     string  `json:"desc"`
	ItemName string  `json:"itemName"`
}

//领取宝箱奖励
type C2SReceiveBoxReward struct {
	Token   string `json:"token"`
	AwardId int    `json:"awardId"`
}

//领取宝箱奖励
type S2CReceiveBoxReward struct {
	S2CCommon
	Cdt float32 `json:"ctd,omitempty"` //用户更新后的CDT
}

//{
//"code": 0,
//"message": "成功",
//"awards": {
//"itemId": 1001,
//"itemNum": 0.05,
//"isGift": false,
//"imgUrl": "https://xx-admin-client.oss-cn-shenzhen.aliyuncs.com/test/common/1608519513584.png",
//"desc": "CDT",
//"attr1": "0",
//"itemName": "CDT",
//"sex": 1,
//"awardList": null
//}
//}

//宝箱记录查询
type C2SStreasureBoxGetRecord struct {
	// Page  int    `json:"page"` //第几页
	// Size  int    `json:"size"` //每页多少数量
	Token string `json:"token"`
}

//宝箱记录
type StreasureBoxRecord struct {
	OpenTime  string          `json:"openTime"`
	Cdt       decimal.Decimal `json:"cdt"`
	WatchTime int             `json:"watchTime"`
}

type Paging struct {
	TotalSize   int `json:"totalSize"`
	CurrentSize int `json:"currentSize"`
	CurrentPage int `json:"currentPage"`
}

//获取宝箱记录
type S2CStreasureBoxRecordList struct {
	S2CCommon
	// Paging   Paging               `json:"paging"`
	Record   []StreasureBoxRecord `json:"data"`
	TotalCdt decimal.Decimal      `json:"totalCdt"`
	WatchNum int64                `json:"watchNum"`
}

type C2SStreasureBoxDayNum struct {
	Token string `json:"token"`
}

//打开宝箱返回
type S2CStreasureBoxDayNum struct {
	S2CCommon
	DayNum        int64 `json:"dayNum"`        //每天的数量
	ResidueDegree int64 `json:"residueDegree"` // 剩余次数
}
