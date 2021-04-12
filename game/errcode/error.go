package errcode

var (
	MSG_SUCCESS                        = int32(0)
	HTTP_SUCCESS                       = int32(200)
	ERROR_REQUEST_NOT_ALLOW            = int32(5001)
	ERROR_SYSTEM                       = int32(5002)
	ERROR_LOGIN_REPEAT                 = int32(5003)
	ERROR_ROLE_EXITS                   = int32(5005)
	ERROR_WALLET                       = int32(5006)
	ERROR_REDIS                        = int32(5007)
	ERROR_MYSQL                        = int32(5008)
	ERROR_PARAM_ILEGAL                 = int32(5009)
	ERROR_OBJ_NOT_EXISTS               = int32(5010)
	ERROR_ACCOUNT_EXIST                = int32(5011)
	ERROR_NOT_MEMBER                   = int32(5012)
	ERROR_KNAPSACK                     = int32(5013)
	ERROR_NO_FREE_MODIFY_COUNT         = int32(5014)
	ERROR_NO_MODIFY_NAME_CARD          = int32(5015)
	ERROR_NOT_ENOUGH_MONEY             = int32(5016)
	ERROR_USER_SYSTEM                  = int32(5017)
	ERROR_NO_CLEAR_CARD                = int32(5018)
	ERROR_CDT_OUT_OF_TODAY_LIMIT       = int32(5019)
	ERROR_PARAM_EMPTY                  = int32(5020)
	ERROR_ROLE_IS_UNLOCK               = int32(5021)
	ERROR_CURDAY_IS_FINISHED           = int32(5022)
	ERROR_NO_SIGNIN_LOST_CARD          = int32(5023)
	ERROR_ITEM_OUT_OF_TODAY_LIMIT      = int32(5024)
	ERROR_NONEED_SIGNIN                = int32(5025)
	ERROR_NONEED_ADDSIGNIN             = int32(5026)
	ERROR_MUST_SINGIN_CURRDAY          = int32(5027)
	ERROR_EMAIL_NOT_FOUND              = int32(5028)
	ERROR_EMAIL_FAILED_TO_CLAIM_REWARD = int32(5029)
	ERROR_EMAIL_REWARD_RECEIVED        = int32(5030)
	ERROR_ITMES_EXPIRED                = int32(5031)
	ERROR_ROLES_EXPIRED                = int32(5032)
	ERROR_PRIZE_EXPIRED                = int32(5033)
	ERROR_BOX_IS_LOCK                  = int32(5034)
	ERROR_BOX_OPEN_NUM                 = int32(5035)
	ERROR_BOX_AWARD_EXCEPTION          = int32(5036)
	ERROR_BOX_AWARD_NOT_FOUND          = int32(5037)
	ERROR_BOX_NOT_FINISHED             = int32(5038)
	ERROR_BOX_DAY_OPEN_NUM             = int32(5039)

	ERROR_USER_NOT_EXIT       = int32(302)
	ERROR_UPDATE_MONEY        = int32(303)
	ERROR_AMOUNTAVAILABLE     = int32(304)
	ERROR_ADD_EXCHANGE_RECORD = int32(305)
	ERROR_MYSQL_COMMIT        = int32(306)
	ERROR_ILLEGAL             = int32(307)
	ERROR_NOT_LOGIN           = int32(308)

	ERROR_HTTP_AMOUNT_ILLEGAL     = int32(40000)
	ERROR_HTTP_SIGNATURE          = int32(40001)
	ERROR_HTTP_USER_NOT_EXIST     = int32(40101)
	ERROR_HTTP_FORBIDDEN          = int32(40102)
	ERROR_HTTP_BALANCE_NOT_ENOUGH = int32(40103)
	ERROR_HTTP_REPEAT_TXID        = int32(40105)
	ERROR_HTTP_USER_NOT_ALLOW     = int32(40106)
	ERROR_HTTP_ACCOUNT_ERROR      = int32(40107)

	ERROR_PUSH_NOTICE = int32(5000)
	ERROR_PUSH_EMAIL  = int32(5000)

	// cdt犯错误码为5101~5200.
	ERROR_CDT_DAY_FULL        = int32(5101)
	ERROR_CDT_LACK_OF_BALANCE = int32(5102)
	ERROR_NOT_ENOUGH_ITEM     = int32(5103)

	ERROR_NOT_START = int32(6000)
	ERROR_END       = int32(6001)

	ERROR_MSG = map[int32]string{
		MSG_SUCCESS:             "成功",
		HTTP_SUCCESS:            "成功",
		ERROR_REQUEST_NOT_ALLOW: "非法请求",
		ERROR_SYSTEM:            "系统错误",
		ERROR_LOGIN_REPEAT:      "重复登录",
		ERROR_ROLE_EXITS:        "角色已存在",
		ERROR_ACCOUNT_EXIST:     "账号已经存在",
		ERROR_WALLET:            "操作钱包错误",
		ERROR_REDIS:             "系统错误",
		ERROR_MYSQL:             "系统错误",
		ERROR_PARAM_ILEGAL:      "参数错误",
		ERROR_KNAPSACK:          "操作背包错误",

		ERROR_HTTP_AMOUNT_ILLEGAL:     "金额非法",
		ERROR_HTTP_SIGNATURE:          "签名校验失败",
		ERROR_HTTP_USER_NOT_EXIST:     "用户不存在",
		ERROR_HTTP_FORBIDDEN:          "禁止划转",
		ERROR_HTTP_BALANCE_NOT_ENOUGH: "余额不足",
		ERROR_HTTP_REPEAT_TXID:        "重复的订单ID",
		ERROR_HTTP_USER_NOT_ALLOW:     "账户被禁用",
		ERROR_HTTP_ACCOUNT_ERROR:      "账户密码错误",
		ERROR_USER_NOT_EXIT:           "用户不存在",
		ERROR_UPDATE_MONEY:            "更新金额失败",
		ERROR_AMOUNTAVAILABLE:         "可用金额不足",
		ERROR_ADD_EXCHANGE_RECORD:     "插入交易记录失败",
		ERROR_MYSQL_COMMIT:            "系统错误",
		ERROR_ILLEGAL:                 "请求参数非法",
		ERROR_NOT_LOGIN:               "未登录",
		ERROR_OBJ_NOT_EXISTS:          "目标不存在",
		ERROR_NOT_ENOUGH_MONEY:        "金额不够",
		ERROR_NOT_MEMBER:              "非会员",
		ERROR_NO_FREE_MODIFY_COUNT:    "免费修改次数已用完",
		ERROR_NO_MODIFY_NAME_CARD:     "没有改名卡",
		ERROR_USER_SYSTEM:             "用户系统接口调用错误",
		ERROR_NO_CLEAR_CARD:           "没有解锁卡",
		ERROR_CDT_OUT_OF_TODAY_LIMIT:  "用户每天获取的CDT已经达到上限",
		ERROR_ITEM_OUT_OF_TODAY_LIMIT: "用户每天获取的物品已经达到上限",
		ERROR_PARAM_EMPTY:             "参数不能为空值",
		ERROR_ROLE_IS_UNLOCK:          "角色已解锁",
		ERROR_CURDAY_IS_FINISHED:      "当天签到已完成",
		ERROR_NONEED_SIGNIN:           "本周签到已完成",
		ERROR_NONEED_ADDSIGNIN:        "当前不需要补签或本周已经补签一次，不能再补签",
		ERROR_MUST_SINGIN_CURRDAY:     "请先完成当日签到",
		ERROR_NO_SIGNIN_LOST_CARD:     "没有补签卡",
		ERROR_PUSH_NOTICE:             "推送公告失败",
		ERROR_PUSH_EMAIL:              "推送邮件失败",

		// cdt状态信息.
		ERROR_CDT_DAY_FULL:        "用户当天领取达到限额",
		ERROR_CDT_LACK_OF_BALANCE: "cdt余额不足",
		ERROR_NOT_ENOUGH_ITEM:     "没有足够的物品来兑换",
		// 邮件
		ERROR_EMAIL_NOT_FOUND:              "邮件不存在",
		ERROR_EMAIL_FAILED_TO_CLAIM_REWARD: "领取奖励失败",
		ERROR_EMAIL_REWARD_RECEIVED:        "奖励已领取",

		//道具卡
		ERROR_ITMES_EXPIRED: "已过期",
		ERROR_ROLES_EXPIRED: "角色已过期",
		ERROR_PRIZE_EXPIRED: "奖品已过期",
		//宝箱
		ERROR_BOX_IS_LOCK:         "宝箱已被其他玩家开启",
		ERROR_BOX_OPEN_NUM:        "打开宝箱次数过于频繁，请稍后再尝试",
		ERROR_BOX_AWARD_EXCEPTION: "网络异常，请稍后再试",
		ERROR_BOX_AWARD_NOT_FOUND: "奖励不存在",
		ERROR_BOX_NOT_FINISHED:    "未完成",
		ERROR_BOX_DAY_OPEN_NUM:    "宝箱开启数量已达到上限",
		//活动
		ERROR_NOT_START: "未开始",
		ERROR_END:       "已结束",
	}
)
