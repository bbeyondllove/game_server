package proto

import (
	"github.com/shopspring/decimal"
)

//创建角色
type C2SCreateRole struct {
	LocationID int    `json:"locationId"` //位置ID
	UserType   int    `json:"userType"`
	RoleId     int    `json:"roleId"`
	NickNname  string `json:"nickName"`
	Sex        int    `json:"sex"`
	Token      string `json:"token"`
}

type S2CCreateRole struct {
	S2CCommon
	CountryCode int    `json:"countryCode"`
	Email       string `json:"email"`
	HouseNum    int    `json:"houseNum"`
	Level       int    `json:"level"`
	Mobile      string `json:"mobile"`
	NickNname   string `json:"nickName"`
	X           int    `json:"X"`
	Y           int    `json:"Y"`
	RoleId      int    `json:"roleId"`
	LocationID  int    `json:"locationId"` //位置ID
	UserId      string `json:"userId"`
	UserType    int    `json:"userType"`
	Sex         int    `json:"sex"`
	Status      int    `json:"status"`
}

//位置更新
type C2SPositionChange struct {
	LocationID   int    `json:"locationId"`   //位置ID
	CurPosX      int    `json:"x"`            //玩家当前x坐标
	CurPosY      int    `json:"y"`            //玩家当前y坐标
	ScreenX      int    `json:"screenX"`      //玩家屏幕x坐标
	ScreenY      int    `json:"screenY"`      //玩家屏幕y坐标
	LeftTopX     int    `json:"leftTopX"`     //玩家屏幕左上角x坐标
	LeftTopY     int    `json:"leftTopY"`     //玩家屏幕左上角y坐标
	RightBottomX int    `json:"rightBottomX"` //玩家屏幕右下角x坐标
	RightBottomY int    `json:"rightBottomY"` //玩家屏幕右下角y坐标
	Type         int    `json:"type"`         //改变类型，0：正常改变；1：传送
	Token        string `json:"token"`
}

type S2CPositionChange struct {
	S2CCommon
	UserId   string `json:"userId"` //消息内容
	RoleId   int    `json:"roleId"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Type     int    `json:"type"` //改变类型，0：正常改变；1：传送
	NickName string `json:"nickName"`
	Sex      string `json:"sex"`
}

//位置获取
type C2SGetPosition struct {
	Token string `json:"token"`
}

type S2CGetPosition struct {
	S2CCommon
	UserId string `json:"userId"` //消息内容
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

//检测昵称
type C2SCheckNickName struct {
	NickName string `json:"nickName"`
	Token    string `json:"token"`
}

type UserPosition struct {
	RoleId   int    `json:"roleId"`
	NickName string `json:"nickName"`
	Sex      string `json:"sex"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

//进入城市地图请求
type C2SEnterCity struct {
	LocationID int    `json:"locationId"` //位置ID
	Token      string `json:"token"`
}

//进入城市地图响应
type S2CEnterCity struct {
	S2CCommon
	LocationID int `json:"locationId"` //位置ID

}

//获取同城玩家请求
type C2SGetCityUser struct {
	LocationID int    `json:"locationId"` //位置ID
	Token      string `json:"token"`
}

//获取同城玩家响应
type S2CGetCityUser struct {
	S2CCommon
	LocationID int                     `json:"locationId"` //位置ID
	EventInfo  map[int]*EventData      `json:"eventInfo"`
	UserList   map[string]UserPosition `json:"userPosition"`
}

//退出城市地图请求
type C2SQuitCity struct {
	LocationID int    `json:"locationId"` //位置ID
	Token      string `json:"token"`
}

//广播随机事件
type S2CBroadRandEvent struct {
	LocationID int                `json:"locationId"` //位置ID
	EventInfo  map[int]*EventData `json:"eventInfo"`  //key坐标序号，value(事件序号)
}

//处理随机事件
type C2SFinishEvent struct {
	LocationID   int    `json:"locationId"`   //位置ID
	ActivityType int    `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件
	X            int    `json:"x"`            //玩家当前x坐标
	Y            int    `json:"y"`            //玩家当前y坐标
	Token        string `json:"token"`
}

//随机事件响应
type S2CFinishEvent struct {
	S2CCommon
	LocationID   int                  `json:"locationId"`   //位置ID
	ActivityType int                  `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件
	X            int                  `json:"x"`            //玩家当前x坐标
	Y            int                  `json:"y"`            //玩家当前y坐标
	UserId       string               `json:"userId"`
	Cdt          decimal.Decimal      `json:"cdt"` //cdt
	ItemInfos    map[int][]*AwardItem `json:"itemInfo"`
}

//宝箱随机事件响应
type S2CTreasureBoxFinishEvent struct {
	S2CFinishEvent
	EventId int `json:"eventId"` //事件ID
}

//修改昵称
type C2SModifyNickName struct {
	Token      string `json:"token"`
	NickName   string `json:"nickName"`
	ModifyType int    `json:"modifyType"` //0：使用免费修改一次的方式 1：使用道具卡
	ItemId     int    `json:"itemId"`     //改名卡id，如果modifyType=0，ItemId忽略。
}
type S2CModifyNickName struct {
	S2CCommon
	NickName      string `json:"nickName"`
	ModifyNameNum int    `json:"modifyNameNum"` //改名次数
}

//实名认证
type C2SCertification struct {
	Token       string `json:"token"`
	Nationality string `json:"nationality"` //国籍
	FirstName   string `json:"firstName"`   //姓
	LastName    string `json:"lastName"`    //名
	IdType      int    `json:"idType"`      //证件类型，1：身份证；2：护照
	IdNumber    string `json:"idNumber"`    //证件号码
	ObjectKey   string `json:"objectKey"`   //文件对象的Key值，类似：aries/12345/temp.png
}

//
type CertificationPhoto struct {
	Front string `json:"front"`
	Back  string `json:"back"`  //国籍
	Other string `json:"other"` //姓
}

//可用角色列表
type C2SGetAllRole struct {
	Token string `json:"token"`
}

//可用角色列表响应
type S2CGetAllRole struct {
	S2CCommon
	//	Total     int         `json:"total"` //总数量
	//	ItemId    int         `json:"itemId" desc:"解锁卡ID"`
	RoleInfos []*RoleInfo `json:"list"  desc:"可用角色列表"` //此用户id对应用户推荐表里的邀请id
}

type RoleInfo struct {
	RoleId     int    `json:"roleId"`
	RoleName   string `json:"roleName"`
	Sex        int    `json:"sex"`
	State      int    `json:"state" desc:"0未解锁，1可解锁，2已解锁"`
	ItemId     int    `json:"itemId"`
	ExpireTime string `json:"expireTime"`
}

//解锁角色
type C2SAddRole struct {
	Token  string `json:"token"`
	RoleId int    `json:"roleId"`
	ItemId int    `json:"itemId"` //卡id
}

type S2CAddRole struct {
	S2CCommon
	RoleId     int    `json:"roleId"`
	ExpireTime string `json:"expireTime"`
}

//角色选择
type C2SSelectRole struct {
	UserId string `json:"userId"`
	RoleId int    `json:"roleId"`
	Sex    int    `json:"sex"  desc:"性别"`
}

//角色变化广播
type S2CRoleChange struct {
	UserID string `json:"userId"`
	RoleID int    `json:"roleId"` //角色ID
}

//获取好友信息
type C2SGetFirendInfo struct {
	SysType int    `json:"sysType"` //系统类型
	UserId  string `json:"userId"`  //好友用户ID
	Token   string `json:"token"`
}

//获取好友信息返回
type S2CGetFirendInfo struct {
	S2CCommon
	UserID      string `json:"userId"`
	NickName    string `json:"nickName"`
	CountryCode int    `json:"countryCode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	RegType     int    `json:"regType"`
	SysType     int    `json:"sysType"`
	InviteCode  string `json:"inviteCode"`
	VipLevel    int    `json:"vipLevel"`
	KycPassed   bool   `json:"kycPassed"`
	Status      int    `json:"status"`
	RegisterIp  string `json:"registerIp"`
	LoginIp     string `json:"loginIp"`
}

//获取好友的邀请关系
type C2SGetUserInvitation struct {
	SysType int    `json:"sysType"` //系统类型
	UserId  string `json:"userId"`  //好友用户ID
	Token   string `json:"token"`
}

//获取好友的邀请关系返回
type S2CUserInvitation struct {
	UserId      string `json:"userId"`      //用户ID，必选
	L1InviterId string `json:"l1InviterId"` //一级邀请人ID，必选
	L2InviterId string `json:"l2InviterId"` //二级邀请人ID，可选
	L3InviterId string `json:"l3InviterId"` //三级邀请人ID，可选
	InviteTime  string `json:"inviteTime"`  //邀请时间，必选
}
