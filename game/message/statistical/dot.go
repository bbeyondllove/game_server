package statistical

import (
	"game_server/core/logger"
	"game_server/db"
	"strconv"
	"time"
)

const (
	DOT_REDIS_TIMEOUT = 60 * 60 * 48 // 两天后过期
)

// 统计打点模块
type StatisticsDot struct {
}

func (this *StatisticsDot) Login(user_id string, bRobot bool, platform int, version string) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102
	now_second := time_now.Unix()

	// 每日活跃数
	key := ACTIVE_PEOPLE_NUM + time_str
	_, err := db.RedisMgr.GetRedisClient().HSet(key, user_id, now_second).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// 登陆次数
	key = LOGIN_NUM + time_str
	db.RedisMgr.Incr(key)
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// 每日每时活跃人数
	key = EVERY_HOUR_ACTIVE_PEOPLE_NUM + time_str + ":" + strconv.Itoa(time_now.Local().Hour())
	_, err = db.RedisMgr.GetRedisClient().HSet(key, user_id, now_second).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// 平台相关
	// 活跃数
	key = PLATFORM_ACTIVE_PEOPLE_NUM + time_str + ":" + strconv.Itoa(platform)
	_, err = db.RedisMgr.GetRedisClient().HSet(key, user_id, now_second).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// 登陆次数
	key = PLATFORM_LOGIN_NUM + time_str + ":" + strconv.Itoa(platform)
	db.RedisMgr.Incr(key)
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)
}

func (this *StatisticsDot) Rebind(user_id string) {
}

func (this *StatisticsDot) Logout(user_id string) {
}

func (this *StatisticsDot) Heartbeat(user_id string) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102
	now_second := time_now.Unix()

	// 每日活跃数
	key := ACTIVE_PEOPLE_NUM + time_str
	old_time_ := db.RedisMgr.HGet(key, user_id)
	db.RedisMgr.GetRedisClient().HSet(key, user_id, now_second)
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)
	if len(old_time_) > 0 {
		old_time, _ := strconv.Atoi(old_time_)
		diff_time := now_second - int64(old_time)

		// 在线时长
		online_key := ONLINE_TIME + time_str
		_, err := db.RedisMgr.GetRedisClient().IncrBy(online_key, diff_time).Result()
		if err != nil {
			logger.Error(err)
		}
		db.RedisMgr.Expire(online_key, DOT_REDIS_TIMEOUT)
	}

	// 每日每时活跃人数
	key = EVERY_HOUR_ACTIVE_PEOPLE_NUM + time_str + ":" + strconv.Itoa(time.Now().Local().Hour())
	_, err := db.RedisMgr.GetRedisClient().HSet(key, user_id, now_second).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// 获取平台
	userInfo, err := db.RedisGame.HMGet(user_id, "platform").Result()
	if err != nil {
		return
	}
	platform := userInfo[0].(string)
	// 平台相关

	// 活跃数
	key = PLATFORM_ACTIVE_PEOPLE_NUM + time_str + ":" + platform
	old_time_ = db.RedisMgr.HGet(key, user_id)
	_, err = db.RedisMgr.GetRedisClient().HSet(key, user_id, now_second).Result()
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	if err != nil {
		logger.Error(err)
	}
	if len(old_time_) > 0 {
		old_time, _ := strconv.Atoi(old_time_)
		diff_time := now_second - int64(old_time)

		// 在线时长
		online_key := PLATFORM_ONLINE_TIME + time_str + ":" + platform
		_, err = db.RedisMgr.GetRedisClient().IncrBy(online_key, diff_time).Result()
		if err != nil {
			logger.Error(err)
		}
		db.RedisMgr.Expire(online_key, DOT_REDIS_TIMEOUT)
	}
}

// 打开宝箱
func (this *StatisticsDot) OpenStreasureBox(UserId string) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102

	// 打开宝箱次数
	peopleNumDayKey := BOX_DAY_OPEN_NUM + time_str
	db.RedisMgr.Incr(peopleNumDayKey)
	db.RedisMgr.Expire(peopleNumDayKey, DOT_REDIS_TIMEOUT)

	// 打开宝箱人数
	peopleNumDayKey = BOX_DAY_OPEN_PEOPLE_NUM + time_str
	_, err := db.RedisMgr.GetRedisClient().HSet(peopleNumDayKey, UserId, 1).Result()
	if err != nil {
		logger.Error(err)
	}

	db.RedisMgr.Expire(peopleNumDayKey, DOT_REDIS_TIMEOUT)
}

// 完成宝箱
func (this *StatisticsDot) FinishStreasureBox(user_id string) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102

	// 标记用户参与宝箱活动
	peopleNumDayKey := BOX_DAY_FINISH_NUM + time_str
	db.RedisMgr.Incr(peopleNumDayKey)
	db.RedisMgr.Expire(peopleNumDayKey, DOT_REDIS_TIMEOUT)

	// 每小时宝箱完成数
	peopleNumDayKey = BOX_DAY_FINISH_HOUR_NUM + time_str + ":" + strconv.Itoa(time_now.Local().Hour())
	db.RedisMgr.Incr(peopleNumDayKey)
	db.RedisMgr.Expire(peopleNumDayKey, DOT_REDIS_TIMEOUT)
}

// 新春活动福娃打点
func (this *StatisticsDot) DoubleYearFuwa(user_id string, itemId int, itemNum int) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102

	key := DOUBLEYEAR_FUWA + time_str
	if itemId == StatisticsConfigIns.FuwaFragmentId {
		_, err := db.RedisMgr.GetRedisClient().HIncrBy(key, user_id, int64(itemNum)).Result()
		if err != nil {
			logger.Error(err)
		}

		db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)
	}
}

// 新春活动福娃PVUV打点
func (this *StatisticsDot) DoubleYearFuwaPVUV(user_id string) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102

	//PV
	key := DOUBLEYEAR_FUWA_PV + time_str
	db.RedisMgr.Incr(key)
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// UV
	key = DOUBLEYEAR_FUWA_UV + time_str
	_, err := db.RedisMgr.GetRedisClient().HSet(key, user_id, 1).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)
}

// 新春活动每日PVUV打点
func (this *StatisticsDot) DoubleYearDayPVUV(user_id string) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102

	//PV
	key := DOUBLEYEAR_DAY_PV + time_str
	db.RedisMgr.Incr(key)
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// UV
	key = DOUBLEYEAR_DAY_UV + time_str
	_, err := db.RedisMgr.GetRedisClient().HSet(key, user_id, 1).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)
}

// 新春活动排行榜PVUV打点
func (this *StatisticsDot) DoubleYearRankingListPVUV(user_id string) {
	time_now := time.Now()
	time_str := GetShortTimeStr(time_now) //20060102

	//PV
	key := DOUBLEYEAR_TOTAL_PV + time_str
	db.RedisMgr.Incr(key)
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

	// UV
	key = DOUBLEYEAR_TOTAL_UV + time_str
	_, err := db.RedisMgr.GetRedisClient().HSet(key, user_id, 1).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

}

// 新春活动日排行总计CDT数
func (this *StatisticsDot) DoubleYearDailyRankingCdt(day string, cdt int) {
	if len(day) <= 0 {
		return
	}
	// time_now := time.Now()
	// time_str := GetShortTimeStr(time_now) //20060102

	key := DOUBLEYEAR_DAILY_RANKING_CDT + day
	_, err := db.RedisMgr.GetRedisClient().IncrBy(key, int64(cdt)).Result()
	if err != nil {
		logger.Error(err)
	}
	db.RedisMgr.Expire(key, DOT_REDIS_TIMEOUT)

}
