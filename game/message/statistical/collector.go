package statistical

import (
	"game_server/db"
	dao "game_server/game/db_service"
	"game_server/game/model"
	"game_server/game/proto"
	"math"
	"strconv"
)

type DataCollector struct {
}

// 总注册数
func (this *DataCollector) GetTotalRegisteCount() int {
	num, err := dao.UserIns.GetAllCount()
	if err != nil {
		return 0
	}
	return num
}

// 单平台总注册数
func (this *DataCollector) GetPlatformTotalRegisteCount(platform int) int {
	num, err := dao.UserIns.GetAllCountByPlatform(platform)
	if err != nil {
		return 0
	}
	return num
}

// 总新增
func (this *DataCollector) GetTotalNewlyAddedCount() int {
	num, err := dao.UserIns.TotalNewlyAddedCount()
	if err != nil {
		return 0
	}
	return num
}

// 单平台总新增
func (this *DataCollector) GetPlatformTotalNewlyAddedCount(platform int) int {
	num, err := dao.UserIns.TotalNewlyAddedCountByPlatform(platform)
	if err != nil {
		return 0
	}
	return num
}

// 注册数:day 格式为:'2017-06-16'
func (this *DataCollector) GetRegisteCount(day string) int {
	num, err := dao.UserIns.CountOfDay(day)
	if err != nil {
		return 0
	}
	return num
}

// 单平台注册数:day 格式为:'2017-06-16'
func (this *DataCollector) GetPlatformRegisteCount(day string, platform int) int {
	num, err := dao.UserIns.CountOfDayAndPlatform(day, platform)
	if err != nil {
		return 0
	}
	return num
}

// 新增数:day 格式为:'2017-06-16'
func (this *DataCollector) GetNewlyAddedCount(day string) int {
	num, err := dao.UserIns.NewlyAddedCount(day)
	if err != nil {
		return 0
	}
	return num
}

// 单平台新增数:day 格式为:'2017-06-16'
func (this *DataCollector) GetPlatformNewlyAddedCount(day string, platform int) int {
	num, err := dao.UserIns.NewlyAddedCountByPlatform(day, platform)
	if err != nil {
		return 0
	}
	return num
}

// 当日登录数,day:20060102
func (this *DataCollector) GetLoginCount(day string) int {
	key := LOGIN_NUM + day
	now := db.RedisMgr.Get(key)
	count, _ := strconv.Atoi(now)
	return count
}

// 当日单一平台登录数,day:20060102
func (this *DataCollector) GetPlatformLoginCount(day string, platform int) int {
	key := PLATFORM_LOGIN_NUM + day + ":" + strconv.Itoa(platform)
	now := db.RedisMgr.Get(key)
	count, _ := strconv.Atoi(now)
	return count
}

// 活跃用户数,day:20060102
func (this *DataCollector) ActiveCount(day string) int {
	// 获取活跃人数
	key := ACTIVE_PEOPLE_NUM + day
	isExists := db.RedisMgr.KeyExist(key)
	if isExists {
		return db.RedisMgr.HLen(key)
	}
	return 0
}

// 单一平台活跃用户数,day:20060102
func (this *DataCollector) PlatformActiveCount(day string, platform int) int {
	// 获取活跃人数
	key := PLATFORM_ACTIVE_PEOPLE_NUM + day + ":" + strconv.Itoa(platform)
	isExists := db.RedisMgr.KeyExist(key)
	if isExists {
		return db.RedisMgr.HLen(key)
	}
	return 0
}

// 活跃用户数,day:20060102
func (this *DataCollector) HourActiveCount(day string, hour int) int {
	// 获取活跃人数
	key := EVERY_HOUR_ACTIVE_PEOPLE_NUM + day + ":" + strconv.Itoa(hour)
	isExists := db.RedisMgr.KeyExist(key)
	if isExists {
		return db.RedisMgr.HLen(key)
	}
	return 0
}

// 在线时长,day:20060102
func (this *DataCollector) OnlineTime(day string) int {
	key := ONLINE_TIME + day
	now := db.RedisMgr.Get(key)
	count, _ := strconv.Atoi(now)
	return count
}

// 平台在线时长,day:20060102
func (this *DataCollector) PlatformOnlineTime(day string, platform int) int {
	key := PLATFORM_ONLINE_TIME + day + ":" + strconv.Itoa(platform)
	now := db.RedisMgr.Get(key)
	count, _ := strconv.Atoi(now)
	return count
}

// 手机累计注册数
func (this *DataCollector) PhoneRegisteCount() int {
	num, err := dao.UserIns.TotalRegisteCount("mobile")
	if err != nil {
		return 0
	}
	return num
}

// 平台手机号累计注册数
func (this *DataCollector) PlatformPhoneRegisteCount(platform int) int {
	num, err := dao.UserIns.TotalRegisteCountByPlatform("mobile", platform)
	if err != nil {
		return 0
	}
	return num
}

// 邮箱累计注册数
func (this *DataCollector) EmailRegisteCount() int {
	num, err := dao.UserIns.TotalRegisteCount("email")
	if err != nil {
		return 0
	}
	return num
}

func (this *DataCollector) PlatformEmailRegisteCount(platform int) int {
	num, err := dao.UserIns.TotalRegisteCountByPlatform("email", platform)
	if err != nil {
		return 0
	}
	return num
}

// 累计实名用户
func (this *DataCollector) RealNameCount() int {
	num, err := dao.UserIns.RealNameCount()
	if err != nil {
		return 0
	}
	return num
}

func (this *DataCollector) PlatformRealNameCount(platform int) int {
	num, err := dao.UserIns.RealNameCountByPlatform(platform)
	if err != nil {
		return 0
	}
	return num
}

// 当日产出CDT:day格式为:"2006-01-02"
func (this *DataCollector) CDTCount(day string) float32 {
	num, err := model.NewCdtRecord().GetOneDayCdt(day, 0)
	if err != nil {
		return 0
	}
	return num
}

// 累计CDT产出总和
func (this *DataCollector) TotalCDTCount() float32 {
	num, err := model.NewCdtRecord().GetTotalCdtProduce()
	if err != nil {
		return 0
	}
	return num
}

// 留存数,day格式为:"2006-01-02"
func (this *DataCollector) GetRetainedCount(day string, diff int) int {
	num, err := dao.UserIns.GetRetainedCount(day, diff)
	if err != nil {
		return 0
	}
	return num
}

// 日宝箱点击总数,day格式为:"20060102"
func (this *DataCollector) GetOpenBoxCount(day string) int {
	// 获取 宝箱数量
	boxKey := BOX_DAY_OPEN_NUM + day
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.Get(boxKey)
		ret, _ := strconv.Atoi(value)
		return ret
	}
	return 0
}

// 日宝箱点击总人数,day格式为:"20060102"
func (this *DataCollector) GetOpenBoxPeopleCount(day string) int {
	// 获取 宝箱数量
	boxKey := BOX_DAY_OPEN_PEOPLE_NUM + day
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		return db.RedisMgr.HLen(boxKey)
	}
	return 0
}

// 日宝箱打开完成总数,day格式为:"20060102"
func (this *DataCollector) GetFinishBoxCount(day string) int {
	// 获取 宝箱数量
	boxKey := BOX_DAY_FINISH_NUM + day
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.Get(boxKey)
		ret, _ := strconv.Atoi(value)
		return ret
	}
	return 0
}

// 日宝箱小时打开完成总数,day格式为:"20060102"
func (this *DataCollector) GetHourFinishBoxCount(day string, hour int) int {
	// 获取 宝箱数量
	boxKey := BOX_DAY_FINISH_HOUR_NUM + day + ":" + strconv.Itoa(hour)
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.Get(boxKey)
		ret, _ := strconv.Atoi(value)
		return ret
	}
	return 0
}

// 日宝箱打开完成获取CDT记录总数,day格式为:"2006-01-02"
func (this *DataCollector) GetBoxCDTCount(day string) int {
	num, err := model.NewCdtRecord().GetOneDayCdtCount(day, proto.MSG_RECEIVE_STREASURE_BOX)
	if err != nil {
		return 0
	}
	return num
}

// 宝箱打开完成获取CDT记录总数
func (this *DataCollector) GetTotalBoxCDTCount() int {
	num, err := model.NewCdtRecord().GetCdtCountByType(proto.MSG_RECEIVE_STREASURE_BOX)
	if err != nil {
		return 0
	}
	return num
}

// 日宝箱打开完成获取CDT总量,day格式为:"2006-01-02"
func (this *DataCollector) GetBoxCDTValue(day string) float32 {
	num, err := model.NewCdtRecord().GetOneDayCdt(day, proto.MSG_RECEIVE_STREASURE_BOX)
	if err != nil {
		return 0
	}
	return num
}

// 日宝箱打开完成获取CDT总量
func (this *DataCollector) GetTotalBoxCDTValue() float32 {
	num, err := model.NewCdtRecord().GetTotalCdtByType(proto.MSG_RECEIVE_STREASURE_BOX)
	if err != nil {
		return 0
	}
	return num
}

// 获取前一天的宝箱总打开次数
func (this *DataCollector) LastGetTotalOpendCount(day string) int {
	result, err := dao.StatisticsTreasureBoxIns.GetLastData()
	if err != nil {
		return 0
	}
	return result.TotalOpenedCount
}

// 参数补签任务每一天人数,day格式为:2020-01-01
func (this *DataCollector) GetSignInPeopleNumEveryDay(dayDate string) ([]int, error) {
	numMap, err := dao.UserTaskIns.EveryDayReward(dayDate+" 00:00:00", dayDate+" 23:59:59")
	if err != nil {
		return nil, err
	}
	everyData := []int{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
	}
	for _, value := range numMap {
		key := value["task_id"]
		key = key[len(key)-1:] //  截取最后一位， 最后一位代表第几天
		str_num := value["number"]
		number, err := strconv.Atoi(str_num)
		if err != nil {
			continue
		}
		index, _ := strconv.Atoi(key)
		if index > 0 {
			index -= 1
		}
		everyData[index] = number
	}
	return everyData, nil
}

// 消耗补签卡数量,day格式为:2020-01-01
func (this *DataCollector) GetUsedMakeUpCardCount(dayDate string) int {
	num, err := dao.UserTaskIns.GetUsedCardCount(dayDate+" 00:00:00", dayDate+" 23:59:59")
	if err != nil {
		return 0
	}
	return num
}

// 获取签到人数,day格式为:2020-01-01
func (this *DataCollector) GetSignInPeopleNum(dayDate string) int {
	number, err := dao.UserTaskIns.SignInData(dayDate+" 00:00:00", dayDate+" 23:59:59")
	if err != nil {
		return 0
	}
	return number
}

// 获取签到次数,day格式为:2020-01-01
func (this *DataCollector) GetSignInCount(dayDate string) int {
	number, err := dao.UserTaskIns.GetSignInCount(dayDate+" 00:00:00", dayDate+" 23:59:59")
	if err != nil {
		return 0
	}
	return number
}

// 新春兑换活动,day格式为:2020-01-01
func (this *DataCollector) GetDoubleYearPV(dayDate string) int {
	num, err := model.NewCdtRecord().GetOneDayCdtCount(dayDate, proto.MSG_SWEET_TREE)
	if err != nil {
		return 0
	}
	return num
}

// 新春兑换活动,day格式为:2020-01-01
func (this *DataCollector) GetDoubleYearUV(dayDate string) int {
	num, err := model.NewCdtRecord().GetOneDayPeopleCount(dayDate, proto.MSG_SWEET_TREE)
	if err != nil {
		return 0
	}
	return num
}

// 新春兑换活动,day格式为:2020-01-01
func (this *DataCollector) GetDoubleYearCdtValue(dayDate string) float32 {
	num, err := model.NewCdtRecord().GetOneDayCdt(dayDate, proto.MSG_SWEET_TREE)
	if err != nil {
		return 0
	}
	return num
}

// 新春兑换活动用户CDT数,day格式为:2020-01-01,结果格式为"user_id","cdt"
func (this *DataCollector) GetDoubleYearUserOneDayCdtValue(dayDate string) []map[string]string {
	num, err := model.NewCdtRecord().UserOneDayCdt(dayDate, proto.MSG_SWEET_TREE)
	if err != nil {
		return nil
	}
	return num
}

// 新春兑换活动用户总CDT数,day格式为:2020-01-01,结果格式为"user_id","cdt"
func (this *DataCollector) GetDoubleYearUserTotalCdtValue(dayDate, userId string) float32 {
	num, err := model.NewCdtRecord().GetUserTotalCdt(dayDate, userId, proto.MSG_SWEET_TREE)
	if err != nil {
		return 0
	}
	return num
}

// 新春活动每日排行榜奖励CDT数,day格式为:2020-01-01
func (this *DataCollector) GetDoubleYearDayCdtValue(dayDate string) float32 {
	num, err := model.NewCdtRecord().GetOneDayCdt(dayDate, proto.MSG_RANK_LIST_DAY_CDT)
	if err != nil {
		return 0
	}
	return num
}

// 新春活动总排行榜奖励CDT数,day格式为:2020-01-01
func (this *DataCollector) GetDoubleYearTotalCdtValue(dayDate string) float32 {
	num, err := model.NewCdtRecord().GetOneDayCdt(dayDate, proto.MSG_RANK_LIST_ALL_CDT)
	if err != nil {
		return 0
	}
	return num
}

// 新春活动福娃碎片每日量,day格式为:20200101,结果格式为"user_id:num"
func (this *DataCollector) GetDoubleYearFuwaUserDayValue(dayDate string) map[string]string {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_FUWA + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.HGetAll(boxKey)
		return value
	}
	return nil
}

// 新春活动福娃碎片每日PV,day格式为:20200101
func (this *DataCollector) GetDoubleYearFuwaPVValue(dayDate string) int {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_FUWA_PV + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.Get(boxKey)
		ret, _ := strconv.Atoi(value)
		return ret
	}
	return 0
}

// 新春活动福娃碎片每日UV,day格式为:20200101
func (this *DataCollector) GetDoubleYearFuwaUVValue(dayDate string) int {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_FUWA_UV + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.HLen(boxKey)
		return value
	}
	return 0
}

func FormatFloat(num float64, decimal int) (float64, error) {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
	return strconv.ParseFloat(res, 64)
}

// 新春活动每日排行榜cdt总数,day格式为:20200101
func (this *DataCollector) GetDoubleYearDayRankingCdt(dayDate string) float32 {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_DAILY_RANKING_CDT + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.Get(boxKey)
		ret, _ := strconv.ParseFloat(value, 64)
		ret, _ = FormatFloat(ret, 4)
		return float32(ret)
	}
	return 0
}

// 新春活动每日排行榜PV,day格式为:20200101
func (this *DataCollector) GetDoubleYearDayRankingPVValue(dayDate string) int {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_DAY_PV + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.Get(boxKey)
		ret, _ := strconv.Atoi(value)
		return ret
	}
	return 0
}

// 新春活动每日排行榜UV,day格式为:20200101
func (this *DataCollector) GetDoubleYearDayRankingUVValue(dayDate string) int {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_DAY_UV + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.HLen(boxKey)
		return value
	}
	return 0
}

// 新春活动总排行榜PV,day格式为:20200101
func (this *DataCollector) GetDoubleYearTotalRankingPVValue(dayDate string) int {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_TOTAL_PV + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.Get(boxKey)
		ret, _ := strconv.Atoi(value)
		return ret
	}
	return 0
}

// 新春活动总排行榜UV,day格式为:20200101
func (this *DataCollector) GetDoubleYearTotalRankingUVValue(dayDate string) int {
	// 获取 宝箱数量
	boxKey := DOUBLEYEAR_TOTAL_UV + dayDate
	isExists := db.RedisMgr.KeyExist(boxKey)
	if isExists {
		value := db.RedisMgr.HLen(boxKey)
		return value
	}
	return 0
}
