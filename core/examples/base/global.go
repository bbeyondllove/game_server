package base

const ReadBufferSize = 10240 //单位byte
const PacketMaxSize = 1024 * 64

const NUM_MSG_TYPES = 1000 //协议总数

const (
	STATUS_NEVER     = 0 // Opcode not accepted from client (deprecated or server side only)
	STATUS_AUTHED    = 1 // Player authenticated
	STATUS_UNHANDLED = 2 // We don' handle this opcode yet
)

//消息类型
const (
	MSG_NULL_ACT       = 0
	MSG_REGISTER_EMAIL = 1 //邮箱注册
	MSG_REGISTER_PHONE = 2 //手机注册
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
	MSG_CHECK_ACCOUNT             = 23 //检测帐号
	MSG_CHECK_ACCOUNT_RSP         = 24 //检测帐号返回
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
	MSG_UPDATE_POINT          = 39 //积分变更
	MSG_UPDATE_POINT_RSP      = 40 //积分变更响应
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

)

var ActReqCount uint64       //操作数
var ActRspSucessCount uint64 //操作成功数
var ActRspFailCount uint64   //操作失败数
var ActReqTimePoint int64    //开始时间
var ActLastTimePoint int64   //最后响应时间

var ActionRspSucessCount uint64       //操作响应成功数
var PreActionRspSucessCount uint64    //上一次操作成功响应数
var PreActionRspSucessTimePoint int64 //上一次响应成功计算时间点

var ActionRepCount uint64       //操作请求数
var PreActionRepCount uint64    //上一次操作请求数
var PreActionRepTimePoint int64 //上一次请求计算时间点

var ActionRspFailCount uint64       //操作失败数
var PreActionRspFailCount uint64    //上一次操作失败数
var PreActionRspFailTimePoint int64 //上一次操作失败响应时间点
