package proto

import (
	"game_server/game/model"
	"sync"

	"github.com/shopspring/decimal"
)

const (
	JSON = iota
	TEXT
)

const (
	USE_REGISTR      = 1 //验证码用来注册
	USE_RESET_PASSWD = 2 //验证码用来重置密码
	USE_BIND         = 3 //验证码用来绑定(手机或邮箱)
	CODE_TYPE_PHONE  = 1 //验证码发给手机
	CODE_TYPE_EMAIL  = 2 //验证码发给邮箱
	STATUS_LOGINED   = 2
	MIN_BUF_SIZE     = 20   //最小合法请求消息长度
	MAX_BUF_SIZE     = 2048 //最大合法请求消息长度
	RAND_STRING_LEN  = 16   //随机字符串长度
	SysType          = 2
	ROBOT_USER       = 9 //机器人类型
	ONLINE_KEY       = "task_daily_online_time:"
)
const (
	TASK_NORMAL = 1 //主线日任务
	TASK_DAILY  = 2 //每日任务
	TASK_SIGNIN = 3 //签到任务

)

var (
	STATIC_KEY = map[int]string{
		TASK_NORMAL: "task_normal:",
		TASK_DAILY:  "task_daily:",
		TASK_SIGNIN: "task_signin:",
	}
)

const (
	TASK_STOP = iota
	TASK_START
)

const (
	ACTIVITY_TYPE_NOMAL        = iota //普通事件
	ACTIVITY_TYPE_DOUBLE_YEAR         //双蛋事件
	ACTIVITY_TYPE_TREASURE_BOX        // 欢乐宝箱
)

const (
	ACTIVITY_NOT_START = iota //双蛋活动未开始
	ACTIVITY_START            //双蛋活动进行中
	ACTIVITY_END              //双蛋活动已结束
)

const (
	ACTIVITY_SIGN_IN         = 10   // 每日签到
	ACTIVITY_TREASURE_BOX    = 7000 // 广告
	ACTIVITY_SPRING_FESTIVAL = 8000 // 春节活动
)

const (
	ITEM_PATCH = 4101 //碎片
	ITEM_SWEET = 6106 //祝福天使
	ITEM_TREE  = 6107 //幸运箭矢
)

const (
	ENTER_CITY = 1 //进入城市
	QUIT_CITY  = 2 //退出城市
)

const (
	LIMITE_NUM   = 1 //限制数量
	LIMITE_COUNT = 2 //限制次数
)

const (
	ITEM_CDT  = 10 //CDT
	ITEM_LOCK = 20 //解锁卡

)

const (
	SIGNIN_NORMAL = 1 //正常签到
	SIGNIN_LOST   = 2 //补签
)

const (
	PLATFORM_WEB     = 1 //web平台
	PLATFORM_ANDROID = 2 //安卓平台
	PLATFORM_IOS     = 3 //IOS平台
)

const (
	PLATFORM_WEB_STR     = "Windows" // web平台,或者其他字符串
	PLATFORM_ANDROID_STR = "Android" // 安卓平台
	PLATFORM_IOS_STR     = "iOS"     //IOS平台
)

const (
	NOTICE_SYSTEM_UPDATE = 0 // 系统更新公告
	NOTICE_NORNAL        = 1 // 普通公告
	NOFICE_FORCE         = 2 // 强制更新公告
	NOTICE_BEFOR_UPDATE  = 3 // 更新前公告
)

const (
	DOT_DOUBLE_YEAR_FRAGMENT = 4000 // 碎片页打点
	DOT_DOUBLE_YEAR_DAY      = 4001 // 新春日统计
	DOT_DOUBLE_YEAR_TOTAL    = 4002 // 新春总统计
)

//消息类型
const (
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
	MSG_GET_TOEKN_PAY_URL     = 79 //获取聚合支付链接 (已废弃)
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
	MSG_LOCAL_LIFE_INDEX_RECOMMEND              = 101 // 首页推荐
	MSG_LOCAL_LIFE_INDEX_RECOMMEND_RESP         = 102
	MSG_LOCAL_LIFE_TOP_SEARCH                   = 103 // 历史和热门搜索
	MSG_LOCAL_LIFE_TOP_SEARCH_RESP              = 104
	MSG_LOCAL_LIFE_HOTEL_SEARCH                 = 105 // 酒店搜索
	MSG_LOCAL_LIFE_HOTEL_SEARCH_RESP            = 106
	MSG_LOCAL_LIFE_HOTEL_DETAIL                 = 107 // 获取酒店详情
	MSG_LOCAL_LIFE_HOTEL_DETAIL_RESP            = 108
	MSG_LOCAL_LIFE_ROOM_DETAIL                  = 109 // 获取房间详情
	MSG_LOCAL_LIFE_ROOM_DETAIL_RESP             = 110
	MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD         = 111 // 删除历史搜索记录
	MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD_RESP    = 112
	MSG_LOCAL_LIFE_CITY_LIST                    = 113 // 获取城市列表
	MSG_LOCAL_LIFE_CITY_LIST_RESP               = 114
	MSG_LOCAL_LIFE_STORE_ClASSIFY               = 115 // 获取店铺分类
	MSG_LOCAL_LIFE_STORE_ClASSIFY_RESP          = 116
	MSG_LOCAL_LIFE_INDEX_SELECT_CITY            = 117 // 首页选择城市
	MSG_LOCAL_LIFE_INDEX_SELECT_CITY_RESP       = 118
	MSG_LOCAL_LIFE_CATEGORY_SEARCH_STORE        = 119 // 分类搜索-店铺
	MSG_LOCAL_LIFE_CATEGORY_SEARCH_STORE_RESP   = 120
	MSG_LOCAL_LIFE_SEARCH_SUGGEST               = 121 // 搜索联想
	MSG_LOCAL_LIFE_SEARCH_SUGGEST_RESP          = 122
	MSG_LOCAL_LIFE_TOP_SEARCH_V2                = 123 // V2版本-历史和热门搜索
	MSG_LOCAL_LIFE_TOP_SEARCH_V2_RESP           = 124
	MSG_LOCAL_LIFE_STORE_TYPE                   = 125 // V2版本-查询店铺分类信息
	MSG_LOCAL_LIFE_STORE_TYPE_RESP              = 126
	MSG_LOCAL_LIFE_STORE_HOTEL_DETAIL           = 127 // 住宿详情
	MSG_LOCAL_LIFE_STORE_HOTEL_DETAIL_RESP      = 128
	MSG_LOCAL_LIFE_STORE_RESTAURANT_DETAIL      = 129 // 美食详情
	MSG_LOCAL_LIFE_STORE_RESTAURANT_DETAIL_RESP = 130
	MSG_LOCAL_LIFE_GOODS_HOTEL_DETAIL           = 131 // 商品住宿(房间详情)/商品团购(套餐详情)/商品代金劵(优惠劵详情)
	MSG_LOCAL_LIFE_GOODS_HOTEL_DETAIL_RESP      = 132
	MSG_LOCAL_LIFE_DISCOUNT_DETAIL              = 133 // 美食优惠劵详情
	MSG_LOCAL_LIFE_DISCOUNT_DETAIL_RESP         = 134
	MSG_LOCAL_LIFE_GOODS_RESTAURANT_DETAIL      = 135 // 美食套餐
	MSG_LOCAL_LIFE_GOODS_RESTAURANT_DETAIL_RESP = 136
	MSG_LOCAL_LIFE_CITY                         = 137 // 查询城市列表
	MSG_LOCAL_LIFE_CITY_RESP                    = 138
	MSG_LOCAL_LIFE_CITY_PICK                    = 139 // 切换城市
	MSG_LOCAL_LIFE_CITY_PICK_RESP               = 140
	MSG_LOCAL_LIFE_CITY_SUGGEST                 = 141 // 城市联想
	MSG_LOCAL_LIFE_CITY_SUGGEST_RESP            = 142
	MSG_LOCAL_LIFE_SEARCH_RECORD_DELETE         = 143 // v2版本－删除历史搜索记录
	MSG_LOCAL_LIFE_SEARCH_RECORD_DELETE_RESP    = 144
	MSG_LOCAL_LIFE_CATEGORY_SEARCH_GOODS        = 145 // 分类搜索-商品
	MSG_LOCAL_LIFE_CATEGORY_SEARCH_GOODS_RESP   = 146

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

	MSG_GET_MERCHANTS_URL     = 213 // 获取商户连接 (已废弃)
	MSG_GET_MERCHANTS_URL_RSP = 214 //
	MSG_GET_CONFIGURE_URL     = 215 // 获取商户连接 (已废弃)
	MSG_GET_CONFIGURE_URL_RSP = 216 //

	// 统计数据用的的 301 - 400
	MSG_STATISTICS_CITY_ICON  = 301 //城市中心 icon 入口统计
	MSG_TREASURE_BOX_ADVERTIS = 302 // 统计宝箱有效广告

	MSG_TREASURE_BOX_DAY_NUM     = 350 //获取当前宝箱打开的次数
	MSG_TREASURE_BOX_DAY_NUM_RSP = 351 //

	// 双旦活动 401~500.
	MSG_SWEET_TREE                      = 401 // 圣诞树＋圣诞糖果对换cdt. (新春活动 祝福天使+幸运箭矢)
	MSG_SWEET_TREE_RSP                  = 402
	MSG_RANK_LIST_UPDATE_PROP           = 403 // 更新道具分数值.
	MSG_RANK_LIST_UPDATE_PROP_RSP       = 404
	MSG_RANK_LIST_DAY                   = 405 // 每日排行榜.
	MSG_RANK_LIST_DAY_RSP               = 406
	MSG_RANK_LIST_ALL                   = 407 // 总排行榜.
	MSG_RANK_LIST_ALL_RSP               = 408
	MSG_RANK_LIST_DAY_PROP              = 409 // 每日双旦值记录.
	MSG_RANK_LIST_DAY_PROP_RSP          = 410
	MSG_RANK_LIST_ALL_PROP              = 411 // 总日双旦值记录.
	MSG_RANK_LIST_ALL_PROP_RSP          = 412
	MSG_RANK_LIST_INVITE_RECORD         = 413 // 邀请记录.
	MSG_RANK_LIST_INVITE_RECORD_RSP     = 414
	MSG_GET_DOUBLE_YEAR_STATUS          = 415 // 或取双旦活动状态
	MSG_GET_DOUBLE_YEAR_STATUS_RSP      = 416
	MSG_BROAD_SANTA_CARD                = 417 //广播圣诞老人卡奖励信息给前端
	MSG_BROAD_ACTIVITY_END              = 418 //广播活动结束给前端
	MSG_GET_SWEET_AND_TREE              = 420 //获取用户双蛋圣诞树，糖果
	MSG_GET_SWEET_AND_TREE_RSP          = 421
	MSG_GET_PATCH                       = 422 //获取用户双蛋碎片
	MSG_GET_PATCH_RSP                   = 423
	MSG_RANK_LIST_DAY_CDT               = 425 // 每日排行榜奖励发放cdt.
	MSG_RANK_LIST_DAY_CDT_RSP           = 426
	MSG_RANK_LIST_ALL_CDT               = 427 // 总排行榜奖励发放cdt.
	MSG_RANK_LIST_ALL_CDT_RSP           = 428
	MSG_CHRISMAS_ROLE                   = 429 // 圣诞老人角色对换cdt.
	MSG_SPECAIL_CDT                     = 449 // 测试特殊情况获取cdt.
	MSG_SPECAIL_CDT_RSP                 = 450
	MSG_SANTA_LAUS_ROLE_EXPIRE_RSP      = 451 // 圣诞老人角色过期通知
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
	MSG_PUSH_NOTICE_LANTERNS_RSP     = 467 //广播公告走马灯消息

	MSG_PUSH_RED_ENVELOPE_RESP = 470 // 兑换红包推送跑马灯
	MSG_SEND_RED_ENVELOPE_RSP  = 471 // 主动推送消息 随机红包，红包

	MSG_DOUBLE_YEAR_DOT     = 472 // 双旦UI打点
	MSG_DOUBLE_YEAR_DOT_RSP = 473 // 双旦UI打点回复
)

//通用返回消息
type S2CCommon struct {
	Code    int32  `json:"code"`    //错误代码
	Message string `json:"message"` //错误信息
}

//心跳消息返回
type HeartBeatRsp struct {
	S2CCommon
	Token string `json:"token"`
}

//心跳消息
type HeartBeat struct {
	UserId string `json:"userId"`
}

//用户系统返回消息
type S2C_HTTP struct {
	Code    int32                  `json:"code"`    //错误代码
	Message string                 `json:"message"` //错误信息
	Data    map[string]interface{} `json:"data"`
}

type S3C_HTTP struct {
	Code    int32       `json:"code"`    //错误代码
	Message string      `json:"message"` //错误信息
	Data    interface{} `json:"data"`
}

type S2CHTTP struct {
	Code    int32       `json:"code"` //错误代码
	Message string      `json:"msg"`  //错误信息
	Data    interface{} `json:"data"`
}

type S3C_PAGE struct {
	S2CHTTP
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type StatisticsDoubleYearCdtRsp struct {
	S3C_PAGE
	Pv  int     `json:"pv"`
	Uv  int     `json:"uv"`
	Cdt float32 `json:"cdt"`
}

type PayLoad struct {
	Exp     int64  `json:"exp"`     //过期时间
	Iat     int64  `json:"iat"`     //签发时间
	Nbf     int64  `json:"nbf"`     //可用开始时间
	UserId  string `json:"userId"`  //用户ID
	SysType int64  `json:"sysType"` //业务系统
}

type EventRate struct {
	ActivityType int     `json:"activity_type"`
	Id           int     `json:"id"`
	ItemId       int     `json:"item_id"`
	ItemNum      float32 `json:"num"`
	LimitType    int     `json:"limit_type"`
	LimitNum     int     `json:"limit_num"`
	UnitNum      int     `json:"unit_num"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	StartDate    string  `json:"start_data"`
	EndDate      string  `json:"end_date"`
}

type EventNode struct {
	ActivityType int    `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件
	Type         string `json:"type"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type S2CActivityEnd struct {
	ActivityType int `json:"activityType"` //活动ID 0:普通随机事件；1：双蛋随机事件

}

// 宝箱活动CDT随机表
type TreasureBoxCdtRand struct {
	Probability  int64           `json:"probability"`
	RewardItem   string          `json:"reward_item"`
	RewardNumber decimal.Decimal `json:"reward_number"`
}

// 宝箱活动表
type TreasureBox struct {
	StateSwitch                 int                  `json:"status_switch"`                   //状态变更开关,0:关闭,1:开启
	EventInterval               int                  `json:"event_interval"`                  //事件间隔
	EventNum                    int                  `json:"event_num"`                       //事件数量
	EventTimeoutSwitch          int                  `json:"event_timeout_switch"`            //事件超时机制开关
	EventTimeout                int                  `json:"event_timeout"`                   //事件超时
	EventTimeoutCheckInterval   int                  `json:"event_timeout_check_interval"`    //事件超时检查间隔
	LockTimeout                 int                  `json:"lock_timeout"`                    //事件锁定超时时长
	ExceptionInterval           int                  `json:"exception_interval"`              //事件处理异常检查间隔
	ExceptionCount              int                  `json:"exception_count"`                 //事件处理异常数量临界值,超过则不允许用户处理宝箱
	ExceptionDayCount           int                  `json:"exception_day_count"`             //事件处理异常数量临界值,每天打开的次数
	CdtRandTableConfigFrom      int                  `json:"cdt_rand_table_config_from"`      // cdt配置加载位置,0:配置文件,1:DB
	CdtRandTableReloadTimeout   int64                `json:"cdt_rand_table_reload_timeout"`   // cdt重加载超时时间
	CdtMin                      float32              `json:"cdt_min"`                         // cdt随机最小值
	CdtMax                      float32              `json:"cdt_max"`                         // cdt随机最大值
	CdtSecondCritical           float32              `json:"cdt_second_critical"`             // cdt二次随机临界值
	CdtSecondCriticalRetryCount int                  `json:"cdt_second_critical_retry_count"` // cdt随机临界值重试次数
	CdtRandTable                []TreasureBoxCdtRand `json:"cdt_rand_table"`
}

type Backstage struct {
	TokenSwitch         int    `json:"token_switch"`          // 登陆token开关
	TokenExpirationTime int    `json:"token_expiration_time"` // 登陆token过期时间
	TokenSecret         string `json:"token_secret"`          // 登陆token密钥
	MaxActivityYear     int    `json:"max_activity_year"`     // 最大的活动年份
	FuwaFragmentId      int    `json:"fuwa_fragment_id"`      // 福娃活动碎片ID
	NoticeLanternLevel  int    `json:"notice_lantern_level"`  // 公告跑马灯等级,用于前端分离不同等级跑马灯消息
	MaxUploadImageSize  int64  `json:"max_upload_image_size"` // 最大的上传图片大小
	UserStatusSwitch    int    `json:"user_status_switch"`    // 用户禁用状态激活开关
}

type BaseConf struct {
	LocationId           []int       `json:"location_id"`
	EventInterval        int         `json:"event_interval"`
	HouseMaxNum          int         `json:"house_max_num"`
	PingminARoles        string      `json:"pingmin_aRoles"`
	PingminDRoles        string      `json:"pingmin_dRoles"`
	EventNum             int         `json:"event_num"`
	EventRate            []EventRate `json:"event_rate"`
	EventOffset          int         `json:"event_offset"`
	EventMap             []Position  `json:"event_map"`
	TokenCode            []string    `json:"token_code"`
	BirthPlace           Position    `json:"birth_place"`
	BirthScreenPlaceInit int         `json:"birth_screen_place_init"`
	BirthScreenPlace     Position    `json:"birth_screen_place"`

	//双蛋活动结构体
	DoubleYear      DoubleYear          `json:"double_year"`
	DoubleYearArea  map[string]Position `json:"double_year_area"`
	DoubleEventRate []EventRate         `json:"double_event_rate"`

	//宝箱结构体
	TreasureBox          TreasureBox `json:"treasure_box"`
	TreasureBoxEventRate []EventRate `json:"treasure_box_event_rate"`

	//后台功能配置
	Backstage Backstage `json:"backstage"`
}

type DoubleYear struct {
	EventNum int `json:"event_num"`
}

type EventData struct {
	EventNode
	EventId int `json:"EventId"` //事件id
}

// 可过期的事件
type ExpiredEventData struct {
	EventData
	Timestamp int64 `json:"timestamp"` //事件开始时间
}

type CityMap struct {
	DataMap map[int]map[int]EventNode
	Mutex   sync.RWMutex
}

type EventMap struct {
	DataMap map[int]map[int]*EventData
	Mutex   sync.RWMutex
}

//用户统计
type C2SGetUserSum struct {
	StartTime string `json:"startTime"`
	EndtTime  string `json:"endtTime"`
}

//用户统计返回
type S2CGetUserSum struct {
	model.HttpCommon
	AllCount int `json:"allCount"`
	NewCount int `json:"newCount"`
}

//前端请求
type C2SBase struct {
	Token string `json:"token"`
}

//红包
type ReaEnvelRand struct {
	Probability  int64   `json:"probability"`
	RewardItem   string  `json:"reward_item"`
	RewardNumber float32 `json:"reward_number"`
}
