

## 游戏接口文档



版本|时间|作者|备注
--|:--:|--:|--:
V0.1.0|2020.07.11|陈友能|1.0版本
V0.1.1|2020.08.12|陈友能|增加背包接口
V0.1.2|2020.08.14|陈友能|增加3个接口：进入城市、广播随机事件、处理随机事件
V0.1.3|2020.08.18|陈友能|统一命令规范
V0.1.4|2020.08.19|陈友能|修改背包协议；增加附录三<<道具类别表>>和附录四<<道具编号表>>
V0.15|2020.08.26|黄贺群|增加会员等级、被抢走战友列表、抢战友
V0.16|2020.08.28|陈友能|增加获取商品列表、购买商品列表、修改昵称接口
V1.1.0|2020.09.07|黄贺群|更改协议包头、包体设计，修复服务端协议校验问题
V1.1.1|2020.09.11|陈友能|对接新用户系统，修改相关协议，签名算法和业务逻辑
V1.2.0|2020.09.28|陈友能|增加实名认证接口
V1.2.1|2020.10.09|陈友能|增加商家模糊查询、好友详细信息、绑定邀请码、用户好友邀请收益等接口
V1.2.2|2020.11.05|陈友能|更新修改协议，并增加聚合支付接口
V1.3|2020.11.09|陈友能|整理协议，去掉多余的userID参数
V1.3|2020.11.21|匡顺辉|增加本地生活请求类型
V 1.3|2020.11.23|陈友能|增加签到、任务等协议接口
V 1.4|2020.12.13|陈友能|增加双旦活动功能设计和相关接口
V 1.4.1|2020.12.14|匡顺辉|增加双旦活动－排行榜消息类型
V 1.4.2|2021.01.07|陈友能|增加同屏玩家处理，优化消息广播
V 1.5|2021.01.08|丁院红，王身维|欢乐宝箱


### 1.通信协议：

json  

### 2.协议结构   
消息的包头和包体结构设计如下：
包头 + 包体
Msgtype（2个字节数字） + 字节数组
因为使用wss协议，没有设计长度字段和加密解密机制，包头用了2个字节数字表示命令字，包体采用json数据交换格式，经编码后转码为字节数组

### 3.协议内容

#### 3.1心跳包  

请求消息命令字：Msgtype =7（MSG_HEARTBEAT）
请求消息内容：

type C2SBase struct {
	Token string `json:"token"`
}

返回消息命令字：Msgtype =8（MSG_HEARTBEAT_RSP）
返回消息内容：

{"code":0,"message":""，"token":""}

#### 3.2挤用户下线

说明：这个消息是由服务端主动下发
请求消息命令字：Msgtype =4（MSG_LOGINANOTHER）
请求消息内容：无

#### 3.3断线重连

说明：登录后，若因网络状况差，服务端将websocket断开时，客户端发此消息进行重连
请求消息命令字：Msgtype =9（MSG_REBIND）
请求消息内容：

```golang
type C2SRebind struct {
	LocationID int32 `json:"locationId"` //位置ID
	Token  string `json:"token"`
}
```
返回消息命令字：Msgtype=10（MSG_REBIND_RSP）
返回消息内容：返回Code，Message

```
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```


#### 3.4发送验证码

请求消息命令字：Msgtype =21 （MSG_GET_VERIFICATION_CODE）
请求消息内容：

```golang
//用户获取验证码
type C2SGetVerificationCode struct {
	SysType     int    `json:"sysType"` //业务系统类型，0：Base系统；1：EChain商城系统；2：游戏系统；3：Aries交易所；											4：CCMYL交易所
	CountryCode int    `json:"countryCode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	UseFor      int    `json:"useFor"`   //1：注册；2：重置密码；3：绑定（手机或邮箱）
	CodeType    int    `json:"codeType"` //1：手机方式；2：邮箱方式
	Language    string `json:"language"` //下发短信的语言，默认为英文：en；简体中文：zh；繁体中文：tc
}

```

返回消息命令字：Msgtype =22（MSG_GET_VERIFICATION_CODE_RSP）
返回消息内容：返回Code，Message：

```
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```

#### 3.5 检测昵称是否存在

请求消息命令字：Msgtype =25（MSG_CHECK_NICK_NAME）
请求消息内容：

```golang
//检测昵称
type C2SCheckNickName struct {
	NickName string `json:" nickName"`
	Token  string `json:"token"`
}
```

返回消息命令字：Msgtype =26（MSG_CHECK_NICK_NAME_RSP）
返回消息内容：返回消息内容：返回Code，Message：

```
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```



#### 3.6注册

请求消息命令字：
Msgtype =1（MSG_REGISTER_PHONE，手机注册）
Msgtype =2（MSG_REGISTER_EMAIL，邮箱注册）

请求消息内容：

```golang
//用户邮箱注册
type C2SRegisterEmail struct {
	SysType   int    `json:"sysType"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	EmailCode string `json:"emailCode"`
	Inviter   string `json:"inviter"`  //邀请人，可以为邀请人的手机号、邮箱或邀请码
	NickNname string `json:"nickName"` //昵称，当不填时则设置为注册的邮箱
    Platform  string    `json:"platform"` //平台

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
    Platform  string    `json:"platform"` //平台
}
```

返回消息命令字：Msgtype =3（MSG_REGISTER_RSP）
返回消息内容：返回消息内容：返回Code，Message 例子：

```golang
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```

#### 3.7登录

请求消息命令字：Msgtype =5（MSG_LOGIN ）
请求消息内容：

```golang
//用户登录
type C2SLogin struct {
	SysType     int    `json:"sysType"`
	CodeType    int    `json:"codeType"` //1：手机方式；2：邮箱方式
	CountryCode int    `json:"countryCode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ClientIp    string `json:"clientIp"`
    Platform    string    `json:"platform"` //平台
	Version     string `json:"version"`  //平台版本
}
```

返回消息命令字：Msgtype =6（MSG_LOGIN_RSP）
返回消息内容：返回Code，Message等。 (返回状态status=2表示前端已有帐号登录，会把前端的帐号踢掉)例子：

```golang
{
	"code": 0,
	"message": "",
	"userId": 310190951450546176,
	"countryCode": 86,
	"email": "",
	"sex": 0,
	"inviteCode": "kcjluI",
	"level": 0,
	"mobile": "18899772602",
	"nickName": "18899772602",
	"regType": 0,
	"userType": 0,
	"roleId": 0,
	"status": 0,
	"sysType": 0,
	"rank": "乐于助人",
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTk4ODkzNjIsImlhdCI6MTU5OTgwMjk2MiwibmJmIjoxNTk5ODAyOTYyLCJ1c2VySWQiOjMxMDE5MDk1MTQ1MDU0NjE3Niwic3lzVHlwZSI6Mn0.bUQhCKyLM92L4rdWjUyAt1UrnTQaU0nxYfheX1IO5xY",
	"positsionX": 0,
	"positsionY": 0,
	"modifyNameNum": 1,
	"houseNum": 0,
	"locationId": 0,
	"inviterId": 0,//邀请人的用户ID，为0则表示没有邀请人
	"kycPassed": 0,// 是否已实名通过\n,0否，1是。,
	"KycStatus":0,//（0：正在审核，1：审核不通过，2：审核通过）
	"suggestion":1,//初始值;1:审核通过 2:姓名有误  3:证件号有误 4:重新提交证件正面照 5:重新提交证件反面照 6:重新提交手持证					件照
	"activityStatus":1//活动关闭或开启。0：关闭；1：开启
}
```

#### 3.8重置密码

请求消息令字：Msgtype =13（MSG_RESET_PASSWORD）
请求消息内容：

```golang
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
```

返回消息命令字：Msgtype =14（MSG_RESET_PASSWORD_RSP）
返回消息内容：返回Code，Message：

```
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```

#### 3.9创建角色	

请求消息命令字：Msgtype =15（MSG_CREATER_ROLE）
请求消息内容：

```golang
//创建角色
type C2SCreateRole struct {
	LocationID int32 `json:"locationId"` //位置ID
	UserType   int32 `json:"userType"`
	RoleId     int32 `json:"roleId"`
	NickNname  string `json:"nickName"`
	Sex        int32 `json:"sex"`
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =16（MSG_CREATER_ROLE_RSP）
返回消息内容：返回Code，Message等 。例子：

```golang
{
  "countryCode": "86",
  "email": "",
  "houseNum": 0,
  "level": 0,
  "mobile": "13434248318",
  "nickName": "昵",
  "x": "0",
  "y": "0",
  "roleId": 27,
  "locationId":3302,
  "userId":683624384,
  "userType":2,
  "code": 0,
  "message": ""
}
```

#### 3.10位置更新

请求消息命令字：Msgtype =17（MSG_POSITION_CHANGE）
请求消息内容：

```golang
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
```

无返回消息，但有广播消息。

广播消息类型：Msgtype =35（MSG_BROAD_POSITION）
广播消息内容：返回Code，Message ：

```golang
{
	"userId": "6462065254263259136",
	"type": 1,//改变类型，0：正常改变；1：传送
	"x": 36,
	"y": 20,
	"nickName": "kk",
	"sex": "1",
	"code": 0,
	"message": ""
}
```

#### 3.11位置获取

请求消息命令字：Msgtype =19（MSG_GET_POSITION）
请求消息内容：

```golang
//位置获取
type C2SGetPosition struct {
	Token  string `json:"token"`
}
```

返回消息命令字：Msgtype =20（MSG_GET_POSITION_RSP）
返回消息内容：返回Code，Message等：

```golang
{
  "x": "0",
  "y": "0",
  "userId": 10000002
  "code": 0,
  "message": ""
}
```

#### 3.12 获取用户游戏子钱包金额

请求消息命令字：Msgtype =27（MSG_GET_AMOUNT）
请求消息内容：

```golang
//获取用户钱包金额
    type C2SGetAmount struct {
	Token     string `json:"token"`
	TokenCode string `json:"tokenCode"`
}
```

返回消息命令字：Msgtype =28（MSG_GET_AMOUNT_RSP）
返回消息内容：返回Code，Message等：

```golang
{
  "userId": 10000002,
  "amount": 200.0,
  "amountAvailable": 200.0,
  "amountBlocked": 100.0,
  "cdt": 200.0,
  "code": 0,
  "message": "",
  "tokenCode":"SAR"
}
```

#### 3.13 用户背包和金额同步

请求消息命令字：Msgtype =39（MSG_UPDATE_ITEM_INFO）
请求消息内容：

```golang
type C2SUpdateItemInfo struct {
	Token  string `json:"token"`
}
```

返回消息命令字：Msgtype =40（MSG_UPDATE_ITEM_INFO_RSP）
返回消息内容：返回Code，Message等：

```golang
{
	"cdt": 200.0,
	"itemInfo": {
		"20": [{
				"itemId": 2000,
				"num": 1,
				"desc":"可解锁商人角色",
				"attr1":"1"
			}]
	},
  "code": 0,
  "message": ""
}
```

#### 3.14获取背包信息

请求消息命令字：Msgtype =29（MSG_GET_KNAPSACK）
请求消息内容：

```golang
//获取背包信息
type C2SGetKnapsack struct {
	Token  string `json:"token"`
}
```

返回消息命令字：Msgtype =30（MSG_GET_KNAPSACK_RSP）
返回消息内容：返回Code，Message等 。例子：

```golang
{
	"itemInfo": {
		"20": [{
				"itemId": 2000,
				"num": 1,
				"desc":"可解锁商人角色",
				"attr1":"1"
			}]
		},
  "userId": "1987978502921617408"
	"code": 0,
	"message": ""
}
```

#### 3.15 进入城市请求

请求消息命令字：Msgtype =31（MSG_ENTER_CITY）
请求消息内容：

```
//进入城市地图
type C2SEnterCity struct {
	LocationID int32 `json:"locationId"` //位置ID
	Token      string `json:"token"`
}
```

该消息触发进入和退出城市广播事件（Msgtype =82 MSG_BROAD_CITY_USER）。

同时返回消息命令字：Msgtype =32（MSG_ENTER_CITY_RSP）
返回消息内容：返回Code，Message和MsgData 。例子：

	{
		"code": 0,
		"message": "",
		"locationId": 3302
	}
#### 3.16 广播随机事件

下发消息类型：Msgtype =33（MSG_BROAD_RAND_EVENT）
下发消息内容：

```

type EventData struct {
	ActivityType int    `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件
	Type         string `json:"type"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
	EventId  int    //事件id
}

//广播随机事件
type S2CBroadRandEvent struct {
	LocationID int32            `json:"locationId"` //位置ID
	EventInfo  map[int]EventData `json:"eventInfo"`  //key坐标序号，value(事件序号)
}

```

#### 3.17 处理随机事件请求

请求消息命令字：Msgtype =37（MSG_FINISH_EVENT）
请求消息内容：

```
//处理随机事件
type C2SFinishEvent struct {
	ActivityType int    `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件
	LocationID int    `json:"locationId"` //位置ID
	X          int    `json:"x"`          //玩家当前x坐标
	Y          int    `json:"y"`          //玩家当前y坐标
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =38（MSG_FINISH_EVENT_RSP），并广播消息MSG_BROAD_FINISH_EVENT给同一城市其它玩家。
返回消息内容：返回Code，Message等：

```json
{
    "activityType":0,
	"location_id": 3302,
	"userIdY": "3222068591507636224",
	"X":38,
	"Y":2,
	"cdt": 200.0,
	"itemInfo": {
		"20": [{
				"itemId": 2000,
				"num": 1,
				"desc":"可解锁商人角色",
				"attr1":"1"
			}]
	},
	"code": 0,
	"message": ""
}
```

广播消息内容：返回Code，Message等：

```golang
{
	"activityType":0,
	"location_id": 3302,
	"userId": "3222068591507636224",
	"X":38,
	"Y":2,
	"cdt": 200.0,
	"itemInfo": {
		"20": [{
				"itemId": 2000,
				"num": 1,
				"desc":"可解锁商人角色",
				"attr1":"1"
			}]
	},
	"code": 0,
	"message": ""
}
```



#### 3.18 获取战友列表

请求消息命令字：Msgtype =41（MSG_GET_INVITE_USERS）
请求消息内容：

```golang
//获取战友列表
type C2SInviteUserList struct {
	SysType         int    `json:"sysType"`         //1商城 2游戏 3交易所
	InviterId       string `json:"inviterId"`       //此用户id对应用户推荐表里的邀请id
	InvitationLevel int    `json:"invitationLevel"` //第几级战友，取值只能是1， 2， 3
	VipLevel        int    `json:"vipLevel"`        //vip等级 0或不传 表示获取所有列表
	Page            int    `json:"page"`            //第几页
	Size            int    `json:"size"`            //每页多少数量
	Token      string `json:"token"`
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
	UserId          string  `json:"userId"`
	InvitationLevel int    `json:"invitationLevel"`
	NickName        string `json:"nickName"`
	InviteTime      string `json:"inviteTime"`
	Level           int    `json:"vipLevel"`
}
```

返回消息命令字：Msgtype =42（MSG_GET_INVITE_USERS_RSP）
返回消息内容：返回Code，Message等：

```golang
{
    "code":0,
    "msg":"success",
    "data":[
        {
            "userId":"323419383110045696",
            "nickName":"山有木兮卿ss",
            "invitationLevel":1,
            "vipLevel":0,
            "inviteTime":1602814723
        },
        {
            "userId":"323168219550453760",
            "nickName":"旧巷故人",
            "invitationLevel":1,
            "vipLevel":0,
            "inviteTime":1602754841
        },
        {
            "userId":"313352897725206528",
            "nickName":"非",
            "invitationLevel":1,
            "vipLevel":0,
            "inviteTime":1600414686
        },
        {
            "userId":"313349127192711168",
            "nickName":"等我长大就",
            "invitationLevel":1,
            "vipLevel":0,
            "inviteTime":1600413787
        },
        {
            "userId":"313347100505018368",
            "nickName":"18576613130",
            "invitationLevel":1,
            "vipLevel":0,
            "inviteTime":1600413303
        }
    ],
    "paging":{
        "totalSize":5,
        "currentSize":5,
        "currentPage":1
    }
}
```

#### 3.19查看玩家信息

请求消息命令字：Msgtype =11（MSG_GET_USERINFO）
请求消息内容：

```golang
//获取用户信息
type C2SGetUserInfo struct {
	UserId   string				`json:"userId"`			//战友ID
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =12（MSG_GET_USERINFO_RSP）
返回消息内容：返回Code，Message等：

```golang 
{
  "code": 0,
  "message": "OK",
  "userId": "2892836006682263552",
  "level": 0,
  "nickName": "goodboy",		
  "rank": "乐于助人",		//头衔
  "roleId":27
}
```

#### 3.20获取建筑简介

请求消息命令字：Msgtype =45（MSG_GET_BUILDING_DESC）
请求消息内容：

```golang
//获取建筑简介
type C2SGetBuildingInfo struct {
  LocationID int32 `json:"locationId"` //位置ID
	LocationX int32 `json:"x"`		//X坐标
	LocationY  int32  `json:"y"`	//Y坐标
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =46（MSG_GET_BUILDING_DESC_RSP）
返回消息内容：返回Code，Message等：

```golang
{
  "code": 0,
  "message": "OK",
  "desc": "腾讯大厦（Tencent Building）位于深圳市南山区高新科技园北区，深南大道北侧，其于2009年8月24日正式落成，是腾讯第一座自建写字楼。其楼体总高193米，地上39层，地下3层，建筑总面积88180.38㎡，办公面积69796.76㎡，是造型新颖、内部功能齐全、人文环境领先的超高层建筑，成为深圳特区闻名全国的深南。。。",
    "url": "https://www.qq.com",		
    "passportAviable": "SAN,ISR,SAR",		//可支持通证
    "imageUrl":"https://www.qq.com/image/1.jpg",
    "smallType":"A1",		//建筑类型
    "buildingName", "腾讯大厦",
    "buildingTypeName":"酒店",
  
}
```

#### 3.21获取会员等级体系数据

请求消息命令字：Msgtype =47（MSG_GET_MEMBER_SYS）
请求消息内容：

```golang
type C2SGetMemberSys struct {
	SysType int   `json:"sysType"`
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =48（MSG_GET_MEMBER_SYS_RSP）
返回消息内容：返回Code，Message等：

```golang
{
  "code": 0,
  "message": "OK",
  "data":[{
      "memberLevel": 1,			//会员等级(0:vip0 1:vip1 2:vip2 3:vip3 4:vip4 5:vip5)
      "requireSelfLocked": 1000,			//float64, 自购通证持有量
      "requireSan": 100, 	//float64, SAN缴纳量（伞下）
      "requireSanOnly":200,	//无自购通证要求时需要的SAN量
      "requireVip1": 1000,		//Vip1总量
      "requireVip2": 900,
      "requireVip3": 500,		//VIP3的人员数量
      "requireVip4": 200,
      "l1InviteVip1BonusPercent": 10，	//1级邀请(直推)为VIP1的返佣百分比
      "l1InviteVip2BonusPercent"10,		//1级邀请(直推)为VIP2的返佣百分比
      "l1InviteVip3BonusPercent":10,	//1级邀请(直推)为VIP3的返佣百分比
      "l1InviteVip4BonusPercent":10,	//1级邀请(直推)为VIP4的返佣百分比
      "l1InviteVip5BonusPercent":10,//1级邀请(直推)为VIP5的返佣百分比
      "l2InviteBonusPercent":10,	//2级邀请(间推)的返佣百分比
      "l3InviteBonusPercent":10,	//3级邀请(间推)的返佣百分比
      "itemList":["商人解锁卡", "改名卡"], 	//该等级享有的道具卡
      },
      {
      "memberLevel": 2,			//会员等级(0:vip0 1:vip1 2:vip2 3:vip3 4:vip4 5:vip5)
      "requireSelfLocked": 1000,			//float64, 自购通证持有量
      "requireSan": 100, 	//float64, 自购通证持有量
      "requireSanOnly":200,	//无自购通证要求时需要的SAN量
      "requireVip1": 1000,		//Vip1总量
      "requireVip2": 900,
      "requireVip3": 500,		//VIP3的人员数量
      "requireVip4": 200,
      "l1InviteVip1BonusPercent": 10，	//1级邀请(直推)为VIP1的返佣百分比
      "l1InviteVip2BonusPercent"10,		//1级邀请(直推)为VIP2的返佣百分比
      "l1InviteVip3BonusPercent":10,	//1级邀请(直推)为VIP3的返佣百分比
      "l1InviteVip4BonusPercent":10,	//1级邀请(直推)为VIP4的返佣百分比
      "l1InviteVip5BonusPercent":10,//1级邀请(直推)为VIP5的返佣百分比
      "l2InviteBonusPercent":10,	//2级邀请(间推)的返佣百分比
      "l3InviteBonusPercent":10,	//3级邀请(间推)的返佣百分比
      "itemList":["商人解锁卡", "改名卡"], 	//该等级享有的道具卡
  }]
  	
	}
```



#### 3.22获取用户当前等级状态

请求消息命令字：Msgtype =49（MSG_GET_USER_LEVEL）
请求消息内容：

```golang
type C2SGetUserLevel struct {
	SysType int   `json:"sysType"`
	Token      string `json:"token"`
}
```



返回消息命令字：Msgtype =50（MSG_GET_USER_LEVEL_RSP）
返回消息内容：返回Code，Message等：

```golang
{
  "code": 0,
  "message": "OK",
  "data": {
  	"level": 1,		//会员等级(0:vip0 1:vip1 2:vip2 3:vip3 4:vip4 5:vip5)
    "accumSelfLocked": 1000,			//float64, 累计自购通证持有量
		"accumSan": 100,				//float64, 累计SAN缴纳量（伞下）
		"vipStartDate"::1599962496,
		"vipEndDate":"",
    "vip1InviteeCount": 1000,
    "vip2InviteeCount": 900,
    "vip3InviteeCount": 500,
    "vip4InviteeCount": 200,
    "vip5InviteeCount": 100
  }
}
```



#### 3.23获取被抢走战友列表

说明：用户系统未提供接口，目前模拟数据
请求消息命令字：Msgtype =43（MSG_GET_GRAB_COMRADES）
请求消息内容：

```golang
type C2SGetGrabComrades struct {
	UserId  string `json:"userId"` 		//此用户id对应用户推荐表里的邀请id
	SysType int    `json:"sysType"`		//1商城 2游戏 3交易所
	Page    int    `json:"page"` 			//第几页
	Size    int    `json:"size"` 			//每页多少数量
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =44（MSG_GET_GRAB_COMRADES_RES）
返回消息内容：返回Code，Message等：

```golang
{
  "code": 0,
  "message": "OK",
  "total": 3,
  "list": [
      {
        "userId": "3222068591507636224",
        "nickName": "pj",
        "level":1,	//用户等级
      },
      {
        "userId": "5974519247143862272",
        "nickName": "",
        "level":1,
      },
      {
        "userId": "9167373859257483264",
        "nickName": "yyl",
        "level":1,
      }
    ]
}
```

#### 3.24抢战友

请求消息命令字：Msgtype =51（MSG_GRAB_COMRADE）
请求消息内容

```golang
type C2SGrabComrade struct {
	UserId   string				`json:"userId"`			//战友ID
	InviteId  string `json:"inviteId"` 		//抢劫者用户ID
	SysType    int    `json:"page"` 			//第几页
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =52（MSG_GRAB_COMRADE_RES）
返回消息内容：返回Code，Message 。例子

```golang
{
	"code": 0,
	"message": ""
}
```
#### 3.25玩家掉线广播

广播消息类型：Msgtype =36（MSG_BROAD_USER_OFFLINE）
广播消息内容

```golang
{
	"actionType":0,
	"382343346234": {
		"roleId": 27,
		"nickName": "风一样的",
		"sex": 0,
		"x": 36,
		"y": 20
	},
	"code": 0,
	"message": ""
}
```

#### 3.26获取商品列表

请求消息命令字：Msgtype =53（MSG_GET_ITEMS_LIST）
请求消息内容：	Token      string `json:"token"`。

返回消息命令字：Msgtype =54（MSG_GET_ITEMS_LIST_RSP）
返回消息内容：返回Code，Message和MsgData 。例子

```golang
{
	"itemList":{
    "20": [{
          "itemId": 2000,
          "itemType": 1,	//道具类型
          "itemName":"解锁卡",	//道具名称
          "isBind":0,	//是否绑定。0账户绑定、1可以交易
          "quality":,	//道具品质。0（其它类型）白色；1（稀有）金色；2（特级）紫色色；3（高级）蓝色；4（普通）绿色
          "isPile":0,	//是否堆叠（0：特殊物品；1不能堆叠;>1的数字可以堆叠）
          "getFrom":,		//获取途径。按功能分类表的编号，用逗号隔开。可为空
          "useJump":,		//使用跳转。按功能分类表的编号，点使用跳转到对应的功能界面
          "price":"0.48",			//售价
          "recommend":1,	
          "desc":"可解锁商人角色",
          "attr1":"0"
        }]
	} 
	"code": 0,
	"message": ""
}
```

#### 3.27购买商品

请求消息命令字：Msgtype =55（MSG_BUY_ITEM）
请求消息内容：

```golang
//购买商品
type C2SBuyItem struct {
	ItemId int    `json:"itemId"` //商品ID
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =56（MSG_BUY_ITEM_RSP）
返回消息内容：返回Code，Message和MsgData 。例子

```golang
{
	"itemId": {
		"id"       :, 
    "item_type"  :,
    "item_name"  :,
    "is_bind"    :,
    "quality"   :,
    "is_pile"    :,
    "get_from"   :,
    "use_jump"   :,
    "price"      :,
    "recommend"  :,
    "desc"       :,
    "create_time":,
    "update_time":,
	},
	"Money":"10.21"
	"code": 0,
	"message": ""
}
```


#### 3.28修改昵称

请求消息命令字：Msgtype =57（MSG_MODIFY_NICKNAME）
请求消息内容：

```golang
//修改昵称
type C2SModifyNickName struct {
	NickName   string `json:"nickName"`
	ModifyType int    `json:"modifyType"` //0：使用免费修改一次的方式 1：使用道具卡
	ItemId     int    `json:"itemId"`     //改名卡id，如果modifyType=0，ItemId忽略。
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =58（MSG_MODIFY_NICKNAME_RSP）
返回消息内容：返回Code，Message和MsgData 。例子

```golang
{
	"code": 0,
	"message": "",
	"nickName":"",
	"modifyNameNum":3,
}
```

#### 3.29可用角色列表

请求消息命令字：Msgtype =59（MSG_GET_ALL_ROLE）
请求消息内容：

```
type C2SGetAllRole struct {
	Token      string `json:"token"`
}
```



返回消息命令字：Msgtype =60（MSG_GET_ALL_ROLE_RSP）
返回消息内容：返回Code，Message和MsgData 。例子

```
{
	"code": 0,
	"message": "",
	"total": 3,
	"list": [{
    "roleId": 1,
    "roleName": "pj",
    "sex":1,	//性别:(0:男，1：女）
    "state":0, 	//0未解锁，1可解锁，2已解锁
    "itemId":2008,//解锁卡ID
    },
    {
    "roleId": 1,
    "roleName": "pj",
    "sex":1,	//性别:(0:男，1：女）
    "state":2,
    "itemId":2008,//解锁卡ID
    }]
}
```

#### 3.30给用户解锁角色

说明：
请求消息命令字：Msgtype =61（MSG_USER_ADD_ROLE）
请求消息内容：

```
//解锁角色
type C2SAddRole struct {
	RoleId     int    `json:"roleId"`   
 	ItemId     int    `json:"itemId"`     //改名卡id，如果modifyType=0，ItemId忽略。
 	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =62（MSG_USER_ADD_ROLE_RSP）
返回消息内容：返回Code，Message。例子

```go
//通用返回消息
type S2CAddRole struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
    RoleId  int `json:"roleId"`
    ExpireTime string `json:"expireTime"` // 过期时间
}
```

#### 3.31角色选择

请求消息命令字：Msgtype =63（MSG_SELECT_ROLE）
请求消息内容：

```
type C2SSelectRole struct {
	RoleId  int `json:"roleId"`
	Sex int `json:"sex"  desc:"性别"`
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =64（MSG_SELECT_ROLE_RSP）
返回消息内容：返回Code，Message。例子

```
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```

广播消息命令字：Msgtype =65（MSG_ROLE_CHANGE）

```
type S2CRoleChange struct {
   UserID    string  `json:"userId"`    
   RoleID  int `json:"roleId"` //角色ID
}
```

#### 3.32实名认证

请求消息命令字：Msgtype =66（MSG_CERTIFICATION）
请求消息内容：

```
//实名认证
type C2SCertification struct {
	Token       string `json:"token"`
	Nationality string `json:"nationality"` //国籍
	FirstName   string `json:"firstName"`   //姓
	LastName    string `json:"lastName"`    //名
	IdType      int    `json:"idType"`      //证件类型，1：身份证；2：护照
	IdNumber    int    `json:"idNumber"`    //证件号码
	ObjectKey   string `json:"objectKey"`   //文件对象的Key值，类似：aries/12345/temp.png
}
```

返回消息命令字：Msgtype =67（MSG_CERTIFICATION_RSP）
返回消息内容：返回Code，Message。例子

```
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```



#### 3.33绑定邀请码

请求消息命令字：Msgtype =69（MSG_BIND_INVITER）
请求消息内容：

```golang
//绑定邀请码
type C2SBindInviter struct {
	InviteCode string `json:"inviteCode"`
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =70（MSG_BIND_INVITER_RSP）
返回消息内容：返回Code，Message。例子

```golang
//通用返回消息
type S2CCommon struct {
    Code    int32  `json:"code"`    //错误代码
    Message string `json:"message"` //错误信息
}
```

#### 3.34 用户好友邀请收益接口（分红）

请求消息命令字：Msgtype =71（MSG_DEPOSIT_REBATE）
请求消息内容：

```golang
type C2SDepositRebate struct {
	Page    int   `json:"page"`            //第几页
	Size    int   `json:"size"`            //每页多少数量
	Token      string `json:"token"`、
}
```

返回消息命令字：Msgtype =72（MSG_DEPOSIT_REBATE_RSP）
返回消息内容：返回Code，Message，data，paging。例子

```golang
{
  "code": 0,
  "message": "OK",
  "data":[{
     		"userId":"3222068591507636224",
        "inviteIncome":"30.1",	//用户邀请总收益
        "friendNickName":"非",	//好友昵称
        "friendLevel":2,		//好友关系层级
        "vipLevel":1,			//vip等级
        "tokenCode":"SAN-BIT",		//货币名称
        "source":"",			//来源
      },
      {
      	"userId":"3222068591507636224",
        "inviteIncome":"31.1",
        "friendNickName":"rfk",
        "friendLevel":1,
        "vipLevel":2,
        "tokenCode":"SAN-BIT",
        "source":"",
  }],
  "paging": {
        "totalSize": 1,
        "currentSize": 1,
        "currentPage": 1
    }
}

```

#### 3.35 商家模糊查找接口

请求消息命令字：Msgtype =73（MSG_QUERY_SHOP）
请求消息内容：

```golang
type C2SQueryShop struct {
	KeyWord string `json:"keyWord"`
	Token   string `json:"token"`
}
```

返回消息命令字：Msgtype =74（MSG_QUERY_SHOP_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
	"code": 0,
	"message": "OK",
	"data": [{
		"Desc": "test",
		"Id": 2000,
		"BuildingName": "",
		"ImageUrl": "",
		"PassportAviable": 1,
		"PositionX": 0,
		"PositionY": 0,
		"ShopName": "shop1",
		"SmallType": "X1",
		"Url": ""
	}]
}

```

#### 3.36 获取好友详细信息接口

请求消息命令字：Msgtype =75（MSG_GET_FRIEND_INFO）
请求消息内容：

```golang
//获取好友信息
type C2SGetFirendInfo struct {
	SysType int    `json:"sysType"` //系统类型
	UserId  string `json:"userId"`
	Token   string `json:"token"`
}
```

返回消息命令字：Msgtype =76（MSG_GET_FRIEND_INFO_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
    "code":0,
    "message":"",
    "data":{
        "countryCode":86,
        "createdAt":1600156064,
        "email":"",
        "inviteCode":"gpqtwr",
        "inviterId":"0",
        "kycPassed":false,
        "lockTime":0,
        "loginIp":"",
        "loginTime":0,
        "mobile":"18576653139",
        "nickName":"静听年华s d搜",
        "regType":1,
        "registerIp":"",
        "status":0,
        "sysType":2,
        "thirdparty":[

        ],
        "userId":"312268161195970560",
        "vipLevel":1
    }
}

```

#### 3.37 获取好友的邀请关系

请求消息命令字：Msgtype =77（MSG_GET_INVITATION）
请求消息内容：

```golang
//获取好友的邀请关系
type C2SGetUserInvitation struct {
	SysType int    `json:"sysType"` //系统类型
	UserId  string `json:"userId"`  //好友用户ID
	Token   string `json:"token"`
}

```

返回消息命令字：Msgtype =78（MSG_GET_INVITATION_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
    "code":0,
    "message":"",
    "data":{
        "userId":"312268161195970560", //用户ID，必选
        "l1InviterId":"312268161195970561",//一级邀请人ID，必选
        "l2InviterId":"",				   //二级邀请人ID，可选
        "l3InviterId":"",				   //三级邀请人ID，可选
        "inviteTime":"2020-10-19 10:20:28"//邀请时间，必选
    }
}
```

#### 3.38 获取聚合支付链接

请求消息命令字：Msgtype =79（MSG_GET_TOEKN_PAY_URL）
请求消息内容：

```golang
//获取聚合支付接口
type C2SGetTokenPayUrl struct {
	Token string `json:"token"`
}
```

返回消息命令字：Msgtype =80（MSG_GET_TOEKN_PAY_URL_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
    "code": 200,
    "message": "操作成功",
	"sharePayUrl": "https://paytest.zifu.vip/payment/member/memberBaseRegister"  
}
```

#### 3.39退出城市

请求消息命令字：Msgtype =81（MSG_USER_QUT_CITY）
请求消息内容：

```golang
//退出城市地图请求
type C2SQuitCity struct {
	LocationID int    `json:"locationId"` //位置ID
	Token      string `json:"token"`
}
```

该消息触发进入和退出城市广播事件（Msgtype =82 MSG_BROAD_CITY_USER）。

#### 3.40进入和退出城市广播

广播消息命令字：Msgtype =82（MSG_BROAD_CITY_USER）
广播消息内容：返回Code，Message，和用户信息。例子

```golang
{
	"actionType":2,//1：进入城市，2：退出城市
	"382343346234": {
		"roleId": 27,
		"nickName": "风一样的",
		"sex": 0,
		"x": 36,
		"y": 20
	},
	"code": 0,
	"message": ""
}
```

#### 3.41 获取同城在线用户

请求消息命令字：Msgtype =83（MSG_GET_CITY_USER）
请求消息内容：

```golang
//获取同城玩家请求
type C2SGetCityUser struct {
	LocationID int    `json:"locationId"` //位置ID
	Token      string `json:"token"`
}
```

返回消息命令字：Msgtype =84（MSG_GET_CITY_USER_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
	"code": 0,
	"message": "",
	"locationId": 3302，
	"eventInfo":{
    "1":{
        "position":3,
        "eventId":6,
    }
	},
	"userList":{
		"138776789":{
			"nickName":"kk",
			"sex":"1",
			"x":38,
			"y":28,
		}
	}
	
}
```

#### 3.42 获取邮件列表

请求消息命令字：Msgtype =85（MSG_GET_EMAIL_LIST）
请求消息内容：

```golang
// 邮件列表
type C2SEmailList struct {
	IsRead int    `json:"isRead"`  // 是否已读  -1 全部数据  0 只返回未读   1 返回已读
	Token      string `json:"token"`
}
```
返回消息命令字：Msgtype =86（MSG_GET_EMAIL_LIST_RSP）
返回消息内容：返回Code，Message，data。例子

```json
// 邮件列表
{
	"code": 0,
	"message": "",
    "isPush": 1,  # 0 拉取列表  1 主动推送
	"readList": [{  // 已读列表
		"id": 2,  // 邮件ID
		"userId": "332551164648230912",  # 用户ID
		"emailType": 2,  # 邮件类型
		"emailTitle": "测试邮件2",   # 邮件标题
		"emailContent": "你好那你哈德杰拉德就Eads",  # 邮件内容
		"isRead": 1,  # 0 未读  1 已读
		"createTime": "2020-11-17T19:57:13+08:00",
		"expireTime": "2021/03/06 23:59:59",  # 过期时间
        "prizeList": [{
                "id": 1,
                "emailId": 44,
    			"prizeId":123, //奖品id
                "prizeType": 1,  // 奖品类型
                "przieName": "测试昵称", // 奖励名称
                "prizeNum": 1, // 奖励数量
                "prizeImg": "https://dss0.bdstatic.com/-0U0bnSm1A5BphGlnYG/tam-ogel/4b313f6304a4f683cca23f7426dc37f0_254_144.png",  // 奖品图片
                "isReceive": 0  // 0 未领取，1. 已领取， 2.领取失败
            }]
	}],
	"unreadList": [{  // 未读列表
		"id": 1,
		"isRead": "332551164648230912",
		"emailType": 1,
		"emailTitle": "测试邮件",
		"emailContent": "测试邮件内容",
		"isRead": 0,
		"createTime": "2020-11-17T19:56:39+08:00",
        "prizeList": []
	}]
}
```
#### 3.43 删除邮件

请求消息命令字：Msgtype =87（MSG_DEL_EMAIL）
请求消息内容：

```go
type C2SDelEmail struct {
	Token    string `json:"token"`
	EmailIds []int  `json:"emailIds"`  // 邮件ID
}
```

返回消息命令字：Msgtype =88（MSG_DEL_EMAIL_RSP）
返回消息内容：返回Code，Message，data。例子

```json
{"code":0,"message":""}
```



#### 3.44 设置邮件为已读

请求消息命令字：Msgtype =89（MSG_SET_EMAIL_READ）
请求消息内容：

```go
// 设置邮件已读
type C2SSetEmailRead struct {
	Token    string `json:"token"`
	EmailIds []int  `json:"emailIds"`  // 邮件ID
}
```

返回消息命令字：Msgtype =90（MSG_SET_EMAIL_READ_RSP）
返回消息内容：返回Code，Message，data。例子

```json
{"code":0,"message":""}
```



#### 3.45 获取邮件数量

请求消息命令字：Msgtype =91（MSG_COUNT_EMAIL）
请求消息内容：

```
// 获取邮件数量
type C2SGetEmialCount struct {
	Token string `json:"token"`
}
```



返回消息命令字：Msgtype =92（MSG_COUNT_EMAIL_RSP）
返回消息内容：返回Code，Message，data。例子

```json
{
	"code": 0,
	"message": "",
	"readNum": 1, // 已读的邮件数量
	"unreadNum": 1  // 未读的邮件数量

}
```



#### 3.46 发送实名邮件 （服务端推送）

返回消息命令字：Msgtype =94（MSG_PUSH_USER_MEIAL_RSP）
返回消息内容：返回Code，Message，data。例子

```json
{
	"code": 0,
	"message": "",
    "isPush": 1,  //  1 推送
	"newEmailList": [{
		"id": 1,  // id
		"userId": "332551164648230912",  // 用户id
		"emailType": 1,  // 邮件类型 1文本邮件
		"emailTitle": "测试邮件",  // 邮件标题
		"emailContent": "测试邮件内容", // 邮件内容
		"isRead": 1, // 0 未读， 1已读
		"createTime": "2020/11/17 19:56:39"  // 邮件发送时间
	}]
}
```

#### 3.47 获取签到列表

请求消息命令字：Msgtype =201（MSG_GET_SIGNIN_LIST）
请求消息内容：

```golang
type C2SBase struct {
	Token string `json:"token"`
}
```

返回消息命令字：Msgtype =202（MSG_GET_SIGNIN_LIST_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
	"code": 0,
	"message": "成功",
	"signinDayNo": 3, //今天要签第几天
	"CurDayStatus": false, //今天是否已完成(false：未完成；true:已完成)
	"CurWeekStatus": false,//本周是否已补签(false：未完成；true:已完成)	
	"taskList": [{
		"id": 30001,
		"taskType": 3,
		"title": "第1阶段第1天签到",
		"desc": "签到任务",
		"status": 0,
		"taskKey": "task_signin",
		"taskValue": 0,
		"messageId": "0",
		"eventId": 0,
		"frontTaskId": 0,
		"awards": {
			"ItemId": 4001,
			"ItemNum": 1000,
			"IsGift": false,
			"imgUrl": "http://zifu-admin-client.oss-cn-shenzhen.aliyuncs.com/test/1606463793580.png",
			"desc": "可用于物品兑换",
			"attr1": "0",
			"itemName": "普通碎片",
			"sex": 1,
			"awardList": null
		}
	}, {
		"id": 30002,
		"taskType": 3,
		"title": "第1阶段第2天签到",
		"desc": "签到任务",
		"status": 0,
		"taskKey": "task_signin",
		"taskValue": 0,
		"messageId": "0",
		"eventId": 0,
		"frontTaskId": 0,
		"awards": {
			"ItemId": 4001,
			"ItemNum": 1000,
			"IsGift": false,
			"imgUrl": "http://zifu-admin-client.oss-cn-shenzhen.aliyuncs.com/test/1606463793580.png",
			"desc": "可用于物品兑换",
			"attr1": "0",
			"itemName": "普通碎片",
			"sex": 1,
			"awardList": null
		}

	}]
}
```

#### 3.48 签到

请求消息命令字：Msgtype =203（MSG_SIGN_IN）
请求消息内容：

```golang
//签到请求
type C2SSignIn struct {
	Token      string `json:"token"`
	SigninType int    `json:"signinType"` //1：正常签到；2：补签
	ItemId     int    `json:"itemId"`     //补签卡id，SigninType=1，ItemId忽略。SigninType=2时，ItemId为补签卡id
}
```

返回消息命令字：Msgtype =204（MSG_SIGN_IN_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
	"code": 0,
	"message": "",
	"awards":[
		{"itemId": 2001,"itemNum": 1},
		{"itemId": 2002,"itemNum": 1}
		]	
}
```

#### 3.49 获取任务列表

请求消息命令字：Msgtype =205（MSG_GET_TASK_LIST）
请求消息内容：

```golang
type C2SBase struct {
	Token string `json:"token"`
}
```

返回消息命令字：Msgtype =206（MSG_GET_TASK_LIST_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
	"code": 0,
	"message": "成功",
	"signinDayNo": 0,
	"taskList": [{
		"id": 20001,
		"taskType": 2,
		"title": "今日完成了每日签到",
		"desc": "今日在游戏内签到成功",
		"status": 0,
		"taskKey": "task_daily_signin",
		"taskValue": 1,
		"messageId": "203",
		"eventId": 0,
		"frontTaskId": 0,
		"awards": {
			"ItemId": 4001,
			"ItemNum": 1000,
			"IsGift": false,
			"imgUrl": "http://zifu-admin-client.oss-cn-shenzhen.aliyuncs.com/test/1606463793580.png",
			"desc": "可用于物品兑换",
			"attr1": "0",
			"itemName": "普通碎片",
			"sex": 1,
			"awardList": null
		}
	}, {
		"id": 20002,
		"taskType": 2,
		"title": "今日清理一次垃圾",
		"desc": "今日清理一次垃圾",
		"status": 0,
		"taskKey": "task_daily_rubbish",
		"taskValue": 1,
		"messageId": "37",
		"eventId": 0,
		"frontTaskId": 0,
		"awards": {
			"ItemId": 4001,
			"ItemNum": 1000,
			"IsGift": false,
			"imgUrl": "http://zifu-admin-client.oss-cn-shenzhen.aliyuncs.com/test/1606463793580.png",
			"desc": "可用于物品兑换",
			"attr1": "0",
			"itemName": "普通碎片",
			"sex": 1,
			"awardList": null
		}

	}]
}
```

#### 3.50 领取任务奖励

请求消息命令字：Msgtype =207（MSG_GET_TASK_AWARD）
请求消息内容：

```golang
//任务奖励领取请求
type C2STaskAward struct {
	Token  string `json:"token"`
	TaskId int    `json:"task_id"`
}
```

返回消息命令字：：Msgtype =208 (MSG_GET_TASK_AWARD_RSP）
返回消息内容：返回Code，Message，data。例子

```golang
{
	"code": 0,
	"message": "",
	"awards":[
		{"itemId": 2001,"itemNum": 1},
		{"itemId": 2002,"itemNum": 1}
		]	
}
```

#### 3.51 商家跳转请求

请求消息命令字：Msgtype =209（MSG_ENTER_SHOP）
请求消息内容：

```golang
type C2SBase struct {
	Token string `json:"token"`
}
```

无返回消息。

#### 3.52 任务完成通知

任务完成通知命令字：Msgtype =210（MSG_BROAD_USER_TASK）
广播消息内容：返回Code，Message，和用户信息。例子

```golang
{
	"code": 0,
	"message": "",
	"task_id":"20004",
	"awards":[
		{"itemId": 2001,"itemNum": 1},
		{"itemId": 2002,"itemNum": 1}
		]	
}
```

#### 3.53 推送在线公告 （服务端推送）

返回消息命令字：Msgtype =96（MSG_PUSH_NOTICE_RSP）
返回消息内容：返回Code，Message，data。例子

```json
{
	"code": 0,  // 0 成功
	"message": "",
	"upgradeNotice": { //升级通知
		"noticeType": 1, // 公告类型,0:版本更新公告,1:普通公告,2:强制更新公告,3:更新前公告
		"noticeTitle": "test123", // 公告标题
		"noticeContent": "test123", // 公告内容
        "noticeUrl":"", // 公告URL
		"version": "1.0.1" // 版本号
	}
}
```

#### 3.54 获取更新公告

请求消息命令字：Msgtype =97（MSG_GET_UPGRADE_NOTICE ）
请求消息内容：

```go
type C2SCheckUpgradeNotice struct {
	Version int    `json:"version"`  // 版本号
}
```

返回消息命令字：Msgtype =96（MSG_PUSH_NOTICE_RSP）
返回消息内容：返回Code，Message，data。例子

```json
{
	"code": 0,  // 0 成功
	"message": "",
	"upgradeNotice": [{ //升级通知
		"noticeType": 1, // 公告类型
		"noticeTitle": "test123", // 公告标题
		"noticeContent": "test123", // 公告内容
		"version": "1.0.1" // 版本号
	}]
}
```

#### 3.55 城市中心 icon 入口统计

请求消息命令字：Msgtype =301（MSG_STATISTICS_CITY_ICON）
请求消息内容：

```go
type C2SBase struct {
  Token string `json:"token"`
}
```

#### 3.56同步用户kyc状态(推送到前端)

消息命令字：Msgtype =98（MSG_SEND_KYC_STATUS）
广播消息内容：返回Code，Message，和用户信息。例子

```golang
{
	"userId":"312268161195970560", 
	"KycStatus":1,
	"Reson":"",
	"code": 0,
	"message": ""
}
```

#### 3.57更新任务信息（后端推送给前端）

消息命令字：Msgtype =99（MSG_UPDATE_TASK_INFO）
广播消息内容：空值

{
	"status":1,//活动关闭或开启。0：关闭；1：开启
	"code": 0,
	"message": ""
}

#### 3.6. 双旦活动

##### 3.6.1 每日排行榜.

请求消息命令字：Msgtype =405（MSG_RANK_LIST_DAY）

请求消息内容：

```
type RankListDay struct {
	token string `json:"token"` // 登入token.
}
```

返回内容：

```
{
    "code":0,
    "msg":"success",
    "data":{
        "frontUserPoint":5, // 距离前一名分数差
        "rankList":[ // 排行榜列表,已经按分数由高到低排好序了
            {
                "award":200, // 奖励
                "name":"erying@qq.com", // 用户昵称
                "rank":1, // 名次
                "scores":810 // 当天双旦值
            },
            {
                "award":190,
                "name":"jiujiu@qq.com",
                "rank":2,
                "scores":610
            }
        ],
        "userPoint":65 // 我的分数
        "userRank":5 // 我的排名, 如果值为0，显示“暂无排名”
    }
}
```

##### 3.6.2 总排行榜.

请求消息命令字：Msgtype =407（MSG_RANK_LIST_ALL）

请求消息内容：

```
type RankListDay struct {
	token string `json:"token"` // 登入token.
}
```

返回内容：

```
{
    "code":0,
    "msg":"success",
    "data":{
        "frontUserPoint":5, // 距离前一名分数差
        "rankList":[ // 排行榜列表,已经按分数由高到低排好序了
            {
                "award":200, // 奖励
                "name":"erying@qq.com", // 用户昵称
                "rank":1, // 排名
                "scores":810 // 当天双旦值
            },
            {
                "award":190,
                "name":"jiujiu@qq.com",
                "rank":2,
                "scores":610
            }
        ],
        "userPoint":65 // 我的分数
        "userRank":5 // 我的排名, 如果值为0，显示“暂无排名”
    }
}
```

##### 3.6.3 每日双旦值记录.

请求消息命令字：Msgtype =409（MSG_RANK_LIST_DAY_PROP）

请求消息内容:

```
type RankListDay struct {
	token string `json:"token"` // 登入token.
}
```

返回值：

```
{
    "code":0,
    "msg":"success",
    "data":{
        "list":[ // 双旦值记录列表
            {
                "award":200, // 可获取奖励的双旦值
                "completedNumber":0, // 完成一次奖励中已经达到的次数，这个字段配合requireNumber字段一起使用
                "name":"成功邀请1名好友", // 道具名称
                "point":400, // 已加分数
                "requireNumber":20, // 完成一次奖励需要的次数
                "used":2 // 已完成奖励次数
                "isInviteFriend": false, // 是否为邀请好友类型
            },
            {
                "award":10,
                "completedNumber":16,
                "name":"收集圣诞袜子数量达到20个",
                "point":10,
                "requireNumber":20,
                "used":1
                "isInviteFriend": false, // 是否为邀请好友类型
            },
            {
                "award":5,
                "completedNumber":13,
                "name":"收集圣诞帽子数量达到20个",
                "point":5,
                "requireNumber":20,
                "used":1
                "isInviteFriend": false, // 是否为邀请好友类型
            }
        ]
    }
}
```

##### 3.6.4 总双旦值记录.

请求消息命令字：Msgtype =411（MSG_RANK_LIST_ALL_PROP）

请求消息内容:

```
type RankListDay struct {
	token string `json:"token"` // 登入token.
}
```

返回值：

```
{
    "code":0,
    "msg":"success",
    "data":{
        "list":[ // 双旦值记录列表
            {
                "award":200, // 可获取奖励的双旦值
                "completedNumber":0, // 完成一次奖励中已经达到的次数，这个字段配合requireNumber字段一起使用
                "name":"成功邀请1名好友", // 道具名称
                "point":400, // 已加分数
                "requireNumber":20, // 完成一次奖励需要的次数
                "used":2 // 已完成奖励次数
                "isInviteFriend": false, // 是否为邀请好友类型
            },
            {
                "award":10,
                "completedNumber":16,
                "name":"收集圣诞袜子数量达到20个",
                "point":10,
                "requireNumber":20,
                "used":1
                "isInviteFriend": false, // 是否为邀请好友类型
            },
            {
                "award":5,
                "completedNumber":13,
                "name":"收集圣诞帽子数量达到20个",
                "point":5,
                "requireNumber":20,
                "used":1
                "isInviteFriend": false, // 是否为邀请好友类型
            }
        ]
    }
}
```

##### 3.6.5 获取用户双蛋圣诞树，糖果.

请求消息命令字：Msgtype =420（MSG_GET_SWEET_AND_TREE）
请求消息内容：

```golang
type C2SBase struct {
	Token  string `json:"token"`
}
```

返回消息命令字：Msgtype =421（MSG_GET_SWEET_AND_TREE_RSP）
返回消息内容：返回Code，Message等 。例子：

```golang
{
	"sweetTree": {
		"60": [{
				"itemId": 6006,
				"num": 10,
				"desc":"圣诞糖果",
				"attr1":"0"
				},
               {
				"itemId": 6007,
				"num": 20,
				"desc":"圣诞树",
				"attr1":"0"
				}
              ]
		},
    "userDayCdt":0.10,//当天当前用户总共兑换的cdt
    "totalDayCdt":0.10,//当天所有用户总共兑换的cdt
    "currentCDT":0.10,//当前用户的cdt
 	"userId": "1987978502921617408"
	"code": 0,
	"message": ""
}
```

##### 3.6.6 获取用户双蛋碎片.

请求消息命令字：Msgtype =422（MSG_GET_PATCH）
请求消息内容：

```golang
type C2SBase struct {
	Token  string `json:"token"`
}
```

返回消息命令字：Msgtype =423（MSG_GET_PATCH_RSP）
返回消息内容：返回Code，Message等 。例子：

```golang
{
    "currentNum":52,//当前数量
    "tradeNeedNum":1000,//兑换圣诞老人所需数量        
    "userId": "1987978502921617408",
	"code": 0,
	"message": ""
}
```



##### 3.6.7  兑换CDT

请求消息命令字：Msgtype =401（MSG_SWEET_TREE）
请求消息内容：

```golang
type C2SBase struct {
	Token  string `json:"token"`
}
```

返回消息命令字：Msgtype =402（MSG_SWEET_TREE_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
  	"userId": "1987978502921617408",
	"code": 0,
	"message": "",
	"changeCdt": 0.2 // 用户兑换的CDT
}
```

##### 3.6.8 获取邀请列表

请求消息命令字：Msgtype =413（MSG_RANK_LIST_INVITE_RECORD）
请求消息内容：

```go
type C2SBase struct {
	Token  string `json:"token"`
    IsDay  int  `json:"isDay"`  // 0总的邀请列表   1 每天的邀请列表
}
```

返回消息命令字：Msgtype =414（MSG_RANK_LIST_INVITE_RECORD_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "createTime": "2020/12/12",
        "nickName": "温2323"
      }
    ]
  }
}
```

##### 3.6.9  圣诞老人角色过期通知

返回消息命令字：Msgtype =451（MSG_SANTA_LAUS_ROLE_EXPIRE_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
	"userId":"123456", # 用户ID
    "roleId":100
}
```

##### 3.6.10  获取活动状态

请求消息命令字：Msgtype =415（MSG_GET_DOUBLE_YEAR_STATUS）
请求消息内容：

```go
type C2SGetActiveStatus struct {
	Token  string `json:"token"`
    StatusType   int `json:"statusType"` //0:双旦;1:排行榜;
}
```

返回消息命令字：Msgtype =416（MSG_GET_DOUBLE_YEAR_STATUS_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
    "statusType":0, //0:双旦;1:排行榜;
    "status":1 //0：双蛋活动未开始 1:双蛋活动进行中 2:双蛋活动已结束
}
```

##### 3.6.11  圣诞老人解锁卡发放通知

返回消息命令字：Msgtype =417（MSG_BROAD_SANTA_CARD）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
	"content":"恭喜xxx集齐圣诞老人碎片，活动过后即可使用圣诞老人角色" 
}
```

##### 3.6.12  双旦活动结束通知

返回消息命令字：Msgtype =418（MSG_BROAD_ACTIVITY_END）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
	"activityType":1 //活动ID 0:普通随机事件；1：双蛋随机事件
}
```

#### 

#### 3.70 领取邮件奖励

请求消息命令字：Msgtype =211（MSG_GET_INVITATION_CODE）
请求消息内容：

```go
type C2SSetReceiveRewards struct {
	Token   string `json:"token"`
	EmailId int    `json:"emailId"`  // 邮件ID
}
```

返回消息命令字：Msgtype =212（MSG_GET_INVITATION_CODE_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功"
}
错误码：code 
5028  // 邮件不存在
5029  // 领取奖励失败
5030  // 奖励已领取
5033  // 奖品已过期
```

#### 3.71 获取入驻商户URL

请求消息命令字：Msgtype =213（MSG_GET_MERCHANTS_URL）
请求消息内容：

```go
type C2SGetMerchantUrl struct {
	Token string `json:"token"`
}
```

返回消息命令字：Msgtype =214（MSG_GET_MERCHANTS_URL_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
	"merchantEnteringUrl": "https://base.zifu.vip/m/#/settlements/list",
	"activityPromotionUrl": "https://base.zifu.vip/m/#/extensionCode"
}
```

#### 3.72 宝箱活动

##### 3.72.1 打开宝箱

请求消息命令字：Msgtype =455（MSG_OPEN_STREASURE_BOX）
请求消息内容：

```go
type C2SOpenStreasureBox struct {
	Token        string `json:"token"`
	ActivityType int    `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件 2.宝箱
	LocationID   int    `json:"locationId"`   //位置ID
	X            int    `json:"x"`            //玩家当前x坐标
	Y            int    `json:"y"`            //玩家当前y坐标
}
```

返回消息命令字：Msgtype =456（MSG_OPEN_STREASURE_BOX_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
  "code": 0,
  "message": "成功",
  "dayNum":0,  //每天打开的次数
  "residueDegree": 100 // 剩余次数
}
```

##### 3.72.2完成宝箱

请求消息命令字：Msgtype =457（MSG_FINISH_STREASURE_BOX）
请求消息内容：

```go
//完成宝箱
type C2SFinishStreasureBox struct {
	Token string `json:"token"`
	FinishStatus int    `json:"finishStatus"` //完成状态 0 未完成， 1 已完成
	ActivityType int `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件 2.宝箱
	LocationID   int `json:"locationId"`   //位置ID
	X            int `json:"x"`            //玩家当前x坐标
	Y            int `json:"y"`            //玩家当前y坐标
}
```

返回消息命令字：Msgtype =458（MSG_FINISH_STREASURE_BOX_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
	"locationId": 3302,  // 位置ID
	"activityType": 2, //活动ID 0:普通随机事件；1：双蛋随机事件 2.宝箱
	"x": 23, //玩家当前x坐标
	"y": 9,  //玩家当前y坐标
	"userId": "335004603600867328", // 用户ID
	"boxAward": { // 箱子奖励信息
		"awardId": 900023,  //领取奖励的ID
		"itemId": 0,   
		"itemNum": 0.05, // 随机生成的奖励
		"imgUrl": "https://zifu-admin-client.oss-cn-shenzhen.aliyuncs.com/test/common/1608519513584.png", // 奖励图片
		"desc": "CDT",
		"itemName": "CDT", //奖励名称
	}
}

错误码:
5001 非法请求
5010 目标不存在
5002 系统异常
```

##### 3.72.3 领取奖励

请求消息命令字：Msgtype =459（MSG_RECEIVE_STREASURE_BOX）
请求消息内容：

```go
type C2SReceiveBoxReward struct {
	Token   string `json:"token"`  // 用户token
	AwardId int    `json:"awardId"`  // 箱子的奖励ID
}
```

返回消息命令字：Msgtype =460（MSG_RECEIVE_STREASURE_BOX_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
	"ctd": 1.5  // 用户的总CDT
}
```



##### 3.72.4 获取宝箱奖励记录 

请求消息命令字：Msgtype =461（MSG_GET_STREASURE_BOX_RECORD）
请求消息内容：

```go
type C2SStreasureBoxGetRecord struct {
	Token string `json:"token"` //用户token
}
```

返回消息命令字：Msgtype =462（MSG_GET_STREASURE_BOX_RECORD_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,  
	"message": "",
	"data": [{
		"openTime": 1610075594, // 打开时间
		"cdt": "0.1",  // cdt
		"watchTime": 0  // 观看时长
	}, {
		"openTime": 1610884101,
		"cdt": "0.01",
		"watchTime": 0
	}],
	"totalCdt": "0.12",  //宝箱获取的总cdt
	"watchNum": 11  // 总记录
}
```

##### 3.72.5清理过期宝箱 (后端推送给前端)

返回消息命令字：Msgtype =466（MSG_BROAD_CLEAN_EVENT）
返回消息内容：返回Code，Message等 。例子：

主动推送过期宝箱: Msgtype =466（MSG_BROAD_CLEAN_EVENT）

```json
{
    "locationId": 3302,
    "eventInfo": {
        "700001": {
            "activityType": 2,
            "type": "",
            "x": 1,
            "y": 7,
            "EventId": 7001
        },
        "700025": {
            "activityType": 2,
            "type": "",
            "x": 25,
            "y": 7,
            "EventId": 7001
        }
    }
}
```

##### 3.72.6 获取打开宝箱的次数

返回消息命令字：Msgtype =350（MSG_TREASURE_BOX_DAY_NUM）
返回消息内容：

```go
type C2SStreasureBoxDayNum struct {
	Token string `json:"token"`
}

```

返回消息命令字：Msgtype =351（MSG_TREASURE_BOX_DAY_NUM_RSP）

返回Code，Message等 。例子：

```json
{
  "code": 0,
  "message": "成功",
  "dayNum":0,  //每天打开的次数
  "residueDegree": 100 // 剩余次数
}
```



#### 3.73 获取活动状态

请求消息命令字：Msgtype =463（MSG_ACTIVITY_STATYS）
请求消息内容：

```go
type C2SActivitStatus struct {
	Token string `json:"token"`
}

```

返回消息命令字：Msgtype =464（MSG_ACTIVITY_STATYS_RSP）
返回消息内容：返回Code，Message等 。例子：

主动推送活动状态: Msgtype =465（MSG_PUSH_ACTIVITY_STATUS_RSP）

```json
{
	"code": 0,
	"message": "",
	"activity": {
		"treasureBox": 1, // 欢乐宝箱状态  1 打开，0 关闭
		"randList": 1,  // 排行榜 0未开始  1 进行中  2 活动已结束
		"exchangeCdt": 1, // 兑换CDT 0未开始  1 进行中  2 活动已结束
		"exchangeRole": 1 // 兑换角色 0未开始  1 进行中  2 活动已结束
	}
}
```

#### 3.74 获取url 配置

请求消息命令字：Msgtype =215（MSG_GET_CONFIGURE_URL）
请求消息内容：

```go
type C2SConfigureUrl struct {
	Token string `json:"token"`
}

```

返回消息命令字：Msgtype =216（MSG_GET_CONFIGURE_URL_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功",
	"data": {
		"activityPromotionUrl": "http://testbase.zifu.vip/h5/#/extensionCode ",
		"aggregatePay": "http://testbase.zifu.vip/h5/#/shared",
		"ariesLink": "https://www.ccmyl.vip/m/",
		"assetsList": "https://base.zifu.vip/h5/#/",
		"deFiLink": "https://ccmdefi.com/",
		"echainLink": "https://www.echainbuy.com/h5/home/home.html",
		"echainLinkH5": "https://www.echainbuy.com/h5/index/index.html",
		"inviteLink": "https://zifu.vip/web-advers/index.html",
		"localLife": "http://106.53.252.137/#/home",
		"merchantEnteringUrl": "http://testbase.zifu.vip/h5/#/settlements/list"
	}
}
```

#### 3.75 统计宝箱

请求消息命令字：Msgtype =302（MSG_GET_CONFIGURE_URL）
请求消息内容：

```json
type C2SCEffectiveAdvertis struct {
	Token     string `json:"token"`
	RequestID string `json:"requestId"` // 唯一请求ID
}
```

#### 3.76 新春活动

##### 3.76.1  推送随机红包

请求消息命令字：Msgtype =471（MSG_SEND_RED_ENVELOPE_RSP）
主动推送：

```
{
	"code": 0,
	"message": "成功",
	"cdt":1, // 随机红包的CDT
	"CdtalCdt":1 // 用户总的CDT
}
```

##### 3.76.2 兑换cdt 广播消息

请求消息命令字：Msgtype =470（MSG_PUSH_RED_ENVELOPE_RESP）
主动推送

```
{
"code":0,
"message":"",
"cdt":0.2,  //用户兑换的CDT
"userName":"温2323"  // 用户昵称
}
```

##### 3.76.3 新春活动页面打点

请求消息命令字：Msgtype =472（MSG_DOUBLE_YEAR_DOT）
请求消息内容：

```go
const (
	DOT_DOUBLE_YEAR_FRAGMENT = 4000 // 碎片页打点
	DOT_DOUBLE_YEAR_DAY      = 4001 // 新春日统计
	DOT_DOUBLE_YEAR_TOTAL    = 4002 // 新春总统计
)
type Dot struct {
	Token    string `json:"token"`
	CodeType int    `json:"codeType"` // 打点类型
}

```

返回消息命令字：Msgtype =473（MSG_DOUBLE_YEAR_DOT_RSP）
返回消息内容：返回Code，Message等 。例子：

```json
{
	"code": 0,
	"message": "成功"
}
```

#### 3.77 推送公告跑马灯消息 （服务端推送）

返回消息命令字：Msgtype =467（MSG_PUSH_NOTICE_LANTERNS_RSP）
返回消息内容：返回Code，Message，data。例子

```json
{
    "code": 0, // 0 成功
    "message": "",
	"lantern": {   // 跑马灯消息
		"level": 9,
        "content": "lanterns3"   // 跑马灯类容
    }
}
```

#### 

#### 

#### 

### 附录一：游戏角色表

| 角色ID | 角色名称 | 性别 |
| ------ | :------: | ---: |
| 1      |   商人   |    0 |
| 2      |   商人   |    1 |
| 3      |   空乘   |    0 |
| 4      |   空乘   |    1 |
| 5      |   模特   |    0 |
| 6      |   模特   |    1 |
| 7      |   警察   |    0 |
| 8      |   警察   |    1 |
| 9      |   强盗   |    0 |
| 10     |   强盗   |    1 |
| 11     |   老人   |    0 |
| 12     |   老人   |    1 |
| 13     |  青少年  |    0 |
| 14     |  青少年  |    1 |
| 15     |   贵族   |    0 |
| 16     |   贵族   |    1 |
| 17     |   演员   |    0 |
| 18     |   演员   |    1 |
| 19     |   工人   |    0 |
| 20     |   工人   |    1 |
| 21     |  运动员  |    0 |
| 22     |  运动员  |    1 |
| 23     |   医生   |    0 |
| 24     |   医生   |    1 |
| 25     |   教授   |    0 |
| 26     |   教授   |    1 |
| 27     | 游客男1  |    0 |
| 28     | 游客女1  |    1 |

### 附录二：.消息类型表

```golang
	MSG_NULL_ACT       = 0
	MSG_REGISTER_PHONE = 1 //手机注册
	MSG_REGISTER_EMAIL = 2 //邮箱注册
	MSG_REGISTER_RSP   = 3 //注册返回
	MSG_LOGINANOTHER   = 4 //挤用户下线

	MSG_LOGIN                     = 5  //登录
	MSG_LOGIN_RSP                 = 6  //登录返回
	MSG_HEARTBEAT                 = 7  //心跳
	MSG_HEARTBEAT_RSP             = 8  //心跳返回
	MSG_REBIND                    = 9  //断线重连
	MSG_REBIND_RSP                = 10 //断线重连响应
	MSG_GET_USERINFO              = 11 //查看玩家信息
	MSG_GET_USERINFO_RSP          = 12 //查看玩家信息回复
	MSG_RESET_PASSWORD            = 13 //重置密码
	MSG_RESET_PASSWORD_RSP        = 14 //重置密码返回
	MSG_CREATER_ROLE              = 15 //创建角色
	MSG_CREATER_ROLE_RSP          = 16 //创建角色返回
	MSG_POSITION_CHANGE           = 17 //位置改变定时上报
	MSG_POSITION_CHANGE_RSP       = 18 //位置改变返回
	MSG_GET_POSITION              = 19 //获取当前位置
	MSG_GET_POSITION_RSP          = 20 //获取当前位置返回
	MSG_GET_VERIFICATION_CODE     = 21 //获取验证码
	MSG_GET_VERIFICATION_CODE_RSP = 22 //获取验证码返回
	MSG_CHECK_NICK_NAME           = 25 //检测昵称
	MSG_CHECK_NICK_NAME_RSP       = 26 //检测昵称返回

	MSG_GET_AMOUNT            = 27 //获取游戏子钱包金额
	MSG_GET_AMOUNT_RSP        = 28 //获取游戏子钱包金额响应
	MSG_GET_KNAPSACK          = 29 //获取用户背包信息
	MSG_GET_KNAPSACK_RSP      = 30 //获取用户背包信息返回
	MSG_ENTER_CITY            = 31 //进入城市
	MSG_ENTER_CITY_RSP        = 32 //进入城市返回
	MSG_BROAD_RAND_EVENT      = 33 //广播随机事件给前端
	MSG_BROAD_FINISH_EVENT    = 34 //广播完成事件给前端
	MSG_BROAD_POSITION        = 35 //广播玩家当前给前端
	MSG_BROAD_USER_OFFLINE    = 36 //广播玩家掉线
	MSG_FINISH_EVENT          = 37 //完成事件请求
	MSG_FINISH_EVENT_RSP      = 38 //完成事件返回
	MSG_UPDATE_ITEM_INFO      = 39 //用户背包和金额同步
	MSG_UPDATE_ITEM_INFO_RSP  = 40 //用户背包和金额同步响应
	MSG_GET_INVITE_USERS      = 41 //获取战友列表
	MSG_GET_INVITE_USERS_RSP  = 42 //获取战友列表响应
	MSG_GET_GRAB_COMRADES     = 43 //获取被抢走战友列表
	MSG_GET_GRAB_COMRADES_RES = 44 //获取被抢走战友列表响应
	MSG_GET_BUILDING_DESC     = 45 //获取建筑简介
	MSG_GET_BUILDING_DESC_RSP = 46 //获取建筑简介响应
	MSG_GET_MEMBER_SYS        = 47 //获取会员等级体系数据
	MSG_GET_MEMBER_SYS_RSP    = 48 //获取会员等级体系响应
	MSG_GET_USER_LEVEL        = 49 //获取用户当前等级状态
	MSG_GET_USER_LEVEL_RSP    = 50 //获取用户当前等级状态响应
	MSG_GRAB_COMRADE          = 51 //抢战友
	MSG_GRAB_COMRADE_RES      = 52 //抢战友响应
	MSG_GET_ITEMS_LIST        = 53 //获取道具商品列表
	MSG_GET_ITEMS_LIST_RSP    = 54 //获取道具商品列表响应
	MSG_BUY_ITEM              = 55 //购买道具
	MSG_BUY_ITEM_RSP          = 56 //购买道具响应
	MSG_MODIFY_NICKNAME       = 57 //修改昵称
	MSG_MODIFY_NICKNAME_RSP   = 58 //修改昵称响应
	MSG_GET_ALL_ROLE          = 59 //所有可用角色
	MSG_GET_ALL_ROLE_RSP      = 60
	MSG_USER_ADD_ROLE         = 61 //解锁角色
	MSG_USER_ADD_ROLE_RSP     = 62
	MSG_SELECT_ROLE           = 63 //角色选择
	MSG_SELECT_ROLE_RSP       = 64
	MSG_ROLE_CHANGE           = 65
	MSG_CERTIFICATION         = 66 //实名认证
	MSG_CERTIFICATION_RSP     = 67
	MSG_BIND_INVITER          = 69 //绑定邀请码
	MSG_BIND_INVITER_RSP      = 70
	MSG_DEPOSIT_REBATE        = 71 //好友邀请收益(好友存币返佣)
	MSG_DEPOSIT_REBATE_RSP    = 72
	MSG_QUERY_SHOP            = 73 //商家查询
	MSG_QUERY_SHOP_RSP        = 74
	MSG_GET_FRIEND_INFO       = 75 //获取好友详细信息
	MSG_GET_FRIEND_INFO_RSP   = 76
	MSG_GET_INVITATION        = 77 //获取用户的邀请关系
	MSG_GET_INVITATION_RSP    = 78
	MSG_GET_TOEKN_PAY_URL     = 79 //获取聚合支付链接
	MSG_GET_TOEKN_PAY_URL_RSP = 80
	MSG_USER_QUT_CITY         = 81 //退出城市
	MSG_BROAD_CITY_USER       = 82 //广播用户进入或离开城市事件
	MSG_GET_CITY_USER         = 83 //获取同城在线用户
	MSG_GET_CITY_USER_RSP     = 84

	MSG_GET_EMAIL_LIST      = 85 //获取邮件列表
	MSG_GET_EMAIL_LIST_RSP  = 86
	MSG_DEL_EMAIL           = 87 // 删除邮件
	MSG_DEL_EMAIL_RSP       = 88
	MSG_SET_EMAIL_READ      = 89 // 设置邮件为已读
	MSG_SET_EMAIL_READ_RSP  = 90
	MSG_COUNT_EMAIL         = 91 // 获取邮件数量
	MSG_COUNT_EMAIL_RSP     = 92
	MSG_PUSH_USER_MEIAL     = 93 // 发送实名邮件 (协议已删除)
	MSG_PUSH_USER_MEIAL_RSP = 94
	MSG_PUSH_NOTICE         = 95 // 请求发送推送公告 （已删除）
	MSG_PUSH_NOTICE_RSP     = 96
	MSG_GET_UPGRADE_NOTICE  = 97  //获取最新升级公告
	MSG_SEND_KYC_STATUS     = 98  //同步kyc状态（后端推送给前端）
	MSG_UPDATE_TASK_INFO    = 99  //更新任务信息（后端推送给前端）
	MSG_UPDATE_TASK_STATUS  = 100 //开启或关闭任务（后台推送给后端）

	// 在游戏中展示本地生活类型配置,　编号为101～200.
	MSG_LOCAL_LIFE_INDEX_RECOMMEND           = 101 // 首页推荐
	MSG_LOCAL_LIFE_INDEX_RECOMMEND_RESP      = 102
	MSG_LOCAL_LIFE_TOP_SEARCH                = 103 // 历史和热门搜索
	MSG_LOCAL_LIFE_TOP_SEARCH_RESP           = 104
	MSG_LOCAL_LIFE_HOTEL_SEARCH              = 105 // 酒店搜索
	MSG_LOCAL_LIFE_HOTEL_SEARCH_RESP         = 106
	MSG_LOCAL_LIFE_HOTEL_DETAIL              = 107 // 获取酒店详情
	MSG_LOCAL_LIFE_HOTEL_DETAIL_RESP         = 108
	MSG_LOCAL_LIFE_ROOM_DETAIL               = 109 // 获取房间详情
	MSG_LOCAL_LIFE_ROOM_DETAIL_RESP          = 110
	MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD      = 111 // 删除历史搜索记录
	MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD_RESP = 112
	MSG_LOCAL_LIFE_CITY_LIST                 = 113 // 获取城市列表
	MSG_LOCAL_LIFE_CITY_LIST_RESP            = 114
	MSG_LOCAL_LIFE_STORE_ClASSIFY            = 115 // 获取店铺分类
	MSG_LOCAL_LIFE_STORE_ClASSIFY_RESP       = 116
	MSG_LOCAL_LIFE_INDEX_SELECT_CITY         = 117 // 首页选择城市
	MSG_LOCAL_LIFE_INDEX_SELECT_CITY_RESP    = 118

	MSG_GET_SIGNIN_LIST     = 201 //获取签到列表
	MSG_GET_SIGNIN_LIST_RSP = 202
	MSG_SIGN_IN             = 203 //签到
	MSG_SIGN_IN_RSP         = 204
	MSG_GET_TASK_LIST       = 205 //获取任务列表
	MSG_GET_TASK_LIST_RSP   = 206
	MSG_GET_TASK_AWARD      = 207 //任务奖励领取
	MSG_GET_TASK_AWARD_RSP  = 208
	MSG_ENTER_SHOP          = 209 //商家跳转请求
	MSG_BROAD_USER_TASK     = 210 //任务完成通知

	MSG_EMAIL_RECEIVE_REWARDS     = 211 //
	MSG_EMAIL_RECEIVE_REWARDS_RSP = 212 //

	MSG_GET_MERCHANTS_URL     = 213 // 获取商户连接
	MSG_GET_MERCHANTS_URL_RSP = 214 //

	// 统计数据用的的 301 - 400
	MSG_STATISTICS_CITY_ICON = 301 //城市中心 icon 入口统计

	// 双旦活动 401~450.
	MSG_SWEET_TREE                  = 401 // 圣诞树＋圣诞糖果对换cdt.
	MSG_SWEET_TREE_RSP              = 402
	MSG_RANK_LIST_UPDATE_PROP       = 403 // 更新道具分数值.
	MSG_RANK_LIST_UPDATE_PROP_RSP   = 404
	MSG_RANK_LIST_DAY               = 405 // 每日排行榜.
	MSG_RANK_LIST_DAY_RSP           = 406
	MSG_RANK_LIST_ALL               = 407 // 总排行榜.
	MSG_RANK_LIST_ALL_RSP           = 408
	MSG_RANK_LIST_DAY_PROP          = 409 // 每日双旦值记录.
	MSG_RANK_LIST_DAY_PROP_RSP      = 410
	MSG_RANK_LIST_ALL_PROP          = 411 // 总日双旦值记录.
	MSG_RANK_LIST_ALL_PROP_RSP      = 412
	MSG_RANK_LIST_INVITE_RECORD     = 413 // 邀请记录.
	MSG_RANK_LIST_INVITE_RECORD_RSP = 414
	MSG_GET_DOUBLE_YEAR_STATUS      = 415 // 或取双旦活动状态
	MSG_GET_DOUBLE_YEAR_STATUS_RSP  = 416
	MSG_BROAD_SANTA_CARD            = 417 //广播圣诞老人卡奖励信息给前端
	MSG_BROAD_ACTIVITY_END          = 418 //广播活动结束给前端
	MSG_GET_SWEET_AND_TREE          = 420 //获取用户双蛋圣诞树，糖果
	MSG_GET_SWEET_AND_TREE_RSP      = 421
	MSG_GET_PATCH                   = 422 //获取用户双蛋碎片
	MSG_GET_PATCH_RSP               = 423
	MSG_RANK_LIST_DAY_CDT           = 425 // 每日排行榜奖励发放cdt.
	MSG_RANK_LIST_DAY_CDT_RSP       = 426
	MSG_RANK_LIST_ALL_CDT           = 427 // 总排行榜奖励发放cdt.
	MSG_RANK_LIST_ALL_CDT_RSP       = 428
	MSG_CHRISMAS_ROLE               = 429 // 圣诞老人角色对换cdt.
	MSG_SPECAIL_CDT                 = 449 // 测试特殊情况获取cdt.
	MSG_SPECAIL_CDT_RSP             = 450
	MSG_SANTA_LAUS_ROLE_EXPIRE_RSP  = 451 // 圣诞老人角色过期通知
	MSG_RANK_LIST_EMAIL_REPAIR_CDT      = 453 // 补发cdt
	MSG_RANK_LIST_EMAIL_REPAIR_CDT_RESP = 454

	//宝箱
	MSG_OPEN_STREASURE_BOX           = 455 //打开宝箱
	MSG_OPEN_STREASURE_BOX_RSP       = 456
	MSG_FINISH_STREASURE_BOX         = 457 //完成宝箱
	MSG_FINISH_STREASURE_BOX_RSP     = 458
	MSG_RECEIVE_STREASURE_BOX        = 459 //领取奖励
	MSG_RECEIVE_STREASURE_BOX_RSP    = 460
	MSG_GET_STREASURE_BOX_RECORD     = 461 //获取宝箱记录
	MSG_GET_STREASURE_BOX_RECORD_RSP = 462
	MSG_ACTIVITY_STATYS              = 463 //获取活动状态
	MSG_ACTIVITY_STATYS_RSP          = 464
	MSG_PUSH_ACTIVITY_STATUS_RSP     = 465 //主动推送活动状态
	MSG_BROAD_CLEAN_EVENT            = 466 //广播清除事件给前端

	MSG_PUSH_RED_ENVELOPE_RESP          = 470 // 兑换红包推送跑马灯
	MSG_SEND_RED_ENVELOPE_RSP           = 471 // 主动推送消息 随机红包，红包

```

### 附录三：道具类别表

| ID   |  角色名称  |
| ---- | :--------: |
| 10   |    cdt    |
| 20   | 角色解锁卡 |
| 30   |  普通道具  |
| 40   |    碎片    |

### 附录四：道具编号表

| 道具编号 | ·类别编号 | 道具名称 		| 是否绑定 | 道具品质 | 是否堆叠 | 道具图片 | 说明文字 | 获取途径 | 使用跳转 | 价格 |
| -------- | -------- | --------------| -------- | -------- | -------- | -------- | -------- | -------- | -------- | ------------ |
| 1001     | 10       | cdt    	  | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2001     | 20       | 男警解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2002     | 20       | 女警解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2003     | 20       | 男工解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2004     | 20       | 女工解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2008     | 20       | 商人解锁卡    | 1        | 1        | 10       |          |          |          |          | 1000         |
| 2009     | 20       | 商人解锁卡    | 1        | 1        | 10       |          |          |          |          | 1000         |
| 2010     | 20       | 教授解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2011     | 20       | 教授解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2012     | 20       | 演员解锁卡    | 1        | 1        | 10       |          |          |          |          | 30000        |
| 2013     | 20       | 演员解锁卡    | 1        | 1        | 10       |          |          |          |          | 30000        |
| 2014     | 20       | 空乘解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2015     | 20       | 空乘解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2016     | 20       | 模特解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 2017     | 20       | 模特解锁卡    | 0        | 1        | 10       |          |          |          |          | 1000         |
| 3001     | 30       | 改名卡   	  | 1        | 2        | 10       |          |          |          |          | 1000         |
| 3002     | 30       | 补签卡 		 | 1 		 | 2 		| 10 	   | 		  |          |          |          | 1000         |
| 4001     | 40       | 一般碎片      | 0        | 1        | 10       |          |          |          |          | 1000         |
| 4002     | 40       | 中级碎片      | 0        | 1        | 10       |          |          |          |          | 1000         |
| 4003     | 40       | 高级碎片      | 0        | 1        | 10       |          |          |          |          | 1000         |
