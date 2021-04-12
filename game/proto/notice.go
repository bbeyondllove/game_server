package proto

import (
	"game_server/game/model"
)

// 检查公告
type C2SCheckUpgradeNotice struct {
	Version int `json:"version"`
}

// 发送公告
type C2SNotice struct {
	NoticeType    int    `json:"noticeType"`
	NoticeTitle   string `json:"noticeTitle"`
	NoticeContent string `json:"noticeContent"`
	Version       string `json:"version"`
}

//返回
type S2CNotice struct {
	S2CCommon
	Notice NoticeInfo `json:"upgradeNotice,omitempty"`
}

// 公告
type NoticeInfo struct {
	NoticeType    int    `json:"noticeType"`
	NoticeTitle   string `json:"noticeTitle"`
	NoticeContent string `json:"noticeContent"`
	NoticeUrl     string `json:"noticeUrl"`
	Version       string `json:"version"`
}

type S2CNotices struct {
	S2CCommon
	Notice []NoticeInfo `json:"upgradeNotice"`
}

// 公告请求
type NoticeReq struct {
	Id            int64  `json:"id"`
	NoticeType    int    `json:"notice_type"` //0:版本更新公告,1:普通公告,2:强制更新公告,3:更新前公告'
	NoticeTitle   string `json:"notice_title"`
	NoticeContent string `json:"notice_content"`
	NoticeUrl     string `json:"notice_url"`
	Version       string `json:"version"`
	Remark        string `json:"remark"`
	NoticeTime    string `json:"notice_time"`

	WalkingLanterns []model.TimeNotice `json:"walking_lanterns" desc:"走马灯"`
}

// 公告返回结果
type NoticeRsp struct {
	Id            int64  `json:"id"`
	NoticeType    int    `json:"notice_type"`
	NoticeTitle   string `json:"notice_title"`
	NoticeContent string `json:"notice_content"`
	NoticeUrl     string `json:"notice_url"`
	Version       string `json:"version"`
	Remark        string `json:"remark"`
	OperatorName  string `json:"operator_name"`
	IsNoticed     int    `json:"is_noticed"`
	NoticeTime    string `json:"notice_time"`
	CreateTime    string `json:"create_time"`
	UpdateTime    string `json:"update_time"`

	WalkingLanterns *[]model.TimeNotice `json:"walking_lanterns" desc:"走马灯"`
}

// 走马灯消息
type LanternInfo struct {
	Level   int    `json:"level"`
	Content string `json:"content"`
}

// 走马灯
type S2CNoticeLaterns struct {
	S2CCommon
	Notice LanternInfo `json:"lantern,omitempty"`
}
