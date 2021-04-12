package proto

import "game_server/game/model"

type EmailBase struct {
	S2CCommon
}

// 邮件列表
type C2SEmailList struct {
	IsRead int    `json:"isRead"`
	Token  string `json:"token"`
}

type S2CEmailList struct {
	IsPush int `json:"isPush"`
	S2CCommon
	ReadList   []EmailItem `json:"readList"`
	UnreadList []EmailItem `json:"unreadList"`
}

// 删除邮件
type C2SDelEmail struct {
	Token    string `json:"token"`
	EmailIds []int  `json:"emailIds"`
}
type S2CDelEmail struct {
	S2CCommon
}

// 设置邮件已读
type C2SSetEmailRead struct {
	Token    string `json:"token"`
	EmailIds []int  `json:"emailIds"`
}

// 领取邮件奖励
type C2SSetReceiveRewards struct {
	Token   string `json:"token"`
	EmailId int    `json:"emailId"`
}

type S2CSetEmailRead struct {
	S2CCommon
	ReadNum   int64 `json:"readNum"`
	UnreadNum int64 `json:"unreadNum"`
}

// 获取邮件数量
type C2SGetEmialCount struct {
	Token string `json:"token"`
}

type S2CGetEmialCount struct {
	S2CCommon
}

type EmailItem struct {
	Id           int64              `json:"id"`
	UserId       string             `json:"userId"`
	EmailType    int                `json:"emailType"`
	EmailTitle   string             `json:"emailTitle"`
	EmailContent string             `json:"emailContent"`
	IsRead       int                `json:"isRead"`
	PrizeList    []model.EmailPrize `json:"prizeList"`
	CreateTime   string             `json:"createTime"`
	ExpireTime   string             `json:"expireTime"`
}

//给用户推送邮件
type S2CSendEmail struct {
	IsPush int `json:"isPush"`
	S2CCommon
	NewEmailList []EmailItem `json:"newEmailList"`
}

type EmailPrize struct {
	Id        int64 `json:"id"`
	EmailId   int   `json:"emailId"`
	EmailType int   `json:"emailType"`
	PrizeType int   `json:"prizeType"`
	PrizeId   int   `json:"prizeId"`
	PrizeNum  int   `json:"prizeNum"`
	PrizeImg  int   `json:"prizeImg"`
	IsReceive int   `json:"isReceive"`
}
