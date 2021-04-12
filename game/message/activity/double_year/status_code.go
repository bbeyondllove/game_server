package double_year

// 状态码.
const (
	ActiveDoubleYearSuccess = 0
	// ActiveDoubleYearArgsError 参数错误.
	ActiveDoubleYearArgsError = iota + 1000
	// ActiveDoubleYearIllegalToken token不合法.
	ActiveDoubleYearIllegalToken
	// ActiveDoubleYearSweetOrTreeIllegal sweet或tree为0.
	ActiveDoubleYearSweetOrTreeIllegal
	// ActiveDoubleYearIsNotOpen 活动未开启.
	ActiveDoubleYearIsNotOpen
	// ActiveDoubleYearTradeIsLimit 当天总对换cdt达到上限.
	ActiveDoubleYearTradeIsLimit
	// ActiveDoubleYearUpdateAllUserCdtFail 更新当天所有玩家cdt到redis失败.
	ActiveDoubleYearUpdateAllUserCdtFail
	// ActiveDoubleYearUpdateUserCdtFail 更新当天某个玩家cdt到redis失败.
	ActiveDoubleYearUpdateUserCdtFail
	// ActiveDoubleYearServerBusy 系统错误.
	ActiveDoubleYearServerBusy

	// ActiveDoubleYearCdtDaYFull 个人一天获取cdt达到上限，这是兼容更新cdt方法状态码.
	ActiveDoubleYearCdtDaYFull = 5101
)

// StatusCodeMessage 状态码对应信息.
var StatusCodeMessage = map[int]string{
	ActiveDoubleYearSuccess:              "success",
	ActiveDoubleYearArgsError:            "args error",
	ActiveDoubleYearIllegalToken:         "illegal token",
	ActiveDoubleYearSweetOrTreeIllegal:   "sweet or tree is negative",
	ActiveDoubleYearIsNotOpen:            "active is not open",
	ActiveDoubleYearTradeIsLimit:         "今天对换CDT已经达到上限!",
	ActiveDoubleYearUpdateAllUserCdtFail: "update tradeCdtDay fail",
	ActiveDoubleYearUpdateUserCdtFail:    "update one user cdt fail",
	ActiveDoubleYearServerBusy:           "the server busy!",
}
