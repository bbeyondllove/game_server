package model

import "time"

//用户表
type User struct {
	Id                     string    `xorm:"int(20) notnull pk" json:"id"`
	UserId                 string    `xorm:"int(20) notnull " json:"user_id"`
	UserType               int       `xorm:"int(1) notnull" json:"user_type"`
	SysType                int       `xorm:"int(1) notnull" json:"sys_type"`
	RoleId                 int       `xorm:"int(6) notnull" json:"role_id"`
	AvailableRoles         string    `xorm:"varchar(200)" json:"available_roles"`
	DeblockedRoles         string    `xorm:"varchar(255)" json:"deblocked_roles"`
	NickName               string    `xorm:"varchar(255)" json:"nick_name"`
	Sex                    int       `xorm:"int(1) notnull" json:"sex"`
	Level                  int       `xorm:"int(1) notnull" json:"level"`
	TopLevel               int       `xorm:"int(1) notnull" json:"top_level"`
	CountryCode            int       `xorm:"varchar(12)" json:"countryCode"`
	Mobile                 string    `xorm:"varchar(11)" json:"mobile"`
	Email                  string    `xorm:"varchar(64)" json:"email"`
	Token                  string    `xorm:"varchar(255) notnull" json:"token"`
	Status                 int       `xorm:"int(11) notnull" json:"status"`
	LocationId             int       `xorm:"int(20) notnull" json:"location_id"`
	PositionX              int       `xorm:"int(6) notnull" json:"position_x"`
	PositionY              int       `xorm:"int(6) notnull" json:"position_y"`
	HouseNum               int       `xorm:"int(6) notnull" json:"house_num"`
	InviteCode             string    `xorm:"varchar(20) notnull" json:"inviteCode"`
	InviterId              string    `xorm:"int(20) notnull" json:"inviter_id"`
	KycPassed              int       `xorm:"int(1) notnull" json:"kyc_passed"`
	KycStatus              int       `xorm:"int(1) notnull" json:"kyc_status"`
	LoginIp                string    `xorm:"varchar(30) notnull" json:"loginIp"`
	ModifyNameNum          int       `xorm:"int(10) notnull" json:"modify_name_num"`
	Point                  int       `xorm:"int(20) notnull" json:"point"`
	Cdt                    float32   `xorm:"decimal(12,4)" json:"cdt" desc:"用户的ctd"`
	TreasureBoxTotalIncome float32   `xorm:"decimal(12,4)" json:"treasure_box_total_income" desc:"用户的ctd"`
	Platform               int       `xorm:"int(10) notnull" json:"platform"`
	Version                string    `xorm:"-" json:"version"`
	CreateTime             time.Time `xorm:"timestamp notnull" json:"createTime"`
	UpdateTime             time.Time `xorm:"timestamp notnull" json:"updateTime"`
}

type GameWallet struct {
	Id              string    `xorm:"int(20) notnull pk" json:"id"`
	UserId          string    `xorm:"int(20) notnull pk" json:"user_id"`
	Amount          float64   `xorm:"Decimal notnull" json:"amount"`
	AmountAvailable float64   `xorm:"Decimal notnull" json:"amount_available"`
	AmountBlocked   float64   `xorm:"Decimal notnull" json:"amount_blocked"`
	CreateTime      time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime      time.Time `xorm:"timestamp notnull" json:"update_time"`
	TokenCode       string    `xorm:"int(20) notnull " json:"token_code"`
}

type ExchangeRecord struct {
	SysOrderSn      string    `xorm:"varchar(32)" json:"sys_order_sn" desc:"用户系统生成的流水号"`
	OrderSn         string    `xorm:"varchar(32) notnull" json:"order_sn" desc:"业务系统生成的流水号"`
	UserId          string    `xorm:"int(20) notnull" json:"user_id" `
	ExchangeType    string    `xorm:"varchar(20) notnull" json:"exchange_type" desc:"交易类型,charge充值， deduce扣款，"`
	CurrencyType    string    `xorm:"varchar(15) notnull" json:"currency_type" desc:"币种"`
	Amount          float64   `json:"amount" desc:"金额数量"`
	Status          int       `xorm:"int(20)" json:"status" desc:"加减状态，1加，2减"`
	UserAmount      float64   `xorm:"Decimal notnull" json:"user_amount" desc:"用户余额"`
	AmountAvailable float64   `xorm:"Decimal notnull" json:"amount_available" desc:"用户可用余额"`
	TargetAccount   int       `xorm:"int(20) notnull"  json:"target_account" desc:"目标账户"`
	Desc            string    `xorm:"varchar(255)" json:"desc" desc:"操作类型描述，如积分兑换"`
	AdminUser       string    `xorm:"varchar(50) notnull" json:"admin_user" desc:"操作用户名"`
	AdminUserId     int       `xorm:"int(20) notnull" json:"admin_user_id" desc:"操作用户ID"`
	CreateTime      time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime      time.Time `xorm:"timestamp notnull" json:"update_time"`
}

type UserLevelConfig struct {
	Id         int       `xorm:"int(20) notnull"  json:"id" desc:"ID"`
	LevelId    int       `xorm:"int(20) notnull"  json:"level_id" desc:"用户等级"`
	ItemId     int       `xorm:"int(20)" json:"item_id" desc:"道具卡id"`
	ItemType   int       `xorm:"int(20)" json:"item_type" desc:"道具卡类型id"`
	ItemName   string    `xorm:"int(20)" json:"item_name" desc:"道具卡名称"`
	ItemNum    int       `xorm:"int(20)" json:"item_num" desc:"道具卡数量"`
	CreateTime time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"update_time"`
}

type RoleInfo struct {
	Id         int       `xorm:"int(20) notnull"  json:"id" desc:"角色ID"`
	RoleName   string    `xorm:"varchar(255) notnull"  json:"role_name" desc:"用户等级"`
	Sex        int       `xorm:"int(1)" json:"sex" desc:"性别"`
	HeadId     string    `xorm:"varchar(255)" json:"head_id" desc:"头像ID"`
	State      int       `xorm:"-" json:"-"`
	CreateTime time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"update_time"`
}

type Certification struct {
	UserId      string    `xorm:"int(20) notnull" json:"user_id" `
	Nationality string    `xorm:"varchar(255)" json:"nationality"` //国籍
	FirstName   string    `xorm:"varchar(255)" json:"first_name"`
	LastName    string    `xorm:"varchar(255)" json:"last_name"`
	IdType      int       `xorm:"int(20)" json:"id_type"`
	IdNumber    string    `xorm:"int(20)" json:"id_number"`
	ObjectKey   string    `xorm:"varchar(255)" json:"objectKey"` //文件对象的Key值，类似：aries/12345/temp.png
	Suggestion  int       `xorm:"int(20)" json:"suggestion"`
	Reson       string    `xorm:"varchar(255)" json:"reson"`
	Status      int       `xorm:"int(20)" json:"status" desc:"0：正在审核，1：审核不通过，2：审核通过 3：再次提交审核"`
	CreateTime  time.Time `xorm:"timestamp notnull" json:"create_time"`
	UpdateTime  time.Time `xorm:"timestamp notnull" json:"update_time"`
}

type Email struct {
	Id           int64     `xorm:"int(20) autoincr pk" json:"id"`
	Display      int       `xorm:"int(1)" json:"-"`
	UserId       string    `xorm:"int(20)" json:"userId"`
	EmailType    int       `xorm:"int(2)" json:"emailType"`
	EmailTitle   string    `xorm:"varchar(64)" json:"emailTitle"`
	EmailContent string    `xorm:"text" json:"emailContent"`
	IsRead       int       `xorm:"int(1)" json:"isRead"`
	CreateTime   time.Time `xorm:"created" json:"createTime"`
	UpdateTime   time.Time `xorm:"updated" json:"-"`
	ExpireTime   time.Time `xorm:"timestamp" json:"expireTime"`
}

type CertificationRecord struct {
	Id              int64     `xorm:"int(20) autoincr pk" json:"id"`
	ExamineUserId   string    `xorm:"int(20) notnull comment('审核人员ID')" json:"examine_user_id" desc:"审核人员ID"`
	ExamineUserName string    `xorm:"varchar(255) notnull comment('审核人员名称')" json:"examine_user_name" desc:"审核人员名称"`
	UserId          string    `xorm:"int(20) notnull comment('用户ID')" json:"user_id" desc:"用户ID"`
	Status          int       `xorm:"int(20) comment('提交状态')" json:"status" desc:"提交状态,0：正在审核，1：审核不通过，2：审核通过 3：再次提交审核"`
	Suggestion      string    `xorm:"varchar(255) comment('审核原因')" json:"suggestion" desc:"审核原因"`
	ExamineStatus   int       `xorm:"int(20) comment('审核状态')" json:"examine_status" desc:"审核状态"`
	CreateTime      time.Time `xorm:"timestamp notnull comment('创建时间')" json:"create_time"`
	UpdateTime      time.Time `xorm:"timestamp notnull comment('更新时间')" json:"update_time"`
}
