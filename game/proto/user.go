package proto

import (
	"github.com/shopspring/decimal"
)

//用户邮箱注册
type C2SRegisterEmail struct {
	SysType   int    `json:"sysType"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	EmailCode string `json:"emailCode"`
	Inviter   string `json:"inviter"`  //邀请人，可以为邀请人的手机号、邮箱或邀请码
	NickNname string `json:"nickName"` //昵称，当不填时则设置为注册的邮箱
	Platform  string `json:"platform"` //平台，1 web 2 安卓 3 IOS
}

//用户手机注册
type C2SRegisterMobile struct {
	SysType     int    `json:"sysType"`
	CountryCode int    `json:"countryCode"`
	Mobile      string `json:"mobile"`
	Password    string `json:"password"`
	SmsCode     string `json:"smsCode"`
	Inviter     string `json:"inviter"`  //邀请人，可以为邀请人的手机号、邮箱或邀请码
	NickNname   string `json:"nickName"` //昵称，当不填时则设置为注册的手机号
	Platform    string `json:"platform"` //平台，1 web 2 安卓 3 IOS
}

//用户登录
type C2SLogin struct {
	SysType     int    `json:"sysType"`
	CodeType    int    `json:"codeType"` //1：手机方式；2：邮箱方式
	CountryCode int    `json:"countryCode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ClientIp    string `json:"clientIp"`
	Platform    string `json:"platform"` //平台，1 web 2 安卓 3 IOS
	Version     string `json:"version"`  //平台版本
}

type S2CLogin struct {
	S2CCommon
	SysType        int             `json:"sysType"`
	UserId         string          `json:"userId"`
	CountryCode    int             `json:"countryCode"`
	Email          string          `json:"email"`
	Sex            int             `json:"sex"`
	InviteCode     string          `json:"inviteCode"`
	Level          int             `json:"level"`
	Mobile         string          `json:"mobile"`
	NickNname      string          `json:"nickName"`
	RegType        int             `json:"regType"`
	UserType       int             `json:"userType"`
	RoleId         int             `json:"roleId"`
	Status         int             `json:"status"`
	Rank           string          `json:"rank"`
	Token          string          `json:"token"`
	PositsionX     int             `json:"positsionX"`
	PositsionY     int             `json:"positsionY"`
	ModifyNameNum  int             `json:"modifyNameNum"` //改名次数
	HouseNum       int             `json:"houseNum"`
	LocationID     int             `json:"locationId"` //位置ID
	InviterId      string          `json:"inviterId"`
	KycPassed      int             `json:"kycPassed"`
	KycStatus      int             `json:"kycStatus"`
	AvailableRoles string          `json:"availableRoles"`
	DeblockedRoles string          `json:"deblockedRoles"`
	Suggestion     int             `json:"suggestion"`
	Point          int             `json:"point"`
	Cdt            decimal.Decimal `json:"cdt"`            //cdt
	ActivityStatus int             `json:"activityStatus"` //活动关闭或开启。0：关闭；1：开启
}

type S2CHttpLogin struct {
	S2CCommon
	UserId    string `json:"userId"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

//重置密码
type C2SResetPwd struct {
	SysType          int32  `json:"sysType"`
	CodeType         int    `json:"codeType"` //1：手机方式；2：邮箱方式
	CountryCode      int    `json:"countryCode"`
	Mobile           string `json:"mobile"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	VerificationCode string `json:"verificationCode"`
}

//检测帐号
type C2SCheckAccount struct {
	CountryCode int    `json:"countryCode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	RegType     int32  `json:"regType"` //1手机注册，2邮箱注册
	SysType     int32  `json:"sysType"` //1 商城 2游戏
}

//获取背包信息
type C2SGetKnapsack struct {
	Token string `json:"token"`
}

//获取背包信息响应
type S2CGetKnapsack struct {
	S2CCommon
	ItemInfos map[int][]ItemInfo `json:"itemInfo"`
	UserId    string             `json:"userId"`
}
type ItemInfo struct {
	ItemId   int    `json:"itemId"`
	Num      int    `json:"num"`
	Desc     string `json:"desc"`
	Attr1    string `json:"attr1"`
	ItemName string `json:"itemName"`
	Sex      int    `json:"sex"`
	ImgUrl   string `json:"imgUrl"`
}

//用户背包和金额同步
type C2SUpdateItemInfo struct {
	Token string `json:"token"`
}

//用户背包和金额同步响应
type S2CUpdateItemInfo struct {
	S2CCommon
	Cdt       decimal.Decimal    `json:"cdt"` //cdt
	ItemInfos map[int][]ItemInfo `json:"itemInfo"`
}

//断线重连
type C2SRebind struct {
	LocationID int32  `json:"locationId"` //位置ID
	Token      string `json:"token"`
}

//获取用户钱包金额
type C2SGetAmount struct {
	Token     string `json:"token"`
	TokenCode string `json:"tokenCode"`
}

type S2CGetAmount struct {
	S2CCommon
	UserId          string          `json:"userId"`
	Amount          float64         `json:"amount"`
	AmountAvailable float64         `json:"amountAvailable"`
	AmountBlocked   float64         `json:"amountBlocked"`
	Cdt             decimal.Decimal `json:"cdt"` //cdt
	TokenCode       string          `json:"tokenCode"`
}

//获取战友列表
type C2SInviteUserList struct {
	SysType         int    `json:"sysType"`         //1商城 2游戏 3交易所
	InviterId       string `json:"inviterId"`       //此用户id对应用户推荐表里的邀请id
	InvitationLevel int    `json:"invitationLevel"` //第几级战友，取值只能是1， 2， 3
	VipLevel        int    `json:"vipLevel"`        //vip等级 0或不传 表示获取所有列表
	Page            int    `json:"page"`            //第几页
	Size            int    `json:"size"`            //每页多少数量
	Token           string `json:"token"`
}

//获取战友列表响应
type S2CInviteUserList struct {
	S2CCommon
	Paging      paging       `json:"paging"` //分页对象
	InviteUsers []InviteUser `json:"data"`   //此用户id对应用户推荐表里的邀请id
}

type paging struct {
	TotalSize   int `json:"totalSize"`
	CurrentSize int `json:"currentSize"`
	CurrentPage int `json:"currentPage"`
}

type InviteUser struct {
	UserId          string `json:"userId"`
	InvitationLevel int    `json:"invitationLevel"`
	NickName        string `json:"nickName"`
	InviteTime      int64  `json:"inviteTime"`
	Level           int    `json:"vipLevel"`
}

//获取被抢走战友列表
type C2SGetGrabComrades struct {
	UserId  string `json:"userId"`  //此用户id对应用户推荐表里的邀请id
	SysType int    `json:"sysType"` //1商城 2游戏 3交易所
	Page    int    `json:"page"`    //第几页
	Size    int    `json:"size"`    //每页多少数量
	Token   string `json:"token"`
}

//获取被抢走战友列表响应
type S2CGetGrabComrades struct {
	S2CCommon
	Total       int          `json:"total"` //总数量
	InviteUsers []InviteUser `json:"list"`  //此用户id对应用户推荐表里的邀请id
}

//获取用户钱包金额
type C2SGetUserInfo struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

type S2CGetUserInfo struct {
	S2CCommon
	UserId   string `json:"userId"`
	RoleId   int    `json:"roleId"`
	NickName string `json:"nickName"`
	Level    int    `json:"level"`
	Rank     string `json:"rank"`
}

//获取会员等级体系数据
type C2SGetMemberSys struct {
	Token   string `json:"token"`
	SysType int    `json:"sysType"`
}

type S2CGetMemberSys struct {
	S2CCommon
	Data []MemberInfo `json:"data"`
}

type UserLevelItem struct {
	ItemId   int    `xorm:"int(20)" json:"item_id" desc:"道具卡id"`
	ItemType int    `xorm:"int(20)" json:"item_type" desc:"道具卡类型id"`
	ItemName string `xorm:"int(20)" json:"item_name" desc:"道具卡名称"`
	ItemNum  int    `xorm:"int(20)" json:"item_num" desc:"道具卡数量"`
}

type MemberInfo struct {
	MemberLevel              int             `json:"vipLevel"`
	RequireSelfLocked        decimal.Decimal `json:"requireSelfLocked"`
	RequireSan               decimal.Decimal `json:"requireSan"`
	RequireSanOnly           decimal.Decimal `json:"requireSanOnly"`
	RequireVip1              int             `json:"requireVip1"`
	RequireVip2              int             `json:"requireVip2"`
	RequireVip3              int             `json:"requireVip3"`
	RequireVip4              int             `json:"requireVip4"`
	L1InviteVip1BonusPercent int             `json:"l1InviteVip1BonusPercent"`
	L1InviteVip2BonusPercent int             `json:"l1InviteVip2BonusPercent"`
	L1InviteVip3BonusPercent int             `json:"l1InviteVip3BonusPercent"`
	L1InviteVip4BonusPercent int             `json:"l1InviteVip4BonusPercent"`
	L1InviteVip5BonusPercent int             `json:"l1InviteVip5BonusPercent"`
	L2InviteBonusPercent     int             `json:"l2InviteBonusPercent"`
	L3InviteBonusPercent     int             `json:"l3InviteBonusPercent"`
	ItemList                 []UserLevelItem `json:"itemList" desc:"等级享有的道具卡"`
}

//获取用户当前等级状态
type C2SGetUserLevel struct {
	Token   string `json:"token"`
	SysType int    `json:"sysType"`
}

type S2CGetUserLevel struct {
	S2CCommon
	Member UserLevel `json:"data"`
}
type UserLevel struct {
	MemberLevel      int             `json:"level"`
	Desc             string          `json:"desc"`
	AccmSelfLocked   decimal.Decimal `json:"accumSelfLocked"`
	AccmSan          decimal.Decimal `json:"accumSan"`
	VipStartDate     int             `json:"vipStartDate"`
	VipEndDate       int             `json:"vipEndDate"`
	Vip1InviteeCount int             `json:"vip1InviteeCount"`
	Vip2InviteeCount int             `json:"vip2InviteeCount"`
	Vip3InviteeCount int             `json:"vip3InviteeCount"`
	Vip4InviteeCount int             `json:"vip4InviteeCount"`
	Vip5InviteeCount int             `json:"vip5InviteeCount"`
}

//获取商品列表
type S2CItemList struct {
	S2CCommon
	ProductItemList map[int][]*ProductItem `json:"itemList"`
}

//获取实名认证用户列表
type C2SGetCertUser struct {
	Token string `json:"token"`
}

//获取实名认证用户列表返回
type S2CCertUserList struct {
	UserId    string `json:"user_id" `
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IdType    int    `json:"id_type"`
	IdNumber  int    `json:"id_number"`
}

type ProductItem struct {
	ItemId    int             `json:"itemId"`
	ItemType  int             `json:"itemType"`
	ItemName  string          `json:"itemName"`
	IsBind    int             `json:"isBind"`
	IsGift    int             `json:"isGift"`
	Quality   int             `json:"quality"`
	IsPile    int             `json:"isPile"`
	GetFrom   string          `json:"getFrom"`
	UseJump   string          `json:"useJump"`
	Price     decimal.Decimal `json:"Price"`
	Recommend int             `json:"recommend"`
	Desc      string          `json:"desc"`
	Attr1     string          `json:"attr1"`
	ImgUrl    string          `json:"imgUrl"`
}

//抢战友
type C2SGrabComrade struct {
	UserId   string `json:"userId"`   //战友ID
	InviteId string `json:"inviteId"` //抢劫者用户ID
	SysType  int    `json:"sysType"`  //系统类型
	Token    string `json:"token"`
}

type CityUser struct {
	UserList map[int32][]string `json:"userList"`
}

type S2CCityUser struct {
	S2CCommon
	ActionType int                     `json:"actionType"` //1:进入城市；2：离开城市
	UserList   map[string]UserPosition `json:"userList"`
}

type C2SBindInviter struct {
	Token      string `json:"token"`
	InviteCode string `json:"inviteCode"`
}

type C2SDepositRebate struct {
	Token string `json:"token"`
	Page  int    `json:"page"` //第几页
	Size  int    `json:"size"` //每页多少数量
}

type S2CDepositRebate struct {
	S2CCommon
	Data   []RebateDetail `json:"data"`
	Paging paging         `json:"paging"` //分页对象
}
type RebateDetail struct {
	UserId       string          `json:"userId"`
	InviteIncome decimal.Decimal `json:"inviteIncome"`
	NickName     string          `json:"friendNickName"`
	FriendLevel  int             `json:"friendLevel"`
	VipLevel     int             `json:"vipLevel"`
	TokenCode    string          `json:"tokenCode"`
	Source       string          `json:"source"`
}

//用户kyc同步请求
type KycSync struct {
	UserId    string `json:"userId"`
	KycStatus int    `json:"kycStatus"`
	Reson     string `json:"reson"`
}

//用户kyc同步
type KycSyncRsp struct {
	S2CCommon
	KycSync
}

type Dot struct {
	Token    string `json:"token"`
	CodeType int    `json:"codeType"` // 打点类型
}
