package double_year

import (
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"math"
	"strconv"
	"time"
)

const (
	// ChristmasPropCdtDayKey 圣诞道具一天对换CDT总数缓存key.
	// 过期时间为第二天零晨00:00.
	ChristmasPropCdtDayKey = "doubleYearCdtDay"
	// ChristmasPropUserKey 用户对换道具所得cdt缓存key.
	// 过期时间为第二天零晨00:00.
	ChristmasPropUserKey = "userPropDay_"
)

// TradeCdt 圣诞活动对换cdt.
type TradeCdt struct {
	B  *Base
	Cr *model.CdtRecord
}

// TradeCdt 道具对换cdt, 返回总对换的cdt和用户对换的cdt.
func (t *TradeCdt) TradeCdt(userId string, sweets, tree int) (map[string]interface{}, int) {
	cdt := db_service.NewCdt()
	userDayCdt, _ := t.GetUserCdt(userId)
	oldTotalDayCdt, _ := t.GetTotalCdt()
	result := map[string]interface{}{
		// 全区当天对换cdt总量.
		"totalDayCdt": formatFloat(oldTotalDayCdt),
		// 用户当天对换cdt总量.
		"userDayCdt": formatFloat(userDayCdt),
		// 用户拥有总的cdt.
		"userAllCdt": formatFloat(cdt.GetUserCdt(userId)),
		// 对换cdt使用sweet数量.
		"sweetUsed": 0,
		// 对换cdt使用tree数量.
		"treeUsed": 0,
		// 本次改变的CDT
		"changeCdt": float32(0),
	}

	if sweets <= 0 || tree <= 0 {
		logger.Errorf("please input legitimate args\n")
		return result, ActiveDoubleYearSweetOrTreeIllegal
	}
	if !t.activityIsOpen() {
		logger.Errorf("the activity is not open!\n")
		return result, ActiveDoubleYearIsNotOpen
	}
	if t.isLimit(oldTotalDayCdt) {
		logger.Errorf("the trade cdt is limited!\n")
		result["totalDayCdt"] = formatFloat(oldTotalDayCdt)
		result["userDayCdt"] = formatFloat(userDayCdt)
		return result, ActiveDoubleYearTradeIsLimit
	}

	// 计算由糖果＋圣诞树组成可对换cdt的数量，取这两个道具中较小者.
	tradeNumber := min(sweets, tree)
	// 使用多少个道具对换cdt，默认为'tradeNumber'个, 当道具对换cdt超过当天cdt上限时，需要计算实际使用了多少个道具.
	usedTradeNumber := tradeNumber
	changeCdt := float32(tradeNumber) * base.Setting.Doubleyear.SweetTreeCdtRate
	if t.isLimit(oldTotalDayCdt + changeCdt) {
		changeCdt = t.limitCdt() - oldTotalDayCdt
		// 四舍五入进一位.
		usedTradeNumber = int(math.Floor(float64(changeCdt / base.Setting.Doubleyear.SweetTreeCdtRate)))
	}
	code, userAllCdt := cdt.UpdateUserCdt(userId, changeCdt, proto.MSG_SWEET_TREE)
	if code != errcode.MSG_SUCCESS {
		logger.Errorf("[tradeCdt] user[%v] trade sweet[%v], tree[%v] to cdt fail.\n", userId, tradeNumber, tradeNumber)
		result["totalDayCdt"] = formatFloat(oldTotalDayCdt)
		result["userDayCdt"] = formatFloat(userDayCdt)
		return result, int(code)
	}

	// 更新当天所有玩家道具对换cdt.
	_, err := db.RedisMgr.GetRedisClient().Set(ChristmasPropCdtDayKey, strconv.FormatFloat(float64(oldTotalDayCdt+changeCdt), 'f', 4, 32), cdt.GetCacheExpireTime(1)).Result()
	if err != nil {
		logger.Errorf("[tradeCdt] user[%v] update totalTradeCdtDay[oldCdt:%v, changeCdt:%v] fail: %v.\n", userId, oldTotalDayCdt, changeCdt, err)
		result["totalDayCdt"] = formatFloat(oldTotalDayCdt)
		result["userDayCdt"] = formatFloat(userDayCdt)
		result["userAllCdt"] = formatFloat(userAllCdt)
		return result, ActiveDoubleYearUpdateAllUserCdtFail
	}

	// 更新当天某个用户对换cdt相关数据.
	_, err = db.RedisMgr.GetRedisClient().Set(ChristmasPropUserKey+userId, strconv.FormatFloat(float64(userDayCdt+changeCdt), 'f', 4, 32), cdt.GetCacheExpireTime(1)).Result()
	if err != nil {
		logger.Errorf("[tradeCdt] user[%v] update userTradeCdtDay[oldCdt:%v, changeCdt:%v] fail: %v.\n", userId, userDayCdt, changeCdt, err)
		result["totalDayCdt"] = formatFloat(oldTotalDayCdt + changeCdt)
		result["userDayCdt"] = formatFloat(userDayCdt)
		result["userAllCdt"] = formatFloat(userAllCdt)
		result["sweetUsed"], result["treeUsed"] = usedTradeNumber, usedTradeNumber
		return result, ActiveDoubleYearUpdateUserCdtFail
	}

	logger.Infof("[tradeCdt] user[%v] trade sweet[%v], tree[%v] to cdt[%v] success.\n", userId, usedTradeNumber, usedTradeNumber, strconv.FormatFloat(float64(changeCdt), 'f', 4, 32))
	result["totalDayCdt"] = formatFloat(oldTotalDayCdt + changeCdt)
	result["userDayCdt"] = formatFloat(userDayCdt + changeCdt)
	result["userAllCdt"] = formatFloat(userAllCdt)
	result["sweetUsed"], result["treeUsed"] = usedTradeNumber, usedTradeNumber
	result["changeCdt"] = changeCdt
	return result, ActiveDoubleYearSuccess
}

// GetTradCdt 获取总对换cdt和个人cdt.
func (t *TradeCdt) GetTradCdt(userId string) (float32, float32) {
	totalCdt, _ := t.GetTotalCdt()
	userCdt, _ := t.GetUserCdt(userId)
	return totalCdt, userCdt
}

// activityIsOpen 判断活动是否开启了.
func (t *TradeCdt) activityIsOpen() bool {
	currentTime := time.Now()
	// 日期判断.
	if currentTime.Unix() < utils.Str2Time(base.Setting.Doubleyear.PropStartDate+" "+base.Setting.Doubleyear.PropStartTime).Unix() || currentTime.Unix() >= utils.Str2Time(base.Setting.Doubleyear.PropEndDate+" "+base.Setting.Doubleyear.PropEndTime).Unix() {
		return false
	}

	// 时间判断.
	//currentDate := time.Unix(currentTime.Unix(), 0).Format("2006-01-02")
	//if currentTime.Unix() < utils.Str2Time(currentDate+" "+base.Setting.Doubleyear.PropStartTime).Unix() || currentTime.Unix() >= utils.Str2Time(currentDate+" "+base.Setting.Doubleyear.PropEndTime).Unix() {
	//	return false
	//}

	return true
}

// GetUserCdt 获取用户道具对换cdt.
func (t *TradeCdt) GetUserCdt(userId string) (float32, error) {
	client := db.RedisMgr.GetRedisClient()
	userDayCdt, err := client.Get(ChristmasPropUserKey + userId).Float32()
	if err != nil {
		logger.Errorf("get prop of user[%v] failed: %v\n", userId, err)
		userDayCdt, err = t.Cr.GetUserOneDayCdt(userId, time.Unix(time.Now().Unix(), 0).Format("2006-01-02"), proto.MSG_SWEET_TREE)
		return userDayCdt, err
	}

	return userDayCdt, nil
}

// GetTotalCdt 获取总对换cdt.
func (t *TradeCdt) GetTotalCdt() (float32, error) {
	client := db.RedisMgr.GetRedisClient()
	totalDayCdt, err := client.Get(ChristmasPropCdtDayKey).Float32()
	if err != nil {
		logger.Errorf("get total cdt of the day from redis fail: %v\n", err)
		// 从数据库查询.
		totalDayCdt, err = t.Cr.GetCdtByEvent(proto.MSG_SWEET_TREE, time.Unix(time.Now().Unix(), 0).Format("2006-01-02"), 1)
		if err != nil {
			return t.limitCdt(), err
		}
		return totalDayCdt, err
	}
	return totalDayCdt, nil
}

// isLimit 判断当天道具对换总cdt是否达到限制值.
func (t *TradeCdt) isLimit(totalCdt float32) bool {
	return totalCdt >= t.limitCdt()
}

// limitCdt 每天道具对换cdt上限值.
func (t *TradeCdt) limitCdt() float32 {
	return base.Setting.Doubleyear.TradeCdtDayLimit
}

// NewTradeCdt 实例化TradeCdt.
func (t *TradeCdt) NewTradeCdt() *TradeCdt {
	return &TradeCdt{}
}

// min 取两个整数中的较小一个.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// formatFloat 格式化float32保住四位小数.
func formatFloat(data float32) string {
	return strconv.FormatFloat(float64(data), 'f', 4, 32)
}
