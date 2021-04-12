package model

import (
	"github.com/shopspring/decimal"
)

type TransferReq struct {
	TxNo      string          `form:"txNo" json:"txNo" binding:"required"  desc:"用户系统生成的流水号"`
	UserId    string          `form:"userId" json:"userId" binding:"required"`
	TokenCode string          `form:"tokenCode" json:"tokenCode" binding:"required" desc:"划转的通证编码"`
	Amount    decimal.Decimal `form:"amount" json:"amount" binding:"required" desc:"金额数量,划转的通证数量，正数表示转入业务系统，负数表示从业务系统转出"`
	Nonce     string          `form:"nonce" json:"nonce" binding:"required" desc:"随机字符串"`
	Sign      string          `form:"sign" json:"sign" binding:"required" desc:"签名值"`
}

type HttpCommon struct {
	Code    int32  `json:"code"` //错误代码
	Message string `json:"msg"`  //错误信息
}

type BalanceReq struct {
	Nonce     string `form:"nonce" json:"nonce" binding:"required" desc:"随机字符串"`
	Sign      string `form:"sign" json:"sign" binding:"required" desc:"签名值"`
	UserId    string `form:"userId" json:"userId" binding:"required"`
	TokenCode string `form:"tokenCode" json:"tokenCode" binding:"required" desc:"划转的通证编码"`
}

type BalanceRsp struct {
	HttpCommon
	Data map[string]decimal.Decimal `json:"data"`
}

type ChangeRecordListReq struct {
	LastRecordId int64  `json:"lastRecordId" validate:"required" desc:"本地最大的对账记录ID，无则写0"`
	UserId       string `json:"userId" validate:"required"`
	DeadLine     int64  `json:"deadLine" desc:"对帐的时间上限"`
	Signature    string `json:"signature" validate:"required" desc:"签名接口,SHA256('deadLine=123456678&lastRecordId=123456567&userId=12345678&key=1qaz2wsx&key=1qaz2wsx')"`
}

type MemberIsActiveReq struct {
	UserId    string `json:"userId" validate:"required"`
	Signature string `json:"sign" validate:"required" desc:"签名接口,SHA256('userId=12345678&key=1qaz2wsx')"`
}

type CurrencySupportReq struct {
	UserId    string `json:"userId" validate:"required"`
	Currency  string `json:"currency" validate:"required" desc:"币种"`
	Signature string `json:"sign" validate:"required" desc:"签名接口,SHA256('currency=SAN&userId=12345678&key=1qaz2wsx')"`
}

type SendGameEmailReq struct {
	UserId       string `json:"userId" validate:"required"`
	EmailType    int    `json:"emailType" validate:"required" desc:"邮件类型 1 文本类型"`
	EmailTitle   string `json:"emailTitle" validate:"required" desc:"邮件标题"`
	EmailContent string `json:"emailContent" validate:"required" desc:"邮件内容"`
	Signature    string `json:"signature" validate:"required" desc:"签名接口,SHA256('emailContent=test123&emailTitle=test123&emailType=1&userId=332551164648230912&key=12345678')"`
}

type PushGameNoticeReq struct {
	NoticeType    int    `json:"noticeType" validate:"required" desc:"公告类型 1系统更新公告"`
	NoticeTitle   string `json:"noticeTitle" validate:"required" desc:"公告标题"`
	NoticeContent string `json:"noticeContent" validate:"required" desc:"公告内容"`
	Version       string `json:"version" validate:"required" desc:"版本号"`
	Signature     string `json:"signature" validate:"required" desc:"签名接口,SHA256('')"`
}

type UserRsp struct {
	UserId        string          `json:"user_id" desc:"用户ID"`
	NickName      string          `json:"nick_name" desc:"用户昵称"`
	CountryCode   int             `json:"countryCode" desc:"手机区号"`
	Mobile        string          `json:"mobile" desc:"手机号码"`
	Email         string          `json:"email" desc:"邮箱地址"`
	Status        int             `json:"status" desc:"是否禁用,0 禁用 1 未禁用"`
	KycPassed     int             `json:"kyc_passed" desc:"是否实名认证,0否，1是"`
	RegisteIp     string          `json:"registeIp" desc:"注册IP地址"`
	Cdt           decimal.Decimal `json:"cdt" desc:"用户的ctd"`
	CreateTime    string          `json:"createTime" desc:"创建时间"`
	LastLoginTime string          `json:"lastLoginTime" desc:"最后登录时间"`
}

type CertificationRsp struct {
	UserId      string `json:"user_id" desc:"用户ID"`
	NickName    string `json:"nick_name" desc:"用户昵称"`
	CountryCode int    `json:"countryCode" desc:"手机区号"`
	Mobile      string `json:"mobile" desc:"手机号码"`
	Email       string `json:"email" desc:"邮箱地址"`

	Nationality string `json:"nationality" desc:"国籍"`
	FirstName   string `json:"first_name" desc:"姓"`
	LastName    string `json:"last_name" desc:"名"`
	IdType      int    `json:"id_type" desc:"1：身份证；2：护照"`
	IdNumber    string `json:"id_number" desc:"证件号码"`
	ObjectKey   string `json:"objectKey" desc:"证件照片,{front:'',back:'',other:''},front:正面,back:反面,other:全身照"`
	Suggestion  string `json:"suggestion" desc:"审核意见"`
	Status      int    `json:"status" desc:"0：首次审核，1：审核不通过，2：审核通过 3：再次提交审核"`
	ApplyTime   string `json:"applyTime" desc:"申请时间"`
	ExamineTime string `json:"examineTime" desc:"最近审核时间"`
}
