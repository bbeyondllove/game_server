package proto

import "sync"

//签到请求
type C2SSignIn struct {
	Token      string `json:"token"`
	SigninType int    `json:"signinType"` //1：正常签到；2：补签
	ItemId     int    `json:"itemId"`     // SigninType=2时，ItemId为补签卡id
}

//任务奖励领取请求
type C2STaskAward struct {
	Token  string `json:"token"`
	TaskId int    `json:"task_id"`
}

type AwardItem struct {
	ItemId   int     `json:"itemId"`   // 奖励id
	ItemNum  float32 `json:"itemNum"`  // 奖励数量
	IsGift   bool    `json:"isGift"`   //是否礼包(0：否；1：是）
	ImgUrl   string  `json:"imgUrl"`   //图片资源
	Desc     string  `json:"desc"`     //描述
	Attr1    string  `json:"attr1"`    //绑定属性1，目前表示解锁卡对应的角色ID
	ItemName string  `json:"itemName"` //名称
	Sex      int     `json:"sex"`      //性别
}

type AwardInfo struct {
	AwardItem
	AwardList []*AwardItem `json:"awardList"` //奖励列表
}

//任务列表
type TaskItem struct {
	Id       int    `json:"id"`
	TaskType int    `json:"taskType"` //任务类型（1：主线任务；2：每日任务；3：签到）
	Title    string `json:"title"`    //标题
	Desc     string `json:"desc"`     //描述
	Status   int    `json:"status"`   //0:未做；1：已做; 2:已领取

	TaskKey     string     `json:"taskKey"`     //任务缓存在redis的key
	TaskValue   int        `json:"taskValue"`   //任务值
	MessageId   string     `json:"messageId"`   //触发任务对应的协议id，多个协议id用|分隔
	EventId     int        `json:"eventId"`     //触发任务同时要满足的事件id(0表示所有事件）
	FrontTaskId int        `json:"frontTaskId"` //前置任务id
	Awards      *AwardInfo `json:"awards"`      //奖励列表
}

//任务获取返回
type S2CTaskList struct {
	S2CCommon
	SigninDayNo   int         `json:"signinDayNo"`   //今天要签第几天
	CurDayStatus  bool        `json:"curDayStatus"`  //今天是否已完成(false：未完成；true:已完成)
	CurWeekStatus bool        `json:"curWeekStatus"` //本周是否已补签(false：未完成；true:已完成)
	TaskList      []*TaskItem `json:"taskList"`
}

//签到返回或任务奖励领取返回
type S2CTaskAward struct {
	S2CCommon
	Awards *AwardInfo `json:"awards"` //奖励列表
}

//
type UserTask struct {
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	TaskId     int    `json:"task_id"`
	TaskType   int    `json:"task_type"`
	SigninType int    `json:"signin_type"`
	AwardInfo  string `json:"award_info"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

//任务完成推送
type PushTaskStatus struct {
	S2CCommon
	TaskId int        `json:"task_id"`
	Awards *AwardInfo `json:"awards"` //奖励列表
}

//后台改变任务开启状态
type ChangeTaskStatus struct {
	Status int `json:"status"` //活动关闭或开启。0：关闭；1：开启
}

type TaskStart struct {
	Status int //活动关闭或开启。0：关闭；1：开启
	Mt     sync.RWMutex
}

//后台改变任务开启状态
type S2CChangeTaskStatus struct {
	S2CCommon
	Status int `json:"status"` //活动关闭或开启。0：关闭；1：开启
}
